package executor

import (
	"fmt"
	"log"
	"time"
)

type commandContext interface {
	SetState(key string, value interface{})
	GetState(key string) interface{}
	Abandon()
}

type Command struct {
	// 任务名
	Name string
	// 命令
	Code string
	// 描述
	Desc string
	// 执行环境
	Session func() (Session, error)
	// 是否允许错误 - 异步执行不关心错误
	AllowError bool
	// 是否异步执行
	Async bool
	// 是否完成
	completed bool
	// 初始化钩子
	Init func(t commandContext)
	// 完成钩子
	Done func(t commandContext)
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
// Check
func (t *Command) Inspect() (err error) {
	if t.Code == "" {
		t.Abandon()
		err = fmt.Errorf("Abandon Command \"%s\" for missing \"Code\" attr.\n",  t.Name)
		return
	}
	if t.Session == nil {
		t.Abandon()
		err = fmt.Errorf("Abandon Command \"%s\" for missing \"Session\" attr.\n", t.Name)
		return
	}
	return
}


// SetState - set payload state
func (t *Command) SetState(key string, value interface{}) {
	payload := *(t.payload)
	payload[key] = value
}

// GetState - set payload state
func (t *Command) GetState(key string) interface{} {
	payload := *(t.payload)
	return payload[key]
}

// Abandon - abandon task
func (t *Command) Abandon() {
	t.abandon = true
}

func (t *Command) Exec() error {
	var (
		err       error
		session   Session
		beginTime = time.Now().UnixNano() / 1e6
	)
	// initialize
	if t.Init != nil {
		t.Init(t)
	}
	// abandon
	if t.abandon {
		goto ABANDON
	}
	log.Printf("[Executing%v]: %s, code: \"%s\"\n", t.index, t.Name, t.Code)
	// delay task if require
	if t.Delay >= 100 {
		<-time.After(time.Duration(t.Delay) * time.Millisecond)
	}
	// create session
	if session, err = t.Session(); err != nil {
		goto ERR
	}
	// exec command
	if err = session.Run(t.Code); err != nil {
		goto ERR
	}
	goto OK

ABANDON:
	log.Printf("[Abandon%d]: %s, code: \"%s\".\n", t.index, t.Name, t.Code)
	return nil
ERR:

	if session == nil {
		fmt.Println("Session create fail.", err.Error())
	} else {
		t.output = session.Output()
		t.errMsg = session.ErrMsg()
	}
	if t.Logging {
		fmt.Printf("%s %s\n", t.output, t.errMsg)
	}
	log.Printf("[Task%d]: %s  X (%dms)\n", t.index, t.Name, time.Now().UnixNano()/1e6-beginTime)
	return err
OK:
	t.output = session.Output()
	t.errMsg = session.ErrMsg()
	if t.Logging {
		fmt.Printf("%s %s\n", t.output, t.errMsg)
	}
	log.Printf("[Task%d]: %s  √ (%dms)\n", t.index, t.Name, time.Now().UnixNano()/1e6-beginTime)
	t.completed = true
	if t.Done != nil {
		t.Done(t)
	}
	return nil
}
