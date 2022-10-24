package jd

import (
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/zrb-channel/jd/lib"
)

type SignRequest interface {
	SetSign(sign string)
}

type (
	Config struct {
		AppSecret string
		AesKey    string
		AppId     string
	}

	BaseRequest struct {
		AppId     string `json:"appId"`
		Method    string `json:"method"`
		Param     string `json:"param"`
		Timestamp int64  `json:"timestamp"`
	}

	BaseResponse struct {
		Param     string `json:"param"`
		Timestamp int64  `json:"timestamp"`
	}
)

type BaseResult struct {
	Code    string `json:"code"`
	Data    string `json:"data"`
	Sign    string `json:"sign"`
	Message string `json:"message"`
}

type CreateOrderResponse struct {
	ProductUrl string `json:"productUrl"`
}

func UnmarshalResult(aesKey string, body []byte) (*BaseResult, error) {
	result := &BaseResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	v, err := lib.AesDecryptECBFromHex(result.Param, []byte(aesKey))
	if err != nil {
		return nil, err
	}

	base := &BaseResult{}
	return base, json.Unmarshal(v, base)
}

type CreateOrderRequest struct {
	ProductId         string `json:"productId" url:"productId"`
	OrderId           string `json:"orderId" url:"orderId"`
	Name              string `json:"name" url:"name"`
	UserTel           string `json:"userTel" url:"userTel"`
	CompanyName       string `json:"companyName" url:"companyName"`
	CompanyCreditCode string `json:"companyCreditCode,omitempty" url:"companyCreditCode,omitempty"`
	Sign              string `json:"sign,omitempty" url:"-"`
}

type (
	OrderItem struct {
		OrderID string `json:"orderId" url:"orderId"`
	}

	Order struct {
		OrderId         string `json:"orderId"`
		BankApplyStatus string `json:"bankApplyStatus"`
		CreditAmount    string `json:"creditAmount"`
		AmountTime      string `json:"amountTime"`
	}

	CreditList []*OrderItem

	QueryCreditRequest struct {
		ProductId  string     `json:"productId" url:"productId"`
		Sign       string     `json:"sign" url:"-"`
		CreditList CreditList `json:"creditList" url:"creditList"`
	}

	QueryCreditResponse struct {
		ProductId  string   `json:"productId"`
		CreditList []*Order `json:"creditList"`
	}

	LoanOrder struct {
		GetMoney         decimal.Decimal `json:"getMoney"`
		GetMoneyTime     string          `json:"getMoneyTime"`
		LoanBankOrderNum string          `json:"loanBankOrderNum"`
		Term             string          `json:"term"`
	}

	OrderLoan struct {
		OrderId string       `json:"orderId"`
		List    []*LoanOrder `json:"list"`
	}

	QueryLoanResponse struct {
		ProductId string       `json:"productId"`
		LoanList  []*OrderLoan `json:"loanList"`
	}

	QueryResult struct {
		Money        decimal.Decimal
		MoneyTime    string
		Status       int8
		StatusText   string
		CreditAmount decimal.Decimal
		CreditTime   string
	}
)

func (req *CreateOrderRequest) SetSign(v string) { req.Sign = v }

func (req *QueryCreditRequest) SetSign(v string) {
	req.Sign = v
}
