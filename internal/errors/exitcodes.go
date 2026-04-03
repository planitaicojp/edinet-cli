package errors

const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitValidation = 2
	ExitAuth       = 3
	ExitAPI        = 4
	ExitNetwork    = 5
	ExitNotFound   = 6
)

type ExitCoder interface {
	ExitCode() int
}

func GetExitCode(err error) int {
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}
