package gmopg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Card struct {
	Seq     int
	No      string
	Forward string
}

func NewCard(seq int, no, string, forward string) *Card {
	return &Card{Seq: seq, No: no, Forward: forward}
}

func GenerateToken(cardNo string, expire string, securityCode string, holder string) (string, error) {
	block, _ := pem.Decode([]byte(strings.Replace(os.Getenv("SITE_PUBLIC_KEY"), `\n`, "\n", -1)))
	if block == nil {
		log.Fatalln("block is nil")
	}
	log.Println(block.Type)
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
		"cardNo":       cardNo,
		"expire":       expire,
		"securityCode": securityCode,
		"holderName":   holder,
	}
	j, err := json.Marshal(card)
	if err != nil {
		log.Fatalf("json.Marshal returns error: %v", err)
	}
	log.Println(rpub.Size())
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rpub, j)
	if err != nil {
		log.Fatalf("rsa.EncryptPKCS1v15 returns error: %v", err)
	}
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	log.Println(encoded)
	return encoded, nil
}

func (c *Card) Parse(body map[string]string) error {
	var err error
	if c.Seq, err = strconv.Atoi(body["CardSeq"]); err != nil {
		return err
	}
	c.No = body["CardNo"]
	c.Forward = body["Forward"]
	return nil
}
