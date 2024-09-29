package server

import (
	"context"
	"io"
	"sync"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/runtimefieldtype"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/goccy/go-json"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
)

var listFileScript = `
String p = doc['path'].value;
p = p.substring(params.dir.length());
def index = p.indexOf('/');
if (index > 0) {
  emit(p.substring(0, index + 1));
} else {
  emit(p)
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
	permissionScriptId = "permission"
	listSearchSize     = 1000
	fp                 = "fp"
	listSearchAggSize  = 1
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

type SearchAggs struct {
	Buckets []struct {
		Key         string `json:"key"`
		AnyDocument struct {
			Hits struct {
				Hits []model.File `json:"hits"`
			} `json:"hits"`
		} `json:"any_document"`
	} `json:"buckets"`
}

func (s *GaCloudServer) ListFiles(ctx context.Context, operator *model.User, dir string) ([]*model.File, func(), error) {
	// groupIds := s.GetUserGroupIds(ctx, operator)
	dir = utils.CleanDirPath(dir)

	searchReq := search.Request{
		Size: &listSearchSize,
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						Prefix: map[string]types.PrefixQuery{
							"path": {
								Value: dir,
							},
						},
					},
				},
			},
		},
		RuntimeMappings: map[string]types.RuntimeField{
			"fp": {
				Type: runtimefieldtype.Keyword,
				Script: &types.Script{
					Source: &listFileScript,
					Params: map[string]json.RawMessage{
						"dir": utils.JsonRaw(dir),
					},
				},
			},
		},
		Aggregations: map[string]types.Aggregations{
			"group_by_fp": {
				Terms: &types.TermsAggregation{
					Field: &fp,
				},
				Aggregations: map[string]types.Aggregations{
					"any_document": {
						TopHits: &types.TopHitsAggregation{
							Size: &listSearchAggSize,
						},
					},
				},
			},
		},
	}

	resp, err := s.es.Search().Index(s.esIndex).Request(&searchReq).Do(ctx)
	if err != nil {
		return nil, func() {}, err
	}

	aggs, ok := resp.Aggregations["group_by_fp"]
	if !ok {
		return nil, func() {}, nil
	}

	searchAggs, ok := aggs.(*types.StringTermsAggregate)
	if !ok {
		return nil, func() {}, nil
	}

	buckets, ok := searchAggs.Buckets.([]types.StringTermsBucket)
	if !ok {
		return nil, func() {}, nil
	}

	files := make([]*model.File, len(buckets))
	for i, bucket := range buckets {
		files[i] = filePool.Get().(*model.File)
		hit := bucket.Aggregations["any_document"].(*types.TopHitsAggregate).Hits.Hits[0]
		err = json.Unmarshal(hit.Source_, files[i])
		if err != nil {
			return nil, func() {}, err
		}

		files[i].Fp = bucket.Key.(string)
	}

	cleanup := func() {
		for _, f := range files {
			filePool.Put(f)
		}
	}

	return files, cleanup, nil
}
