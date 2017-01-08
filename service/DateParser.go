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

	ret := `1900-01-01`

	if strings.Contains(str, `yesterday`) {
		ts := time.Now()
		return ts.AddDate(0, 0, -1).Format(`2006-01-02`)
	}

	if strings.Contains(str, `ago`) {
		re := regexp.MustCompile(`[0-9]+`)
		match := re.FindString(str)
		sub, err := strconv.Atoi(match)
		if err != nil {
			ret = `Error parsing date ` + match
		}

		ts := time.Now()
		if strings.Contains(str, `day`) {
			ts = ts.AddDate(0, 0, -1 * sub)

		} else if strings.Contains(str, `week`) {
			ts = ts.AddDate(0, 0, -7 * sub)

		} else if strings.Contains(str, `sec`) || strings.Contains(str, `min`) || strings.Contains(str, `hour`) {
		} else {
			ts = ts.AddDate(0, -1 * sub, 0)
		}

		return fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())
	} else {

		formats := []string{
			`2006-01-02`,
			`Jan 02, 2006`,
		}

		for _, format := range formats {
			ts, err := time.Parse(format, str)
			if err == nil {
				return fmt.Sprintf("%d-%02d-%02d", ts.Year(), ts.Month(), ts.Day())
			}
		}
	}

	return ret
}
