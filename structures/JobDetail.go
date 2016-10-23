package structures


type JobDetail struct {
	Salary *SalaryRange
	Description string
	PostedDate string
	JobType []string
	Employer string
	Location string
	Skills []string
	Source string
	Title string
	Link string

	Telecommute int
	Travel int
}
