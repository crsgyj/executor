package executor

import (
  "log"
  "testing"
)

func TestNew(t *testing.T) {
  e, err := New([]Command{
    Command{
      Name:    "go env",
      Code:    "go env",
      Desc:    "go env",
      Session: Sessions.Local(),
    },
  }, true)
  if err != nil {
    t.Fail()
    return
  }
  if err := e.Run(); err != nil {
    t.Fail()
    return
  }
}

func TestLocal(t *testing.T) {
  e, err := Local().Default([]CmdSample{
    CmdSample{
      Name:    "go help build",
      Code:    "go help build",
      Desc:    "go help build",
      Logging: false,
    },
  })
  if err != nil {
    t.Error(err)
    return
  }
  if err := e.Run(); err != nil {
    t.Error(err)
    return
  }
}

func TestRemote(t *testing.T) {
  s, err := Remote(SessionConfig{
    User:     "",
    Password: "",
    Host:     "",
    Port:     0,
    Timeout:  0,
  }).Default([]CmdSample{
    CmdSample{
      Name: "ls",
      Code: "ls",
    },
  })
  if err != nil {
    t.Error(err)
    return
  }
  if err := s.Run(); err == nil {
    log.Print(err)
    t.Fail()
    return
  }
}
