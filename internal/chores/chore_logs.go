package chores

import (
	"fmt"
	"github.com/marvin-automator/marvin/internal/db"
	"strings"
	"time"
)

const ErrorLog = "error"
const InfoLog = "info"

const logsStoreName = "choreLogs"

type ChoreLog struct {
	Id	string		`json:"id"`
	Time time.Time	`json:"-"`
	TimeStr string	`json:"time"`
	Type string		`json:"type"`
	Message string	`json:"message"`
}

func (c *Chore) Log(logType, message string, args ...interface{}) error {
	t := time.Now()
	id := c.makeLogId(t)

	log := ChoreLog{
		Id: id,
		Time: t,
		TimeStr: t.Format(time.RFC3339),
		Type: logType,
		Message: fmt.Sprintf(message, args...),
	}

	s := db.GetStore(logsStoreName)
	return s.SetWithExpiration(id, log, 7*24*time.Hour)
}

func (c *Chore) makeLogId(t time.Time) string {
	return strings.Join([]string{c.Id, t.Format(time.RFC3339Nano)}, "|")
}

// GetLogsUpTo returns the n latest logs befofe a specific time t.
func (c *Chore) GetLogsUpTo(t time.Time, n int) ([]ChoreLog, error) {
	s := db.GetStore(logsStoreName)

	loaded := 0
	result := make([]ChoreLog, n)
	var cl ChoreLog
	err := s.EachKeyBefore(c.makeLogId(t), &cl, func(key string) error {
		loaded += 1
		result[n - loaded] = cl

		if loaded == n {
			return db.StopIterating
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if loaded == 0 {
		return []ChoreLog{}, nil
	}

	if loaded < n {
		result = result[n-loaded:n-1]
	}

	return result, nil
}
