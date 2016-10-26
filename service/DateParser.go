package service

import (
	"strings"
	"regexp"
	"strconv"
	"time"
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

		} else if strings.Contains(ret, `sec`) || strings.Contains(ret, `min`) || strings.Contains(ret, `hour`) {
		} else {
			ts = ts.AddDate(0, -1 * sub, 0)
		}

		return ts.Format(`2006-01-02`)
	} else {

		formats := []string{
			`2006-01-02`,
			`Jan 02, 2006`,
		}

		for _, format := range formats {
			ts, err := time.Parse(format, str)
			if err == nil {
				return ts.Format(`2006-01-02`)
			}
		}
	}

	return ret
}
