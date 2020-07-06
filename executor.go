package executor

type Executor interface {
  Run() error
  Init() error
}

type ExecCreator struct {
  IsRemote bool
  session  func() (Session, error)
}

type CmdSample struct {
  Name    string
  Code    string
  Desc    string
  Logging bool
}

var Sessions = struct {
  Remote func(config SessionConfig) func() (Session, error)
  Local  func() func() (Session, error)
}{
  Remote: newRemoteSession,
  Local:  newLocalSession,
}

func New(cmdList []Command, async bool) (Executor, error) {
  h := &commandHandler{
    list:       cmdList,
    resultChan: make(chan *Command),
    closeChan:  make(chan bool, 1),
    async:      async,
  }
  if err := h.Init(); err != nil {
    return nil, err
  }

  return h, nil
}

func Local() *ExecCreator {
  creator := &ExecCreator{
    IsRemote: false,
    session:  Sessions.Local(),
  }

  return creator
}

func Remote(config SessionConfig) *ExecCreator {
  creator := &ExecCreator{
    IsRemote: true,
    session:  Sessions.Remote(config),
  }

  return creator
}

func (e *ExecCreator) Default(commandLines []CmdSample) (Executor, error) {
  var (
    list     []Command       = []Command{}
    async    bool            = e.IsRemote
    executor *commandHandler = &commandHandler{
      list:       nil,
      resultChan: make(chan *Command),
      closeChan:  make(chan bool, 1),
      async:      async,
    }
    err error
  )
  for _, sample := range commandLines {
    list = append(list, Command{
      Name:    sample.Name,
      Code:    sample.Code,
      Desc:    sample.Desc,
      Session: e.session,
      Async:   async,
      Logging: sample.Logging,
    })
  }
  executor.list = list
  if err = executor.Init(); err != nil {
    return nil, err
  }

  return executor, nil
}
