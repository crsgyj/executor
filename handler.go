package opst

import (
  "errors"
  "fmt"
  "log"
)

type handlerConfig struct {
  optsType      string
  sessionConfig SessionConfig
}

type taskHandler struct {
  list        []Task
  Initialized bool
  async       bool
  resultChan  chan *Task
  closeChan   chan bool
  config      handlerConfig
}

func (h *taskHandler) Default(commands []string) Opst {
  return h
}

func (h *taskHandler) loopCheckingResult() {
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

func (h *taskHandler) Run() error {
  var (
    err error
  )
  if !h.Initialized {
    if err = h.Init(); err != nil {
      goto ERR
    }
  }
  log.Printf("[.Start]=> %d tasks in all.\n", len(h.list))
  go h.loopCheckingResult()
  for i := range h.list {
    task := &h.list[i]
    if task.Async {
      go func() {
        task.Exec()
        h.resultChan <- task
      }()
    } else {
      if err = task.Exec(); !task.AllowError && err != nil {
        err = errors.New("[..Stop]=> " + err.Error())
        goto ERR
      } else {
        h.resultChan <- task
      }
    }
  }
  <-h.closeChan
  log.Println("[Finish]=> All works are done.")
  return nil
ERR:
  log.Println(err.Error())
  h.closeChan <- true
  return err
}

// init taskHandler
func (h taskHandler) Init() error {
  var (
    err     error
    payload = make(map[string]interface{})
  )
  for i := range h.list {
    task := &h.list[i]
    // Abandon Empty Task
    if task.Command == "" {
      task.Abandon()
      err = fmt.Errorf("Abandon Task[%v] \"%s\" for missing \"Command\" attr.\n", string(i), task.Name)
      goto ERR
    }
    if task.Session == nil {
      task.Abandon()
      err = fmt.Errorf("Abandon Task[%v] \"%s\" for missing \"Session\" attr.\n", i, task.Name)
      goto ERR
    }
    if h.async {
      task.Async = true
    }
    if task.Async {
      task.Logging = true
    }
    task.payload = &payload
    task.index = i
  }
  h.Initialized = true
  return nil
ERR:
  return errors.New("Init error: " + err.Error())
}
