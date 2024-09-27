package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/wintbiit/gacloud/utils"
)

func TestEsConn(t *testing.T) {
	t.Log("TestEsConn")

	c, err := utils.OpenElasticSearch("http://172.30.162.53:9200", "elastic", "RobotLab2024", "test")
	if err != nil {
		t.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get all existing indices
	indices, err := c.Cat.Indices().Do(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	for _, index := range indices {
		t.Logf("%s", *index.Index)
	}

	// close connection, remove test index
	resp, err := c.Indices.Delete("test").Do(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp)
}
