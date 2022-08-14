package api

import (
	"errors"

	"github.com/hashicorp/go-multierror"
)

func IsError(result []map[string]*string) bool {
	_, exist := result[0]["ErrCode"]
	return exist
}

func NewError(result []map[string]*string) error {
	var err error
	for _, row := range result {
		err = multierror.Append(err, errors.New(*row["ErrInfo"]))
	}
	return err
}
