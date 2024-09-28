package server

import (
	"context"
	"io"
	"sync"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/goccy/go-json"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
)

const listFileScript = `
def path = doc['path'].value;
def dir = params.dir;
if (path.startsWith(dir)) {
	  path = path.substring(dir.length());
	  def index = path.indexOf('/');
	  if (index != -1) {
		    return path.substring(0, index);
	  } else {
		    return path;
	  }
} else {
	  return '';
}
`

const permissionScript = `
if (doc['owner_type'].value == 0 && doc['owner_id'].value == params.operator_id) {
	  return true;
} else if (doc['owner_type'].value == 1 && params.user_group_ids.contains(doc['owner_id'].value)) {
	  return true;
} else {
	  return false;
}
`

var (
	filePool = sync.Pool{
		New: func() interface{} {
			return &model.File{}
		},
	}
	listFileScriptId   = "list_file"
	permissionScriptId = "permission"
)

func (s *GaCloudServer) PutFile(ctx context.Context, f *model.File, reader io.ReadCloser) error {
	// Write file content to provider first
	// TODO: writer lock
	err := s.WriteFile(ctx, f, reader)
	if err != nil {
		return err
	}

	// Write file metadata to elasticsearch
	// if document already exists, update it
	_, err = s.es.Index(s.esIndex).Id(utils.EncodeElasticSearchID(f.Path)).Refresh(refresh.True).Request(f).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *GaCloudServer) DeleteFile(ctx context.Context, path string) error {
	_, err := s.es.Delete(s.esIndex, path).Do(ctx)

	return err
}

func (s *GaCloudServer) GetFileBySum(ctx context.Context, sum string) (*model.File, func(), error) {
	query := types.Query{
		Match: map[string]types.MatchQuery{
			"sum": {
				Query: sum,
			},
		},
	}

	resp, err := s.es.Search().Index(s.esIndex).Query(&query).Do(ctx)
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

func (s *GaCloudServer) GetFileByPath(ctx context.Context, path string) (*model.File, func(), error) {
	resp, err := s.es.Get(s.esIndex, utils.EncodeElasticSearchID(path)).Do(ctx)
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
	provider, ok := s.fileProviders[f.ProviderId]
	if !ok {
		return nil, false, utils.ErrorFileProviderNotFound
	}

	return provider.Get(ctx, f.Sum)
}

func (s *GaCloudServer) WriteFile(ctx context.Context, f *model.File, reader io.Reader) error {
	provider, ok := s.fileProviders[f.ProviderId]
	if !ok {
		return utils.ErrorFileProviderNotFound
	}

	return provider.Put(ctx, f.Sum, reader)
}

func (s *GaCloudServer) FileExists(ctx context.Context, path string) (bool, error) {
	return s.es.Exists(s.esIndex, utils.EncodeElasticSearchID(path)).Do(ctx)
}

func (s *GaCloudServer) ListFiles(ctx context.Context, operator *model.User, dir string) ([]*model.File, func(), error) {
	// groupIds := s.GetUserGroupIds(ctx, operator)
	dir = utils.CleanDirPath(dir)

	searchReq := &search.Request{
		ScriptFields: map[string]types.ScriptField{
			"fd": {
				Script: types.Script{
					Id: &listFileScriptId,
					Params: map[string]json.RawMessage{
						"dir": json.RawMessage(`"` + dir + `"`),
					},
				},
			},
		},
		Source_: types.SourceConfig(true),
		Query: &types.Query{
			MatchAll: &types.MatchAllQuery{},
		},
	}

	resp, err := s.es.Search().Index(s.esIndex).Request(searchReq).Do(ctx)
	if err != nil {
		return nil, func() {}, err
	}

	listFiles := make([]*model.File, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		listFile := filePool.Get().(*model.File)
		err = json.Unmarshal(hit.Source_, listFile)
		if err != nil {
			return nil, func() {}, err
		}
		fd, ok := hit.Fields["fd"]
		if ok {
			listFile.Fd = string(fd)
		}
		listFiles[i] = listFile
	}

	cleanup := func() {
		for _, listFile := range listFiles {
			filePool.Put(listFile)
		}
	}

	return listFiles, cleanup, nil
}
