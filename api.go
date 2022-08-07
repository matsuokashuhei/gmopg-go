package gmopg

import (
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
)

func (p API) url() url.URL {
	url := url.URL{Scheme: "https", Path: string(p)}
	switch p {
	case SaveMember:
		url.Host = os.Getenv("SITE_DOMAIN")
	case SearchMember:
		url.Host = os.Getenv("SITE_DOMAIN")
	}
	return url
}

func (p API) Call(values *url.Values) (*http.Response, error) {
	values.Set("SiteID", os.Getenv("SITE_ID"))
	values.Set("SitePass", os.Getenv("SITE_PASS"))
	url := p.url()
	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=windows-31j")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("StatusCode: %d", res.StatusCode)
	return res, nil
}
