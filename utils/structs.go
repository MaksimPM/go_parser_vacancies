package utils

type Vacancy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Employer struct {
		Name string `json:"name"`
	} `json:"employer"`
	Salary struct {
		From     int    `json:"from"`
		To       int    `json:"to"`
		Currency string `json:"currency"`
	} `json:"salary"`
	Area struct {
		Name string `json:"name"`
	} `json:"area"`
	Employment struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"employment"`
	URL        string `json:"alternate_url"`
	Experience struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:experience`
	Schedule struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:schedule`
}

type VacanciesResponse struct {
	Items []Vacancy `json:"items"`
}

type VacancyDetails struct {
	Description string `json:"description"`
}

type VacancyResult struct {
	ID          string
	Description string
}
