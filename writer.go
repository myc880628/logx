package logx

type (
	Writer interface {
		Write(Record) error
	}
)
