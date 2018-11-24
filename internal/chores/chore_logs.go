package chores

import "time"

type ChoreLog struct {
	Id	string		`json:"id"`
	Time time.Time	`json:"time"`
	Type string		`json:"type"`
	Message string	`json:"message"`
}
