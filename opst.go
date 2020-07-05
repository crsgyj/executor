package opst

type Opst interface {
  Run() error
  Init() error
}

func New(taskList []Task, async bool) (Opst, error) {
  handler := &taskHandler{
    list:       taskList,
    resultChan: make(chan *Task),
    closeChan:  make(chan bool, 1),
    async:      async,
  }
  if err := handler.init(); err != nil {
    return nil, err
  }

  return handler, nil
}

type PreOpst interface {
  Default(commands []string) Opst
}

func Local() PreOpst {
  handler := &taskHandler{
    resultChan: make(chan *Task),
    closeChan:  make(chan bool, 1),
    config: handlerConfig{
      optsType: "local",
    },
  }

  return handler
}

func Remote(config SessionConfig) PreOpst {
  handler := &taskHandler{
    resultChan: make(chan *Task),
    closeChan:  make(chan bool, 1),
    config: handlerConfig{
      optsType:      "remote",
      sessionConfig: config,
    },
  }
  return handler
}
