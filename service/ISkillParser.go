package service

import "github.com/moazzamk/moz-tech/structures"

type ISkillParser interface {
	ParseFromTags(tags *structures.UniqueSlice) *structures.UniqueSlice
	ParseFromDescription(description string) *structures.UniqueSlice
}
