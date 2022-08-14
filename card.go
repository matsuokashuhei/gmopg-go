package gmopg

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"matsuokashuhei/gmopg-go/api"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Seq         int
	Name        string
	No          string
	HolderName  string
	Expire      string
	DefaultFlag int
	DeleteFlag  int
}

const (
	Default   int = 1
	Undecided int = 0
)

func CreateCard(ctx context.Context, memberId string, holderName string, number string, expiryDate string, securityCode string) (*Card, error) {
	token, err := generateToken(&holderName, &number, &expiryDate, &securityCode)
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"MemberID": {memberId},
		"Token":    {*token},
	}
	result, err := api.SaveCard.Call(&values)
	if err != nil {
		return nil, err
	}
	values = url.Values{
		"MemberID": {memberId},
		"CardSeq":  {*result[0]["CardSeq"]},
	}
	result, err = api.SearchCard.Call(&values)
	if err != nil {
		return nil, err
	}
	card := &Card{}
	card.parse(result[0])
	return card, nil
}

func FindCard(ctx context.Context, memberId string, seq int) (*Card, error) {
	values := url.Values{
		"MemberID": {memberId},
		"CardSeq":  {strconv.Itoa(seq)},
	}
	result, err := api.SearchCard.Call(&values)
	if err != nil {
		return nil, err
	}
	log.Printf("res: %v", result)
	card := &Card{}
	card.parse(result[0])
	return card, nil
}

func FindCards(ctx context.Context, memberId string) ([]*Card, error) {
	values := url.Values{
		"MemberID": {memberId},
	}
	result, err := api.SearchCard.Call(&values)
	if err != nil {
		return nil, err
	}
	log.Printf("res: %v", result)
	cards := make([]*Card, len(result))
	for i, row := range result {
		card := &Card{}
		card.parse(row)
		cards[i] = card
	}
	return cards, nil
}

func DeleteCard(ctx context.Context, memberId string, seq int) error {
	values := url.Values{
		"MemberID": {memberId},
		"CardSeq":  {strconv.Itoa(seq)},
	}
	if _, err := api.DeleteCard.Call(&values); err != nil {
		return err
	}
	return nil
}

func generateToken(holderName *string, number *string, expiryDate *string, securityCode *string) (*string, error) {
	encrypted, err := encrypt(holderName, number, expiryDate, securityCode)
	if err != nil {
		return nil, err
	}
	token, err := getToken(encrypted)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (c *Card) parse(body map[string]*string) error {
	seq, err := strconv.Atoi(*body["CardSeq"])
	if err != nil {
		return err
	}
	c.Seq = seq
	c.No = *body["CardNo"]
	v, exist := body["CardName"]
	if exist {
		c.Name = *v
	}
	v, exist = body["DefaultFlg"]
	if exist {
		defaultFlag, err := strconv.Atoi(*v)
		if err != nil {
			return err
		}
		c.DefaultFlag = defaultFlag
	}
	v, exist = body["Expire"]
	if exist {
		c.Expire = *v
	}
	v, exist = body["HolderName"]
	if exist {
		c.HolderName = *v
	}
	v, exist = body["DeleteFlag"]
	if exist {
		deleteFlag, err := strconv.Atoi(*v)
		if err != nil {
			return err
		}
		c.DeleteFlag = deleteFlag
	}
	return nil
}

func encrypt(holder *string, number *string, expiryDate *string, securityCode *string) (*string, error) {
	block, _ := pem.Decode([]byte(strings.Replace(os.Getenv("SITE_PUBLIC_KEY"), `\n`, "\n", -1)))
	if block == nil {
		log.Fatalln("block is nil")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParsePKIXPublicKey returns error: %v", err)
	}
	rpub, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("key is not rsa.PublicKey type")
	}
	card := map[string]*string{
		"holderName":   holder,
		"cardNo":       number,
		"expire":       expiryDate,
		"securityCode": securityCode,
	}
	j, err := json.Marshal(card)
	if err != nil {
		log.Fatalf("json.Marshal returns error: %v", err)
	}
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rpub, j)
	if err != nil {
		log.Fatalf("rsa.EncryptPKCS1v15 returns error: %v", err)
	}
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	return &encoded, nil
}

func getToken(encrypted *string) (*string, error) {
	endpoint := fmt.Sprintf("https://%s/ext/api/credit/getToken", os.Getenv("SITE_DOMAIN"))
	log.Printf("endpoint: %s", endpoint)
	values := url.Values{
		"Encrypted": {*encrypted},
		"ShopID":    {os.Getenv("SHOP_ID")},
		"KeyHash":   {os.Getenv("SITE_PUBLIC_KEY_HASH")},
	}
	resp, err := http.PostForm(endpoint, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	resultCode := result["resultCode"].([]interface{})[0].(string)
	if resultCode != "000" {
		return nil, errors.New(resultCode)
	}
	token := result["tokenObject"].(map[string]interface{})["token"].([]interface{})[0].(string)
	return &token, nil
}
