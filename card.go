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
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type CardInput struct {
	No           string
	Expire       string
	SecurityCode string
	Holder       string
}

type Card struct {
	Seq     int
	No      string
	Forward string
}

func CreateCard(ctx context.Context, memberId string, cardInput *CardInput) (*Card, error) {
	token, err := generateToken(cardInput)
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"MemberID": {memberId},
		"Token":    {*token},
	}
	result, err := SaveCard.Call(&values)
	if err != nil {
		return nil, err
	}
	card := &Card{}
	card.parse(result)
	return card, nil
}

func generateToken(cardInput *CardInput) (*string, error) {
	encrypted, err := encrypt(cardInput)
	if err != nil {
		return nil, err
	}
	token, err := getToken(encrypted)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (c *Card) parse(body map[string]string) error {
	var err error
	if c.Seq, err = strconv.Atoi(body["CardSeq"]); err != nil {
		return err
	}
	c.No = body["CardNo"]
	c.Forward = body["Forward"]
	return nil
}

func encrypt(input *CardInput) (string, error) {
	block, _ := pem.Decode([]byte(strings.Replace(os.Getenv("SITE_PUBLIC_KEY"), `\n`, "\n", -1)))
	if block == nil {
		log.Fatalln("block is nil")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParsePKIXPublicKey returns error: %v", err)
	}
	log.Println(reflect.TypeOf(pub))
	rpub, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("key is not rsa.PublicKey type")
	}
	card := map[string]string{
		"cardNo":       input.No,
		"expire":       input.Expire,
		"securityCode": input.SecurityCode,
		"holderName":   input.Holder,
	}
	j, err := json.Marshal(card)
	if err != nil {
		log.Fatalf("json.Marshal returns error: %v", err)
	}
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rpub, j)
	if err != nil {
		log.Fatalf("rsa.EncryptPKCS1v15 returns error: %v", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func getToken(encrypted string) (*string, error) {
	endpoint := fmt.Sprintf("https://%s/ext/api/credit/getToken", os.Getenv("SITE_DOMAIN"))
	log.Printf("endpoint: %s", endpoint)
	values := url.Values{
		"Encrypted": {encrypted},
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
