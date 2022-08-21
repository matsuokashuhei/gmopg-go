package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type API string

const (
	SaveMember   API = "/payment/SaveMember.idPass"
	SearchMember API = "/payment/SearchMember.idPass"
	UpdateMember API = "/payment/UpdateMember.idPass"
	DeleteMember API = "/payment/DeleteMember.idPass"
)
const (
	SaveCard   API = "/payment/SaveCard.idPass"
	SearchCard API = "/payment/SearchCard.idPass"
	DeleteCard API = "/payment/DeleteCard.idPass"
)
const (
	EntryTran   API = "/payment/EntryTran.idPass"
	ExecTran    API = "/payment/ExecTran.idPass"
	SearchTrade API = "/payment/SearchTrade.idPass"
)

func (p API) url() *url.URL {
	url := url.URL{Scheme: "https", Path: string(p)}
	url.Host = os.Getenv("SITE_DOMAIN")
	return &url
}

func (p API) setParameters(values *url.Values) {
	values.Set("SiteID", os.Getenv("SITE_ID"))
	values.Set("SitePass", os.Getenv("SITE_PASS"))
	switch p {
	case EntryTran, ExecTran, SearchTrade:
		values.Set("ShopID", os.Getenv("SHOP_ID"))
		values.Set("ShopPass", os.Getenv("SHOP_PASS"))
	}
}

func (p API) Call(values *url.Values) ([]map[string]*string, error) {
	p.setParameters(values)
	url := p.url()
	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=windows-31j")
	log.Printf("url: %s, body: %s", url.String(), values.Encode())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	log.Printf("StatusCode: %d, Body: %s", res.StatusCode, string(body))
	result := parse(&body)
	if IsError(result) {
		return nil, NewError(result)
	}
	return result, nil
}

func parse(body *[]byte) []map[string]*string {
	params := strings.Split(string(*body), "&")
	if len(params) == 0 {
		return nil
	}
	c := countRow(params[0])
	result := make([]map[string]*string, c)
	for i := 0; i < c; i++ {
		row := make(map[string]*string)
		for _, param := range params {
			kv := strings.Split(param, "=")
			key := kv[0]
			values := kv[1]
			value := strings.Split(values, "|")[i]
			row[key] = &value
		}
		result[i] = row
	}
	return result
}

func countRow(param string) int {
	return len(strings.Split(strings.Split(param, "=")[1], "|"))
}
