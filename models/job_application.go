package models

type JobApplication struct {
	ID          int    `json:"id" db:"id"`
	FullName    string `json:"full_name" db:"full_name"`
	Age         int    `json:"age" db:"age"`
	Degree      string `json:"degree" db:"degree"`
	Experience  string `json:"experience" db:"experience"`
	ContactTime string `json:"contact_time" db:"contact_time"`
	Resume      string `json:"resume" db:"resume"`
}
