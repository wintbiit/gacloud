package model

import (
	mime2 "mime"
	"path"
	"strings"
	"time"

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
	Path       string    `json:"path"`
	Size       uint64    `json:"size"`
	Mime       string    `json:"mime"`
	Sum        string    `json:"sum"`
	ProviderId uint      `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func (f *File) Name() string {
	return path.Base(f.Path)
}

var FileTypeMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"sum":         types.NewTextProperty(),
		"path":        types.NewKeywordProperty(),
		"size":        types.NewIntegerNumberProperty(),
		"mime":        types.NewTextProperty(),
		"provider_id": types.NewIntegerNumberProperty(),
		"created_at":  types.NewDateProperty(),
		"updated_at":  types.NewDateProperty(),
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

func NewFile(owner Chrootable, p string, size uint64, sum string, providerId uint) (*File, error) {
	if !path.IsAbs(p) {
		p = path.Join(owner.HomeDir(), p)
	}

	if !strings.HasPrefix(p, UserScopeDir) && !strings.HasPrefix(p, GroupScopeDir) && !strings.HasPrefix(p, ShareScopeDir) {
		return nil, utils.ErrorInvalidPath
	}

	return newFile(p, size, sum, providerId), nil
}

func newFile(p string, size uint64, sum string, providerId uint) *File {
	ext := path.Ext(p)
	mime := mime2.TypeByExtension(ext)
	if mime == "" {
		mime = "application/octet-stream"
	}

	return &File{
		Path:       p,
		Size:       size,
		Mime:       mime,
		Sum:        sum,
		ProviderId: providerId,
		CreatedAt:  time.Now(),
	}
}
