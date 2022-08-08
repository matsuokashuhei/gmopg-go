package gmopg

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
)

func IsError(body map[string]string) bool {
	_, exist := body["ErrCode"]
	return exist
}

func NewError(body map[string]string) error {
	var result error
	for _, info := range strings.Split(body["ErrInfo"], "|") {
		result = multierror.Append(result, errors.New(info))
	}
	return result
}
