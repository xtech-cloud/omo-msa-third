package cache

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/micro/go-micro/v2/logger"
	"io"
	"omo.msa.third/config"
	"omo.msa.third/proxy/nosql"
	"os"
	"time"
)

type baseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator string
	Operator string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
	//boxes []*OwnerInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		num := nosql.GetPartnerCount()
		count := nosql.GetChannelCount()
		logger.Infof("the favorite count = %d and the repertory count = %d", num, count)
	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func (mine *cacheContext)MD5File(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
