package executor

import (
  "fmt"
  "log"
  "strings"
  "time"
)

type CmdController struct {
  cmd *Command
}

type ProcessHandler = func(m *CmdController)

// Command 命令
type Command struct {
  // 任务名
  Name string
  // 命令
  Code string
  // 描述
  Desc string
  // 执行环境
  Session func() (execSession, error)
  // 是否允许错误 - 异步执行不关心错误
  AllowError bool
  // 是否异步执行
  Async bool
  // 是否完成
  completed bool
  // 初始化钩子
  Init ProcessHandler
  // 完成钩子
  Done ProcessHandler
  // 是否打印log
  Logging bool
  // 延迟执行，单位毫秒， 100毫秒以下设置无效
  Delay int // Delay - run after ${delay} millisecond, ignore when < 100
  // payload
  payload *map[string]interface{}
  // 是否放弃执行
  abandon bool
  // 输出内容
  output string
  // 错误内容
  errMsg string
  // 任务序号
  index int
}



func (c *CmdController) SetDelay(n int) {
  if c.cmd.completed {
    return
  }
  c.cmd.Delay = n
}

func (c *CmdController) SetCode(code string) {
  if c.cmd.completed {
    fmt.Println("Can not SetCode, command is completed.")
    return
  }
  c.cmd.Code = code
}

// ReplaceCode
func (c *CmdController) ReplaceCode(old string, new string, n int) {
  if c.cmd.completed {
    fmt.Println("Can not ReplaceCode, command is completed.")
    return
  }
  c.cmd.Code = strings.Replace(c.cmd.Code, old, new, n)
}

// SetState
func (c *CmdController) SetState(key string, value interface{}) {
  payload := *(c.cmd.payload)
  payload[key] = value
}

// GetState
func (c *CmdController) GetState(key string) interface{} {
  payload := *(c.cmd.payload)
  return payload[key]
}

// Abandon
func (c *CmdController) Abandon() {
  if c.cmd.completed {
    fmt.Println("Can not Abandon, command is completed.")
    return
  }
  c.cmd.abandon = true
}
// GetOutput - get command Output
func (c *CmdController) GetOutput() string {
  if !c.cmd.completed {
    return ""
  }
  return c.cmd.output + c.cmd.errMsg
}

// Inspect - inspect if command is valid
func (c *Command) Inspect() (err error) {
  if c.Code == "" {
    c.abandon = true
    err = fmt.Errorf("Abandon Command \"%s\" for missing \"Code\" attr.\n", c.Name)
    return
  }
  if c.Session == nil {
    c.abandon = true
    err = fmt.Errorf("Abandon Command \"%s\" for missing \"Session\" attr.\n", c.Name)
    return
  }
  return
}

func (c *Command) Exec() error {
  var (
    err       error
    session   execSession
    beginTime = time.Now().UnixNano() / 1e6
  )

  // initialize
  if c.Init != nil {
    log.Printf("[INIT](%d): %s\n", c.index, c.Name)
    c.Init(&CmdController{c})
  }
  // abandon
  if c.abandon {
    goto ABANDON
  } else {
    log.Printf("[START](%d): %s, code: \"%s\"\n", c.index, c.Name, c.Code)
  }
  // delay task if require
  if c.Delay >= 100 {
    <-time.After(time.Duration(c.Delay) * time.Millisecond)
  }
  // create session
  if session, err = c.Session(); err != nil {
    goto ERR
  }
  // exec command
  if err = session.Run(c.Code); err != nil {
    goto ERR
  }
  goto OK

ABANDON:
  log.Printf("[DROP](%d): %s", c.index, c.Name)
  return nil
ERR:
  if session == nil {
    log.Printf("[LOG](%d): Session create fail.%s\n", c.index, err.Error())
  } else {
    c.errMsg = session.ErrMsg()
    c.output = session.Output()
  }
  if c.Logging {
    fmt.Printf("[LOG](%d): %s %s\n", c.index, c.output, c.errMsg)
  }
  if c.AllowError {
    if c.Done != nil {
      log.Printf("[DONE](%d): %s\n", c.index, c.Name)
      c.Done(&CmdController{c})
    }
  }
  log.Printf("[END](%d): %s  X (%dms)\n", c.index, c.Name, time.Now().UnixNano()/1e6-beginTime)
  return err
OK:
  c.output = session.Output()
  c.errMsg = session.ErrMsg()
  if c.Logging {
    fmt.Printf("[LOG](%d): %s %s\n", c.index, c.output, c.errMsg)
  }
  log.Printf("[END](%d): %s  √ (%dms)\n", c.index, c.Name, time.Now().UnixNano()/1e6-beginTime)
  c.completed = true
  if c.Done != nil {
    c.Done(&CmdController{c})
  }
  return nil
}
