package structures

type JobDetail struct {
	Salary *SalaryRange
	Description string
	PostedDate string
	Employer string
	Location string
	Skills []string
	JobType string
	Source string
	Title string
	Link string

	Telecommute int
	Travel int
}
