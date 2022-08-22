package transaction

import (
	"fmt"

	"github.com/lucsky/cuid"
)

type Job string

const (
	CHECK   Job = "CHECK"
	CAPTURE Job = "CAPTURE"
	AUTH    Job = "AUTH"
	SAUTH   Job = "SAUTH"
	CANCEL  Job = "CANCEL"
	SALES   Job = "SALES"
	VOID    Job = "VOID"
)

func ConvertToJob(jobCd string) (Job, error) {
	switch jobCd {
	case "CHECK":
		return CHECK, nil
	case "CAPTURE":
		return CAPTURE, nil
	case "AUTH":
		return AUTH, nil
	case "SAUTH":
		return SAUTH, nil
	case "CANCEL":
		return CANCEL, nil
	case "SALES":
		return SALES, nil
	case "VOID":
		return VOID, nil
	}
	return "", fmt.Errorf("unknown JobCd: %s", jobCd)
}

func GenerateOrderId() string {
	return cuid.New()
}
