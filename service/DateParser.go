package service

import (
	"strings"
	"regexp"
	"strconv"
	"time"
	"fmt"
)

type DateParser struct {

}

func (r *DateParser) Parse(str string) string {

	ret := str

	if strings.Contains(ret, `ago`) {
		re := regexp.MustCompile(`[0-9]+`)
		match := re.FindString(ret)
		sub, err := strconv.Atoi(match)
		if err != nil {
			ret = `Error parsing date ` + match
		}

		ts := time.Now()
		if strings.Contains(ret, `day`) {
			ts = ts.AddDate(0, 0, -1 * sub)

		} else if strings.Contains(ret, `week`) {
			ts = ts.AddDate(0, 0, -7 * sub)

		} else {
			ts = ts.AddDate(0, -1 * sub, 0)
		}

		ret = fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())
	}

	return ret
}
