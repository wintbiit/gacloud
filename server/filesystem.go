package server

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"io"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

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

var (
	filePool = sync.Pool{
		New: func() interface{} {
			return &model.File{}
		},
	}
	listSearchSize    = 1000
	fp                = "fp"
	listSearchAggSize = 1
)

// FsPermCheck checks if the operator has permission to access the path
// both on user scope and group scope
func (s *GaCloudServer) FsPermCheck(ctx context.Context, operator *model.User, p *string) bool {
	if operator == nil {
		return false
	}

	// if relative path, convert to absolute path in user scope
	if !path.IsAbs(*p) {
		*p = path.Join(operator.HomeDir(), *p)
	}

	// in user scope
	if strings.HasPrefix(*p, operator.HomeDir()) {
		return true
	}

	// in group scope
	if strings.HasPrefix(*p, model.GroupScopeDir) {
		sp := strings.SplitN(*p, "/", 5)
		if len(sp) != 5 {
			return false
		}

		gidStr := sp[3]
		gid, err := strconv.ParseUint(gidStr, 10, 64)
		if err != nil {
			return false
		}

		if !s.IsUserInGroup(ctx, operator, uint(gid)) {
			return false
		}

		return true
	}

	return false
}

// PutFile puts file metadata to elasticsearch
func (s *GaCloudServer) PutFile(ctx context.Context, operator *model.User, f *model.File) error {
	if !s.FsPermCheck(ctx, operator, &f.Path) {
		return utils.ErrorPermissionDenied
	}

	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now()
	}

	f.UpdatedAt = time.Now()

	// Write file metadata to elasticsearch
	// if document already exists, update it
	_, err := s.es.Index(s.esIndex).Id(utils.EncodeElasticSearchID(f.Path)).Refresh(refresh.True).Request(f).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

// WriteFile writes file content to file provider and returns the checksum
func (s *GaCloudServer) WriteFile(ctx context.Context, operator *model.User, f *model.File, reader io.Reader) (string, error) {
	if !s.FsPermCheck(ctx, operator, &f.Path) {
		return "", utils.ErrorPermissionDenied
	}

	provider, err := s.GetProvider(f.ProviderId)
	if err != nil {
		return "", err
	}

	return provider.Put(ctx, reader)
}

func (s *GaCloudServer) DeleteFile(ctx context.Context, operator *model.User, path string) error {
	if !s.FsPermCheck(ctx, operator, &path) {
		return utils.ErrorPermissionDenied
	}

	_, err := s.es.Delete(s.esIndex, path).Do(ctx)

	return err
}

func (s *GaCloudServer) GetFileBySum(ctx context.Context, operator *model.User, sum string) (*model.File, func(), error) {
	if !s.FsPermCheck(ctx, operator, &sum) {
		return nil, nil, utils.ErrorPermissionDenied
	}

	return s.GetFileBySumNoCheck(ctx, sum)
}

func (s *GaCloudServer) GetFileBySumNoCheck(ctx context.Context, sum string) (*model.File, func(), error) {
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

func (s *GaCloudServer) GetFileReader(ctx context.Context, operator *model.User, f *model.File) (io.Reader, bool, error) {
	if !s.FsPermCheck(ctx, operator, &f.Path) {
		return nil, false, utils.ErrorPermissionDenied
	}

	provider, err := s.GetProvider(f.ProviderId)
	if err != nil {
		return nil, false, err
	}

	return provider.Get(ctx, f.Sum)
}

func (s *GaCloudServer) FileExists(ctx context.Context, operator *model.User, path string) (bool, error) {
	if !s.FsPermCheck(ctx, operator, &path) {
		return false, utils.ErrorPermissionDenied
	}

	return s.es.Exists(s.esIndex, utils.EncodeElasticSearchID(path)).Do(ctx)
}

func (s *GaCloudServer) ListFiles(ctx context.Context, operator *model.User, dir string) ([]*model.File, func(), error) {
	if !s.FsPermCheck(ctx, operator, &dir) {
		return nil, func() {}, utils.ErrorPermissionDenied
	}

	return s.listFiles(ctx, dir)
}

// listFiles lists files in a directory, returns a list of files and a cleanup function
// with no permission check
func (s *GaCloudServer) listFiles(ctx context.Context, dir string) ([]*model.File, func(), error) {
	dir = utils.CleanDirPath(dir)
	size := "size"

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
					"total_size": {
						Sum: &types.SumAggregation{
							Field: &size,
						},
					},
					"latest_document": {
						TopHits: &types.TopHitsAggregation{
							Size:    &listSearchAggSize,
							Source_: []string{"updated_at", "mime", "sum"},
							Sort: []types.SortCombinations{
								types.SortOptions{
									SortOptions: map[string]types.FieldSort{
										"updated_at": {
											Order: &sortorder.Desc,
										},
									},
								},
							},
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
		totalSize := bucket.Aggregations["total_size"].(*types.SumAggregate).Value
		files[i].Size = uint64(*totalSize)
		latestDoc := bucket.Aggregations["latest_document"].(*types.TopHitsAggregate).Hits.Hits[0]
		err = json.Unmarshal(latestDoc.Source_, files[i])
		if err != nil {
			return nil, func() {}, err
		}

		files[i].Path = bucket.Key.(string)
	}

	cleanup := func() {
		for _, f := range files {
			filePool.Put(f)
		}
	}

	return files, cleanup, nil
}

func (s *GaCloudServer) ClearIndex(ctx context.Context) error {
	_, err := s.es.DeleteByQuery(s.esIndex).Query(&types.Query{
		MatchAll: &types.MatchAllQuery{},
	}).Do(ctx)
	return err
}
