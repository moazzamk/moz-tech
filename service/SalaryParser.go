package service

import (
	"github.com/moazzamk/moz-tech/structures"
	"regexp"
	"strings"
	"strconv"
)

type SalaryParser struct {

}

func (r SalaryParser) Parse(str string) *structures.SalaryRange {
	ret := new(structures.SalaryRange)
	re := regexp.MustCompile(`[$0-9,.kK]+\s*(-|to)*\s*[$0-9,.kK]+`)
	charsToReplace := map[string]string{
		`k`:  `000`,
		`K`:  `000`,
		`,`:  ``,
		`$`:  ``,
		`to`: `-`,
		` `:  ``,
	}

	ret.OriginalSalary = str
	tmp := re.FindString(str)

	if tmp == `` {
		return ret
	}

	for j, v := range charsToReplace {
		tmp = strings.Replace(tmp, j, v, -1)
	}

	rangeArray := strings.Split(tmp, `-`)
	rangeArrayLen := len(rangeArray)
	if rangeArrayLen == 2 {
		ret.MinSalary, _ = strconv.ParseFloat(rangeArray[0], 64)
		ret.MaxSalary, _ = strconv.ParseFloat(rangeArray[1], 64)
	} else if rangeArrayLen == 1 { // Salary is not a range
		ret.Salary, _ = strconv.ParseFloat(rangeArray[0], 64)
	}

	// Calculate yearly salary if its an hourly position
	if strings.Contains(str, `hr`) || strings.Contains(str, `hour`) {
		ret.CalculatedMinYearlySalary = ret.MinSalary * 40 * 52
		ret.CalculatedMaxYearlySalary = ret.MaxSalary * 40 * 52
		ret.CalculatedSalary = ret.Salary * 40 * 52
	}


	return ret
}
