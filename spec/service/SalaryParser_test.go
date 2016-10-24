package service

import (
	. "github.com/moazzamk/moz-tech/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSalaryParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MozTech Suite")
}

var _ = Describe("SalaryParser", func () {
	Context("somthing", func () {
		It("ss", func () {
			salaryParser := SalaryParser{}
			rs := salaryParser.Parse("95K")

			Expect(float64(rs.Salary)).To(Equal(float64(95000)))
		})
	})
})

