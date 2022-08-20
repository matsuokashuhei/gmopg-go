package transaction

type Job string

const (
	CHECK   Job = "CHECK"
	CAPTURE Job = "CAPTURE"
	AUTH    Job = "AUTH"
	SAUTH   Job = "SAUTH"
)
