package gmopg

import (
	"context"
	"matsuokashuhei/gmopg-go/api"
	"matsuokashuhei/gmopg-go/transaction"
	"net/url"
	"strconv"
	"time"
)

type Transaction struct {
	OrderId     string
	Status      string
	ProcessDate time.Time
	JobCd       transaction.Job
	AccessID    string
	AccessPass  string
	ItemCode    string
	Amount      int
	Tax         int
	MemberId    string
	CardNo      string
	Expire      string
	Method      string
	PayTimes    int
	Forward     string
	TranId      string
	Approve     string
}

func FindTransaction(ctx context.Context, id string) (*Transaction, error) {
	t := &Transaction{OrderId: id}
	if err := t.load(ctx); err != nil {
		return nil, err
	}
	return t, nil
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
	return FindTransaction(ctx, *result2[0]["OrderID"])
}

func (t *Transaction) Cancel(ctx context.Context) error {
	values := url.Values{
		"AccessID":   {t.AccessID},
		"AccessPass": {t.AccessPass},
		"JobCd":      {string(transaction.CANCEL)},
	}
	if _, err := api.AlterTran.Call(&values); err != nil {
		return err
	}
	if err := t.load(ctx); err != nil {
		return err
	}
	return nil
}

func beginTransaction(ctx context.Context, orderId string, job transaction.Job, amount int, tax int) (map[string]*string, error) {
	if len(orderId) == 0 {
		orderId = transaction.GenerateOrderId()
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

func (t *Transaction) load(ctx context.Context) error {
	values := url.Values{"OrderID": {t.OrderId}}
	result, err := api.SearchTrade.Call(&values)
	if err != nil {
		return err
	}
	t.parse(result[0])
	return nil
}

func (t *Transaction) parse(result map[string]*string) error {
	var err error
	t.OrderId = *result["OrderID"]
	t.Status = *result["Status"]
	if t.ProcessDate, err = time.Parse("20060102030405", *result["ProcessDate"]); err != nil {
		return err
	}
	if t.JobCd, err = transaction.ConvertToJob(*result["JobCd"]); err != nil {
		return err
	}
	t.AccessID = *result["AccessID"]
	t.AccessPass = *result["AccessPass"]
	t.ItemCode = *result["ItemCode"]
	if t.Amount, err = strconv.Atoi(*result["Amount"]); err != nil {
		return err
	}
	if t.Tax, err = strconv.Atoi(*result["Tax"]); err != nil {
		return err
	}
	t.CardNo = *result["CardNo"]
	t.Expire = *result["Expire"]
	t.Method = *result["Method"]
	v, exist := result["PayTimes"]
	if exist && len(*v) > 0 {
		if t.PayTimes, err = strconv.Atoi(*v); err != nil {
			return err
		}
	}
	t.Forward = *result["Forward"]
	t.TranId = *result["TranID"]
	t.Approve = *result["Approve"]
	return nil
}
