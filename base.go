package jd

import (
	"encoding/hex"
	"github.com/google/go-querystring/query"
	json "github.com/json-iterator/go"
	"net/url"
	"github.com/zrb-channel/jd/lib"
	"github.com/zrb-channel/utils/hash"
	"strings"
	"time"
)

func CreateBaseRequest(appId string, method string) *BaseRequest {
	return &BaseRequest{
		AppId:     appId,
		Method:    method,
		Timestamp: time.Now().UnixNano() / 1e6,
	}
}

// Sign
// @param conf
// @param body
// @date 2022-09-21 16:53:33
func (req *BaseRequest) Sign(conf *Config, body SignRequest) error {
	urls, err := query.Values(body)
	if err != nil {
		return err
	}

	var value string
	if value, err = url.QueryUnescape(urls.Encode()); err != nil {
		return err
	}

	sign := strings.ToUpper(hash.MD5String(value + "&key=" + conf.AppSecret))
	body.SetSign(sign)

	var bodyJSON []byte
	if bodyJSON, err = json.Marshal(body); err != nil {
		return err
	}

	var enData []byte
	if enData, err = lib.AesEncryptECB(bodyJSON, []byte(conf.AesKey)); err != nil {
		return err
	}

	req.Param = hex.EncodeToString(enData)
	return nil
}
