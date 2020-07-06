package executor

import (
  "errors"
  "fmt"
  "log"
)

type commandHandler struct {
  list        []Command
  initialized bool
  async       bool
  resultChan  chan *Command
  closeChan   chan bool
}

func (h *commandHandler) loopCheckingResult() {
  var n = 0
  for {
    select {
    case <-h.resultChan:
      n++
      if n >= len(h.list) {
        h.closeChan <- true
        break
      }
    }
  }
}

func (h *commandHandler) Run() error {
  var (
    err   error
    count int = len(h.list)
  )
  if !h.initialized {
    if err = h.Init(); err != nil {
      goto ERR
    }
  }
  logStart(count)
  go h.loopCheckingResult()
  for i := range h.list {
    command := &h.list[i]
    if command.Async {
      go func() {
        command.Exec()
        h.resultChan <- command
      }()
    } else {
      if err = command.Exec(); !command.AllowError && err != nil {
        err = errors.New(fmt.Sprintf("Break up by command[%d]: %s, %s", i, command.Name, err.Error()))
        goto ERR
      } else {
        h.resultChan <- command
      }
    }
  }
  <-h.closeChan
  logDone()
  return nil
ERR:
  log.Println(err.Error())
  h.closeChan <- true
  return err
}

// init commandHandler
func (h commandHandler) Init() error {
  var (
    err     error
    payload = make(map[string]interface{})
  )
  for i := range h.list {
    c := &h.list[i]
    // Abandon Empty Task
    if err = c.Inspect(); err != nil {
      goto ERR
    }
    if h.async {
      c.Async = true
    }
    if c.Async {
      c.Logging = false
    }
    c.payload = &payload
    c.index = i
  }
  h.initialized = true
  return nil
ERR:
  return errors.New("Init error: " + err.Error())
}

func logStart(count int) {
  if count > 1 {
    log.Printf("Start. with %d commands in all.\n", count)
  } else {
    log.Printf("Start. with %d command in all.\n", count)
  }
}

func logDone() {
  log.Println("All works are done.\n --------------------------------\n ")
}