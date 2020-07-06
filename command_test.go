package executor

import (
  "testing"
)

func TestCommand_Abandon(t *testing.T) {
  cmd := &Command{
    Name:    "go env",
    Code:    "go env",
    Session: Sessions.Local(),
  }
  controller := CmdController{cmd}
  controller.Abandon()
  if !controller.cmd.abandon {
    t.Fail()
  }
}

func TestCommand_SetState(t *testing.T) {
  payload := make(map[string]interface{})

  cmd := &Command{
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
  controller := &CmdController{cmd}

  kv := map[string]interface{}{
    "keyA": "valueA",
    "keyB": struct{}{},
    "keyC": nil,
    "keyD": 1,
  }

  for k, v := range kv {
    controller.SetState(k, v)
    if getV := controller.GetState(k); getV != v {
      goto ERR
    }
    if getV := (*controller.cmd.payload)[k]; getV != v {
      goto ERR
    }
  }

  return
ERR:
  t.Fail()
}
