package server

import (
	"context"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
	"testing"
)

func TestFileSystem(t *testing.T) {
	s, err := NewLocalGaCloudServer(elasticSearchIndex)
	if err != nil {
		t.Fatalf("failed to setup server: %v", err)
	}

	testUser := &model.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "admin",
	}

	testFiles := make([]*model.File, 5)
	for i := 0; i < 5; i++ {
		name := utils.RandStr(10)
		sum := utils.Md5SumBytes([]byte(name))
		testFiles[i], err = model.NewFile(testUser, name, 100, sum, DefaultFileProviderId)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	ctx := context.TODO()
	_, err = s.es.Indices.Delete(elasticSearchIndex).Do(ctx)
	if err != nil {
		t.Fatalf("failed to clear index: %v", err)
	}

	s, err = NewLocalGaCloudServer(elasticSearchIndex)
	if err != nil {
		t.Fatalf("failed to setup server: %v", err)
	}

	for _, f := range testFiles {
		err = s.PutFile(ctx, testUser, f)
		if err != nil {
			t.Fatalf("failed to put file: %v", err)
		}
	}

	files, clean, err := s.ListFiles(ctx, testUser, "")
	defer clean()
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}

	t.Logf("files: %v", files)
	for i, f := range files {
		t.Logf("file %d: %+v", i, f)
	}
}
