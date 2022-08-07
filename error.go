package gmopg

import (
	"errors"
	"regexp"
	"strings"
)

func IsError(body []byte) bool {
	s := string(body)
	m, _ := regexp.MatchString(`^ErrCode=(E[0-9]{2}|E[0-9]{2}\|)+&ErrInfo=(E[0-9]{8}|E[0-9]{8}\|)+$`, s)
	return m
}

func NewError(body []byte) error {
	s := string(body)
	params := strings.Split(s, "&")
	info := strings.Split(params[len(params)-1], "=")
	return errors.New(info[len(info)-1])
}
