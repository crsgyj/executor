package opst

import (
	"fmt"
	"log"
	"time"
)

type TaskContext interface {
	SetState(key string, value interface{})
	GetState(key string) interface{}
	Abandon()
}

type Task struct {
	// 任务名
	Name string
	// 命令
	Command string
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
	Init func(t TaskContext)
	// 完成钩子
	Done func(t TaskContext)
	// 是否打印log
	Logging bool
	// 延迟执行，单位毫秒， 100毫秒以下设置无效
	Delay int // Delay - run after ${delay} millisecond, ignore when < 100
	// peyload
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

// SetState - set payload state
func (t *Task) SetState(key string, value interface{}) {
	payload := *(t.payload)
	payload[key] = value
}

// GetState - set payload state
func (t *Task) GetState(key string) interface{} {
	payload := *(t.payload)
	return payload[key]
}

// Abandon - abandon task
func (t *Task) Abandon() {
	t.abandon = true
}

func (t *Task) Exec() error {
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
	//fmt.Printf("[Output]=> Doing Task[%v] \"%s\", log:\n", t.index, t.Name)
	// delay task if require
	if t.Delay >= 100 {
		<-time.After(time.Duration(t.Delay) * time.Millisecond)
	}
	// create session
	if session, err = t.Session(); err != nil {
		goto ERR
	}
	// exec command
	if err = session.Run(t.Command); err != nil {
		goto ERR
	}
	goto OK

ABANDON:
	log.Printf("[Log-t%d]=> %s | abandoned.\n", t.index, t.Name)
	return nil
ERR:
	t.output = session.Output()
	t.errMsg = session.ErrMsg()
	if session == nil {
		log.Println("Session create fail.")
	} else if t.Logging {
		fmt.Printf("%s %s\n", t.output, t.errMsg)
	}
	log.Printf("[Log-t%d]=> %s | fail (%dms)\n", t.index, t.Name, time.Now().UnixNano()/1e6-beginTime)
	return err
OK:
	t.output = session.Output()
	t.errMsg = session.ErrMsg()
	if t.Logging {
		fmt.Printf("%s %s\n", t.output, t.errMsg)
	}
	log.Printf("[Log-t%d]: %s | success (%dms)\n", t.index, t.Name, time.Now().UnixNano()/1e6-beginTime)
	t.completed = true
	if t.Done != nil {
		t.Done(t)
	}
	return nil
}
