package opst

import (
	"testing"
)

func TestTask_Abandon(t *testing.T) {
	task := Task{
		Name:    "go env",
		Command: "go env",
		Session: TerminalSessions.Local(),
	}
	task.Abandon()
	if !task.abandon {
		t.Fail()
	}
}

func TestTask_SetState(t *testing.T) {
	var (
		payload = make(map[string]interface{})
		task    = Task{
			Name:       "TEST setState",
			Command:    "command1",
			Session:    nil,
			payload:    &payload,
			AllowError: false,
			Async:      false,
			completed:  false,
			Init:       nil,
			Done:       nil,
		}
		kv = map[string]interface{}{
			"keyA": "valueA",
			"keyB": struct{}{},
			"keyC": nil,
			"keyD": 1,
		}
	)

	for k, v := range kv {
		task.SetState(k, v)
		if getV := task.GetState(k); getV != v {
			goto ERR
		}
	}

	return
ERR:
	t.Fail()
}
