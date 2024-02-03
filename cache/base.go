package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"omo.msa.third/config"
	"omo.msa.third/proxy/nosql"
	"os"
)

const DefaultScene = "system"

type baseInfo struct {
	ID       uint64 `json:"-"`
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Creator  string
	Operator string
	Created  int64
	Updated  int64
}

type cacheContext struct {
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		num := nosql.GetPartnerCount()
		count := nosql.GetChannelCount()
		logger.Infof("the partner count = %d and the channel count = %d", num, count)
		//nosql.CheckTimes()
	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func (mine *cacheContext) MD5File(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}

func PostHttp(address string, bts []byte, more bool) (*gjson.Result, int, error) {
	msg := string(bts)
	if more {
		logger.Info("post request that address = " + address + ";params = " + msg)
	} else {
		//logger.Info("post request that address = " + address + ";params = " + msg)
		logger.Info(fmt.Sprintf("push request that address = %s; params length = %d;params =%s", address, len(msg), msg))
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", address, bytes.NewReader(bts))
	if err != nil {
		return nil, -7, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("apikey", config.Schema.Basic.OgmToken)
	resp, err := client.Do(req)
	//resp, err := http.Post(address, "application/json;charset=UTF-8", bytes.NewReader(bts))
	if err != nil {
		return nil, -8, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -9, err
	}
	data := bytes.NewBuffer(body).String()
	if len(data) < 200 {
		logger.Info("post response of " + address + " ;that data = " + data)
	} else {
		logger.Infof("post response of %s that data length = %d", address, len(data))
	}

	result := gjson.Parse(data)
	errCode := result.Get("status.code").Int()
	errMsg := result.Get("status.message").String()
	if errCode != 0 {
		return nil, int(errCode), errors.New(errMsg)
	}
	return &result, 0, nil
}

func CheckPage[T any](page, number uint32, all []T) (uint32, uint32, []T) {
	if len(all) < 1 {
		return 0, 0, make([]T, 0, 1)
	}
	if number < 1 {
		number = 10
	}
	total := uint32(len(all))
	if len(all) <= int(number) {
		return total, 1, all
	}
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}
	list := make([]T, 0, number)
	list = append(all[start:end])
	return total, maxPage, list
}
