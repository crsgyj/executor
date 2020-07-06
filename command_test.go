package executor

import (
  "testing"
)

func TestCommand_Abandon(t *testing.T) {
  task := Command{
    Name:    "go env",
    Code:    "go env",
    Session: Sessions.Local(),
  }
  task.Abandon()
  if !task.abandon {
    t.Fail()
  }
}

func TestCommand_SetState(t *testing.T) {
  var (
    payload = make(map[string]interface{})
    task    = Command{
      Name:       "TEST setState",
      Code:       "command1",
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
