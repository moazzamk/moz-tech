package service

import (
	. "github.com/moazzamk/moz-tech/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var _ = Describe("Storage", func () {
	Context("stores stuff", func () {
		It("something", func () {
			salaryParser := SalaryParser{}
			rs := salaryParser.Parse("95K")

			Expect(float64(rs.Salary)).To(Equal(float64(95000)))
		})

		It("parses salaries with commas", func () {
			salaryParser := SalaryParser{}
			rs := salaryParser.Parse("95,000 ")

			Expect(float64(rs.Salary)).To(Equal(float64(95000)))
		})
	})
})




