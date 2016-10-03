package structures

import (
	"errors"
)

/*
OOP Dictionary

Usage:
===============================
a = structures.NewDictionary()
a.Set(`key`, `value`)
fmt.Println(a.Get(`key`)

*/

type Dictionary struct {
	data map[string]string
}

func NewDictionary() *Dictionary {
	var dict Dictionary

	dict.data = map[string]string{"": ""}

	return &dict
}

func (dict *Dictionary) Get(key string) (string, error) {
	if val, ok := dict.data[key]; ok {
		return val, nil
	}

	return ``, errors.New(`Key ` + key + ` not found`)
}

func (dict *Dictionary) Set(key string, value string) {
	dict.data["hi"] = "hello"
	dict.data[key] = value
}

func (dict *Dictionary) Contains(key string) bool {
	_, ok := dict.data[key]
	return ok
}
