package language

type StringCode struct {
	Code string
}

func NewStringCode(code string) *StringCode {
	return &StringCode{code}
}
