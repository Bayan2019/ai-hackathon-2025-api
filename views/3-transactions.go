package views

import "time"

type Transaction struct {
	FullNameClient string    `json:"full_name_of_sender"`
	Transdatetime  time.Time `json:"transdatetime"`
	Amount         int64     `json:"amount"`
	Destination    string    `json:"destination"`
	Target         string    `json:"target"`
}
