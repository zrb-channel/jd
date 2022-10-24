package jd

import (
	"context"
	"errors"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/zrb-channel/utils"
	log "github.com/zrb-channel/utils/logger"
	"go.uber.org/zap"
	"net/http"
)

// QueryCredit
// @param ctx
// @param conf
// @param orderId
// @date 2022-10-10 18:14:02
func QueryCredit(ctx context.Context, conf *Config, orderId string) (*Order, error) {
	base := CreateBaseRequest(conf.AppId, "open_queryCreditStatus")
	req := &QueryCreditRequest{
		ProductId: ProductId,
		CreditList: []*OrderItem{
			{OrderID: orderId},
		},
	}

	if err := base.Sign(conf, req); err != nil {
		log.WithError(err).Error("签名失败")
		return nil, err
	}

	resp, err := utils.
		Request(ctx).
		SetBody(base).
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		Post(Addr)

	if err != nil {
		log.WithError(err).Error("请求失败")
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.WithError(err).Error("返回状态吗有误", zap.Int("statusCode", resp.StatusCode()))
		return nil, errors.New("返回状态码有误")
	}

	msg := &QueryCreditResponse{}
	result, err := UnmarshalResult(conf.AesKey, resp.Body())
	if err != nil {
		log.WithError(err).Error("解析返回数据失败", zap.Any("data", map[string]any{"body": base, "req": req}))
		return nil, err
	}

	if result.Data == "" {
		log.WithError(err).Error("解析返回数据失败", zap.Any("data", map[string]any{"body": base, "req": req, "uri": resp.Request.URL}))
		return nil, errors.New("解析后的数据为空")
	}

	if err = json.Unmarshal([]byte(result.Data), msg); err != nil {

		log.WithError(err).Error("解析返回数据失败1")
		return nil, err
	}

	if len(msg.CreditList) == 0 {
		return nil, errors.New("未查询到订单")
	}

	return msg.CreditList[0], nil
}

// QueryLoan
// @param ctx
// @param conf
// @param orderId
// @date 2022-10-10 18:14:01
func QueryLoan(ctx context.Context, conf *Config, orderId string) (*OrderLoan, error) {
	base := CreateBaseRequest(conf.AppId, "open_queryLoanStatus")
	req := &QueryCreditRequest{
		ProductId:  ProductId,
		CreditList: []*OrderItem{{OrderID: orderId}},
	}

	if err := base.Sign(conf, req); err != nil {
		log.WithError(err).Error("签名失败")
		return nil, err
	}

	resp, err := utils.
		Request(ctx).
		SetBody(base).
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		Post(Addr)

	if err != nil {
		log.WithError(err).Error("请求失败")
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.WithError(err).Error("返回状态吗有误", zap.Int("statusCode", resp.StatusCode()))
		return nil, errors.New("返回状态码有误")
	}

	msg := &QueryLoanResponse{}
	result, err := UnmarshalResult(conf.AesKey, resp.Body())
	if err != nil {
		log.WithError(err).Error("解析返回数据失败")
		return nil, err
	}

	if err = json.Unmarshal([]byte(result.Data), msg); err != nil {

		log.WithError(err).Error("解析返回数据失败")
		return nil, err
	}

	return msg.LoanList[0], nil
}

// Query
// @param ctx
// @param conf
// @param orderId
// @date 2022-10-10 18:19:10
func Query(ctx context.Context, conf *Config, orderId string) (*QueryResult, error) {
	var result = new(QueryResult)

	statusResult, err := QueryCredit(ctx, conf, orderId)
	if err != nil {
		return nil, err
	}

	if statusResult.BankApplyStatus != "Y" {
		switch statusResult.BankApplyStatus {
		case "N":
			result.Status = 6
			result.StatusText = "审批拒绝"
		case "D":
			result.Status = 0
			result.StatusText = "待审批"
		case "E":
			result.Status = 0
			result.StatusText = "审批中"
		case "P":
			result.Status = 0
			result.StatusText = "进件未完成"
		}
		return result, nil
	}

	if statusResult.CreditAmount != "" {
		result.CreditAmount, _ = decimal.NewFromString(statusResult.CreditAmount)
	}
	result.CreditTime = statusResult.AmountTime

	loanResult, err := QueryLoan(ctx, conf, orderId)
	if err != nil {
		return nil, err
	}

	var listItem = loanResult.List
	if len(listItem) <= 0 {
		result.StatusText = "授信成功"
		return result, nil
	}

	result.Money = listItem[0].GetMoney
	result.MoneyTime = listItem[0].GetMoneyTime
	result.Status = 1
	result.StatusText = "放款成功"
	return result, nil
}
