package transaction

import "fmt"

type Job string

const (
	CHECK   Job = "CHECK"
	CAPTURE Job = "CAPTURE"
	AUTH    Job = "AUTH"
	SAUTH   Job = "SAUTH"
	CANCEL  Job = "CANCEL"
	SALES   Job = "SALES"
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
	}
	return "", fmt.Errorf("unknown JobCd: %s", jobCd)
}
