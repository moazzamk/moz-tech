package service

import (
	. "github.com/moazzamk/moz-tech/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	//"github.com/moazzamk/moz-tech/structures"

)

func TestSkillParser(t *testing.T) {

	RegisterTestingT(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Salary Parser Suite")
}

var _ = Describe("SkillParser", func () {
	Context("parses skills from tags", func () {
		It("parses abbreviated salaries", func () {
/*			slice := structures.NewUniqueSlice([]string{
				`Python`,
				`MySQL`,
				`DJango`,
				`javascript`,
				`ios`,
				`android`,
				`aws`,
				`iot`,
			})
*/


//			skillParser := NewSkillParser()/
//			rs := skillParser.ParseFromTags("95K")


		})

		It("parses salaries with commas", func () {
			salaryParser := SalaryParser{}
			rs := salaryParser.Parse("95,000 ")

			Expect(float64(rs.Salary)).To(Equal(float64(95000)))
		})

	})
})

