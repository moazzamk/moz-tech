package moz_tech

import (
	"github.com/moazzamk/moz-tech/structures"
	"io/ioutil"
	"strings"
)

func NewAppConfig(file string) *structures.Dictionary {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	content := string(data)
	appConfig := structures.NewDictionary()
	lines := strings.Split(content, "\n")
	for i := range lines {
		tmp := strings.Split(lines[i], "=")
		if len(tmp) != 2 {
			continue
		}
		appConfig.Set(tmp[0], tmp[1])
		//ta[]
	}

	return appConfig
}
