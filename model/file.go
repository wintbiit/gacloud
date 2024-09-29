package model

import (
	mime2 "mime"
	"path"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
)

const (
	UserScopeDir  = "/home/user/"
	GroupScopeDir = "/home/group/"
	ShareScopeDir = "/share/"
)

type File struct {
	Path       string `json:"path"`
	Size       uint64 `json:"size"`
	Mime       string `json:"mime"`
	Sum        string `json:"sum"`
	ProviderId uint   `json:"provider_id"`
	Fp         string `json:"fp,omitempty"`
}

var FileTypeMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"sum":         types.NewTextProperty(),
		"path":        types.NewKeywordProperty(),
		"size":        types.NewIntegerNumberProperty(),
		"mime":        types.NewTextProperty(),
		"provider_id": types.NewIntegerNumberProperty(),
	},
}

type FileProvider struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Type       string `gorm:"not null"`
	Credential string `gorm:"not null"`
}

type Chrootable interface {
	HomeDir() string
}

func NewFile(owner Chrootable, dir, name string, size uint64, sum string, providerId uint) (*File, error) {
	if !path.IsAbs(dir) {
		dir = path.Join(owner.HomeDir(), dir)
	}

	if !strings.HasPrefix(dir, UserScopeDir) {
		return nil, utils.ErrorInvalidPath
	}

	if !strings.HasPrefix(dir, GroupScopeDir) {
		return nil, utils.ErrorInvalidPath
	}

	if !strings.HasPrefix(dir, ShareScopeDir) {
		return nil, utils.ErrorInvalidPath
	}

	return newFile(dir, name, size, sum, providerId), nil
}

func newFile(dir, name string, size uint64, sum string, providerId uint) *File {
	ext := path.Ext(name)
	mime, _, err := mime2.ParseMediaType(ext)
	if err != nil {
		mime = "application/octet-stream"
	}

	return &File{
		Path:       path.Join(dir, name),
		Size:       size,
		Mime:       mime,
		Sum:        sum,
		ProviderId: providerId,
	}
}
