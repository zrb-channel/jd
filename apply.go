package jd

import (
	"context"
	"errors"
	"github.com/zrb-channel/utils"
	"net/http"

	json "github.com/json-iterator/go"
)

func Apply(ctx context.Context, conf *Config, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	base := CreateBaseRequest(conf.AppId, "open_pushApplyData")

	if err := base.Sign(conf, req); err != nil {
		return nil, errors.New("提交失败")
	}

	resp, err := utils.
		Request(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(base).
		Post(CreateOrderAddr)

	if err != nil {
		return nil, errors.New("提交失败")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("提交失败")
	}

	result, err := UnmarshalResult(conf.AesKey, resp.Body())
	if err != nil {
		return nil, errors.New("提交失败")
	}

	if result.Code != "S00000" {
		return nil, errors.New("提交失败")
	}

	createResp := &CreateOrderResponse{}
	if err = json.Unmarshal([]byte(result.Data), createResp); err != nil {
		return nil, errors.New("提交失败")
	}

	return createResp, nil
}
