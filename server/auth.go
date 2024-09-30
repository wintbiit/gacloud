package server

import (
	"context"
	"fmt"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"strconv"
	"strings"
	"time"
)

type DownloadSigner struct {
	Pkey     string
	Duration time.Duration
}

func (s *DownloadSigner) GenerateSign(path string, uid uint) string {
	rand := utils.RandStr(16)
	timestamp := time.Now().Unix()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	pattern := fmt.Sprintf("%s-%d-%s-%d-%s", path, timestamp, rand, uid, s.Pkey)
	sign := utils.Md5SumBytes([]byte(pattern))

	return fmt.Sprintf("%d-%s-%d-%s", timestamp, rand, uid, sign)
}

func (s *DownloadSigner) VerifySign(sum string, timestamp int64, rand string, uid uint, sign string) bool {
	if !strings.HasPrefix(sum, "/") {
		sum = "/" + sum
	}

	ts := time.Unix(timestamp, 0)
	if time.Since(ts) > s.Duration {
		return false
	}

	pattern := fmt.Sprintf("%s-%d-%s-%d-%s", sum, timestamp, rand, uid, s.Pkey)
	return sign == utils.Md5SumBytes([]byte(pattern))
}

func (s *GaCloudServer) GenerateDownloadSign(ctx context.Context, operator *model.User, file *model.File) (string, error) {
	if file == nil {
		return "", utils.ErrorFileNotFound
	}

	return s.downloadSigner.GenerateSign(file.Sum, operator.ID), nil
}

func (s *GaCloudServer) VerifyDownloadSign(ctx context.Context, sum string, sign string) bool {
	parts := strings.SplitN(sum, "-", 4)
	if len(parts) != 4 {
		return false
	}

	timestamp, rand, uid, sum := parts[0], parts[1], parts[2], parts[3]
	timeInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	uidInt, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return false
	}

	return s.downloadSigner.VerifySign(sum, timeInt, rand, uint(uidInt), sign)
}
