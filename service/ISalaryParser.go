package service

import "github.com/moazzamk/moz-tech/structures"

type ISalaryParser interface {
	Parse(str string) *structures.SalaryRange
}
