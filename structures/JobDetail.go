package structures

type JobDetail struct {
	Salary *SalaryRange
	Description string
	PostedDate string
	Employer string
	Location string
	Skills []string
	JobType string
	Link string

	Telecommute int
	Travel int

}
