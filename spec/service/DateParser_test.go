package service

import (
	. "github.com/moazzamk/moz-tech/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
	"fmt"
)

func TestDateParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MozTech Suite")
}

var _ = Describe("DateParser", func () {
	Context("parse date strings", func () {
		It("parses 'day ago' dates", func() {

			ts := time.Now()
			ts = ts.AddDate(0, 0, -1 * 1)
			tsString := fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

			dateParser := DateParser{}
			Expect(dateParser.Parse("1 day ago")).To(Equal(tsString))
		})

		It("parses 'days ago' dates", func () {

			ts := time.Now()
			ts = ts.AddDate(0, 0, -1 * 2)
			tsString := fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

			dateParser := DateParser{}
			Expect(dateParser.Parse("2 days ago")).To(Equal(tsString))

		})

		It("parses 'week ago' dates", func () {
			ts := time.Now()
			ts = ts.AddDate(0, 0, -1 * 7)
			tsString := fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

			dateParser := DateParser{}
			Expect(dateParser.Parse("1 week ago")).To(Equal(tsString))
		})

		It("parses 'weeks ago' dates", func () {
			ts := time.Now()
			ts = ts.AddDate(0, 0, -1 * 14)
			tsString := fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

			dateParser := DateParser{}
			Expect(dateParser.Parse("2 week ago")).To(Equal(tsString))
		})

		It("parses 'seconds ago' dates", func () {
			ts := time.Now()
			tsString := fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

			dateParser := DateParser{}
			Expect(dateParser.Parse("2 secs ago")).To(Equal(tsString))
		})

		It("parses 'mins ago' dates", func () {
			ts := time.Now()
			tsString := fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

			dateParser := DateParser{}
			Expect(dateParser.Parse("2 mins ago")).To(Equal(tsString))
		})

		It("parses 'hours ago' dates", func () {
			ts := time.Now()
			tsString := fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

			dateParser := DateParser{}
			Expect(dateParser.Parse("2 hours ago")).To(Equal(tsString))
		})

		It("parses db dates", func () {
			dateParser := DateParser{}
			Expect(dateParser.Parse("2015-01-01")).To(Equal("2015-01-01"))
		})

		It("parses 'Oct 01, 2016' dates", func () {
			dateParser := DateParser{}
			Expect(dateParser.Parse("Oct 01, 2016")).To(Equal("2016-10-01"))
		})
	})
})
