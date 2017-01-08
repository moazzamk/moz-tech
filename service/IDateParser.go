package service

type IDateParser interface {
	Parse(str string) string
}
