package gmopg

import "errors"

func IsError(body map[string]string) bool {
	_, exist := body["ErrCode"]
	return exist
}

func NewError(body map[string]string) error {
	return errors.New(body["ErrInfo"])
}
