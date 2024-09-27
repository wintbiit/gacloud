package server

import (
	"context"
	"io"
	"sync"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/goccy/go-json"
	"github.com/wintbiit/gacloud/fs"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
)

const listFileScript = `
if (doc['path'].value.startsWith(params.dir)) {
  def path = doc['path'].value.substring(params.dir.length());
  def index = path.indexOf('/');
  if (index == -1) {
	return ['name': path, 'is_dir': false];
  } else {
	return ['name': path.substring(0, index), 'is_dir': true];
  }
} else {
  return null;
}
`

var (
	fileProviders map[int64]fs.FileProvider
	filePool      = sync.Pool{
		New: func() interface{} {
			return &model.File{}
		},
	}
	listFilePool = sync.Pool{
		New: func() interface{} {
			return &model.ListFile{}
		},
	}
	listFileScriptId = "list_file"
)

func (s *GaCloudServer) PutFile(ctx context.Context, f *model.File, reader io.ReadCloser) error {
	// Write file content to provider first
	writer, err := s.GetFileWriter(ctx, f)
	if err != nil {
		return err
	}

	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}

	// Write file metadata to elasticsearch
	// if document already exists, update it
	exists, err := s.FileExists(ctx, f.Sum)
	if err != nil {
		return err
	}

	if !exists {
		_, err = s.es.Index(elasticSearchIndex).Id(f.Sum).Request(f).Do(ctx)
		if err != nil {
			return err
		}
	} else {
		bytes, err := json.Marshal(f)
		if err != nil {
			return err
		}

		_, err = s.es.Update(elasticSearchIndex, f.Sum).Request(&update.Request{
			Doc: bytes,
		}).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GaCloudServer) DeleteFile(ctx context.Context, sum string) error {
	_, err := s.es.Delete(elasticSearchIndex, sum).Do(ctx)

	return err
}

func (s *GaCloudServer) GetFileByPath(ctx context.Context, p string) (*model.File, func(), error) {
	query := types.Query{
		Match: map[string]types.MatchQuery{
			"path": {
				Query: p,
			},
		},
	}

	resp, err := s.es.Search().Index(elasticSearchIndex).Query(&query).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	if resp.Hits.Total.Value == 0 {
		return nil, nil, utils.ErrorFileNotFound
	}

	file := filePool.Get().(*model.File)
	err = json.Unmarshal(resp.Hits.Hits[0].Source_, file)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		filePool.Put(file)
	}

	return file, cleanup, nil
}

func (s *GaCloudServer) GetFileBySum(ctx context.Context, sum string) (*model.File, func(), error) {
	resp, err := s.es.Get(elasticSearchIndex, sum).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	if !resp.Found {
		return nil, nil, utils.ErrorFileNotFound
	}

	file := filePool.Get().(*model.File)
	cleanup := func() {
		filePool.Put(file)
	}
	err = json.Unmarshal(resp.Source_, file)
	if err != nil {
		return nil, cleanup, err
	}

	return file, cleanup, nil
}

func (s *GaCloudServer) GetFileReader(ctx context.Context, f *model.File) (io.Reader, bool, error) {
	provider, ok := fileProviders[f.ProviderId]
	if !ok {
		return nil, false, utils.ErrorFileProviderNotFound
	}

	return provider.Get(ctx, f.Sum)
}

func (s *GaCloudServer) GetFileWriter(ctx context.Context, f *model.File) (io.WriteCloser, error) {
	provider, ok := fileProviders[f.ProviderId]
	if !ok {
		return nil, utils.ErrorFileProviderNotFound
	}

	return provider.Put(ctx, f.Sum)
}

func (s *GaCloudServer) FileExists(ctx context.Context, sum string) (bool, error) {
	return s.es.Exists(elasticSearchIndex, sum).Do(ctx)
}

func (s *GaCloudServer) ListFiles(ctx context.Context, operator *model.User, dir string) ([]*model.ListFile, func(), error) {
	query := types.Query{
		Script: &types.ScriptQuery{
			Script: types.Script{
				Id: &listFileScriptId,
				Params: map[string]json.RawMessage{
					"dir": json.RawMessage(`"` + dir + `"`),
				},
			},
		},
	}

	resp, err := s.es.Search().Index(elasticSearchIndex).Query(&query).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	listFiles := make([]*model.ListFile, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		listFile := listFilePool.Get().(*model.ListFile)
		err = json.Unmarshal(hit.Source_, listFile)
		if err != nil {
			return nil, nil, err
		}
		listFiles[i] = listFile
	}

	cleanup := func() {
		for _, listFile := range listFiles {
			listFilePool.Put(listFile)
		}
	}

	return listFiles, cleanup, nil
}
