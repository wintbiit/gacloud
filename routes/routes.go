package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/server"
)

type registry struct {
	party  string
	cb     func(iris.Party)
	behind bool
	auth   bool
}

var hooks = make([]registry, 0)

func addHook(party string, cb func(iris.Party)) {
	hooks = append(hooks, registry{party, cb, true, false})
}

func addHookFront(party string, cb func(iris.Party)) {
	hooks = append(hooks, registry{party, cb, false, false})
}

func addHookAuth(party string, cb func(iris.Party)) {
	hooks = append(hooks, registry{party, cb, true, true})
}

func RegisterRoutes(app iris.Party) {
	for _, hook := range hooks {
		party := app.Party(hook.party)
		if hook.behind {
			party.Use(coreCheck())
		}
		if hook.auth {
			party.Use(userMiddleware)
		}
		hook.cb(party)
	}
}

func coreCheck() iris.Handler {
	return func(ctx iris.Context) {
		if server.GetServer() == nil {
			ctx.StatusCode(iris.StatusSiteFrozen)
			ctx.WriteString("server not ready")
			ctx.StopExecution()
			return
		}

		if server.GetServer().Maintenance {
			ctx.StatusCode(iris.StatusServiceUnavailable)
			ctx.WriteString("server maintenance")
			ctx.StopExecution()
			return
		}

		ctx.Next()
	}
}

func userMiddleware(ctx iris.Context) {
	s := server.GetServer()

	verify := s.GetUserMiddleware()
	if verify == nil {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	verify(ctx)
}
