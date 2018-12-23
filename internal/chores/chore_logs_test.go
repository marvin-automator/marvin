package chores

import (
	"github.com/marvin-automator/marvin/internal/db"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestChore_GetLogsUpTo(t *testing.T) {
	db.SetupTestDB()
	defer db.TearDownTestDB()

	r := require.New(t)

	c := Chore{Id: "id"}

	err := c.Log(ErrorLog, "error")
	r.NoError(err)
	err = c.Log(InfoLog, "info")
	r.NoError(err)
	err = c.Log(ErrorLog, "error2")
	r.NoError(err)
	err = c.Log(InfoLog, "info2")
	r.NoError(err)

	logs, err := c.GetLogsUpTo(time.Now(), 3)

	r.NoError(err)
	r.Equal(3, len(logs))

	expectedlogTypes := []string{InfoLog, ErrorLog, InfoLog}
	expectedmessages := []string{"info", "error2", "info2"}
	logTypes := make([]string, 0, 3)
	messages := make([]string, 0, 3)

	for _, log := range logs {
		logTypes = append(logTypes, log.Type)
		messages = append(messages, log.Message)
	}

	r.Equal(expectedlogTypes, logTypes)
	r.Equal(expectedmessages, messages)
}

func TestChore_ClearLogs(t *testing.T) {
	db.SetupTestDB()
	defer db.TearDownTestDB()

	r := require.New(t)

	c := Chore{Id: "id"}
	c.Log("", "")
	c.Log("", "")
	c.ClearLogs()

	l, err := c.GetLogsUpTo(time.Now(), 100)

	r.NoError(err)
	r.Equal(0, len(l))
}
