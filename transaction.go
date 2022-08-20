package gmopg

import (
	"context"
	"matsuokashuhei/gmopg-go/api"
	"matsuokashuhei/gmopg-go/transaction"
	"net/url"
	"strconv"
	"time"

	"github.com/lucsky/cuid"
)

type Transaction struct {
	OrderId     string
	Forward     string
	Method      string
	PayTimes    int
	Approve     string
	TranId      string
	TranDate    time.Time
	CheckString string
}

func CreateTransaction(ctx context.Context, memberId string, cardSeq int, orderId string, job transaction.Job, amount int, tax int) (*Transaction, error) {
	result1, err := beginTransaction(ctx, orderId, job, amount, tax)
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"AccessID":   {*result1["AccessID"]},
		"AccessPass": {*result1["AccessPass"]},
		"OrderID":    {orderId},
		"JobCd":      {string(job)},
		"Amount":     {strconv.Itoa(amount)},
		"Tax":        {strconv.Itoa(tax)},
		"Method":     {"1"},
		"MemberID":   {memberId},
		"CardSeq":    {strconv.Itoa(cardSeq)},
	}
	result2, err := api.ExecTran.Call(&values)
	if err != nil {
		return nil, err
	}
	t := &Transaction{}
	if err := t.parse(result2[0]); err != nil {
		return nil, err
	}
	return t, nil
}

func beginTransaction(ctx context.Context, orderId string, job transaction.Job, amount int, tax int) (map[string]*string, error) {
	if len(orderId) == 0 {
		orderId = generateOrderId()
	}
	values := url.Values{
		"OrderID": {orderId},
		"JobCd":   {string(job)},
		"Amount":  {strconv.Itoa(amount)},
		"Tax":     {strconv.Itoa(tax)},
	}
	result, err := api.EntryTran.Call(&values)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

// func createTransaction(ctx context.Context, memberId string, cardSeq int, orderId string, job Job, amount int, tax int) (map[string]*string, error) {
// }

func generateOrderId() string {
	return cuid.New()
}

func (t *Transaction) parse(result map[string]*string) error {
	var err error
	t.OrderId = *result["OrderID"]
	t.Forward = *result["Forward"]
	t.Method = *result["Method"]
	v, exist := result["PayTimes"]
	if exist && len(*v) > 0 {
		if t.PayTimes, err = strconv.Atoi(*v); err != nil {
			return err
		}
	}
	t.Approve = *result["Approve"]
	t.TranId = *result["TranID"]
	return nil
}
