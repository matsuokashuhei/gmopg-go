package gmopg

import (
	"fmt"
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

func (p API) url() *url.URL {
	url := url.URL{Scheme: "https", Path: string(p)}
	switch p {
	case SaveMember, SearchMember, UpdateMember, DeleteMember:
		url.Host = os.Getenv("SITE_DOMAIN")
	}
	return &url
}

func (p API) Call(values *url.Values) (map[string]string, error) {
	values.Set("SiteID", os.Getenv("SITE_ID"))
	values.Set("SitePass", os.Getenv("SITE_PASS"))
	url := p.url()
	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=windows-31j")
	log.Printf("url: %s", url.String())
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

func parse(body *[]byte) map[string]string {
	params := strings.Split(string(*body), "&")
	fmt.Println(params)
	result := make(map[string]string)
	for _, param := range params {
		kv := strings.Split(param, "=")
		result[kv[0]] = kv[1]
	}
	return result
}
