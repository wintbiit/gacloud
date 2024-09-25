package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/internal"
	"time"
)

type registry struct {
	party  string
	cb     func(iris.Party)
	behind bool
}

var hooks = make([]registry, 0)

func addHook(party string, behind bool, cb func(iris.Party)) {
	hooks = append(hooks, registry{party, cb, behind})
}

func RegisterRoutes(app iris.Party) {
	for _, hook := range hooks {
		party := app.Party(hook.party)
		party.Use(coreFunction())
		hook.cb(party)
	}
}

func coreFunction() iris.Handler {
	lastCheckedTime := time.Now()
	lastCheckPass := false

	return func(ctx iris.Context) {
		if !lastCheckPass && time.Since(lastCheckedTime) > 30*time.Second {
			checkPass := internal.GetGaCloudServer() != nil && internal.GetGaCloudServer().HealthCheck()
			lastCheckPass = checkPass
			lastCheckedTime = time.Now()
		}

		if !lastCheckPass {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("server not ready")
			ctx.StopExecution()
			return
		}

		ctx.Next()
	}
}
