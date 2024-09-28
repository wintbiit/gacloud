package server

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"

	"github.com/wintbiit/gacloud/fs"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
)

func TestFileSystem(t *testing.T) {
	db, err := utils.OpenDB("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	err = model.MigrateModels(db)
	if err != nil {
		t.Fatal(err)
	}

	es, err := utils.OpenElasticSearch("http://172.30.162.53:9200", "elastic", "G8gAdnzJ57Wls5feiKI2")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()

	resp, err := es.PutScript(listFileScriptId).Script(&types.StoredScript{
		Lang: scriptlanguage.ScriptLanguage{
			Name: "painless",
		},
		Source: listFileScript,
	}).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Acknowledged {
		t.Fatal(utils.ErrorElasticSearchScriptNotAcknowledged)
	}

	resp, err = es.PutScript(permissionScriptId).Script(&types.StoredScript{
		Lang: scriptlanguage.ScriptLanguage{
			Name: "painless",
		},
		Source: permissionScript,
	}).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Acknowledged {
		t.Fatal(utils.ErrorElasticSearchScriptNotAcknowledged)
	}

	TestClearIndex(t)

	getResp, err := es.GetScript(listFileScriptId).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("GetScript listFileScript: %+v", getResp.Script)

	localFs, err := fs.NewLocalFileProvider([]byte(`{"mount_dir": "./data/gacloud"}`))
	if err != nil {
		t.Fatal(err)
	}

	server := GaCloudServer{
		db:      db,
		es:      es,
		esIndex: elasticSearchIndex + "_test",
		fileProviders: map[uint]fs.FileProvider{
			0: localFs,
		},
	}

	files := []*model.File{
		randomFile("/home/u1/testFile1", model.FileOwnerTypeUser, 1),
		randomFile("/home/u1/testDir1/testFile2", model.FileOwnerTypeUser, 1),
		randomFile("/home/u1/testDir1/testFile3", model.FileOwnerTypeUser, 1),
		randomFile("/home/g1/testFile3", model.FileOwnerTypeGroup, 1),
		randomFile("/home/g1/testDir2/testFile4", model.FileOwnerTypeGroup, 1),
	}

	for _, f := range files {
		err := server.PutFile(ctx, f, utils.ToReadCloser(strings.NewReader(utils.RandStr(1024))))
		if err != nil {
			t.Fatal(err)
		}
	}

	u1 := model.User{
		Model: gorm.Model{
			ID: 1,
		},
	}

	g1 := model.Group{
		Model: gorm.Model{
			ID: 1,
		},
	}

	u2 := model.User{
		Model: gorm.Model{
			ID: 2,
		},
	}

	g2 := model.Group{
		Model: gorm.Model{
			ID: 2,
		},
	}

	err = server.UserAddGroup(ctx, &u1, &g1)
	if err != nil {
		t.Fatal(err)
	}
	err = server.UserAddGroup(ctx, &u2, &g2)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// Test ListFiles
	listFiles, cancel, err := server.ListFiles(ctx, &u1, "/home/u1")
	defer cancel()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ListFiles u1 count: %v", len(listFiles))

	for _, f := range listFiles {
		t.Logf("ListFiles u1: %+v", f)
	}

	listFiles, cancel, err = server.ListFiles(ctx, &u1, "/home/g1")
	defer cancel()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ListFiles g1 count: %v", len(listFiles))

	for _, f := range listFiles {
		t.Logf("ListFiles g1: %+v", f)
	}

	listFiles, cancel, err = server.ListFiles(ctx, &u2, "/home/u2")
	defer cancel()

	if err != nil {
		t.Fatal(err)
	}

	for _, f := range listFiles {
		t.Logf("ListFiles u2: %v", f)
	}

	listFiles, cancel, err = server.ListFiles(ctx, &u2, "/home/g2")
	defer cancel()

	if err != nil {
		t.Fatal(err)
	}

	for _, f := range listFiles {
		t.Logf("ListFiles g2: %v", f)
	}
}

func randomFile(path string, ownerType int8, ownerId uint) *model.File {
	content := utils.RandStr(1024)
	return &model.File{
		Path:      path,
		Size:      int64(len(content)),
		Mime:      "text/plain",
		OwnerType: ownerType,
		OwnerId:   ownerId,
		Sum:       utils.Sha256SumBytes([]byte(content)),
	}
}

func TestClearIndex(t *testing.T) {
	es, err := utils.OpenElasticSearch("http://172.30.162.53:9200", "elastic", "G8gAdnzJ57Wls5feiKI2")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	es.Indices.Delete(elasticSearchIndex).Do(ctx)
	es.Indices.Delete(elasticSearchIndex + "_test").Do(ctx)
}
