package event

const (
	Min Priority = iota
	Low
	BelowNormal
	Normal
	AboveNormal
	High
	Max
)

type Priority int

func (p Priority) Value() int {
	return int(p)
}
