# Executor
run commands in your local or remote(SSH) session.


## Expamle
execute a serial of  commands:
```go
import "github.com/crsgyj/executor"

func main() {
  var (
    err error
    exor Excutor
    async bool          = false
  )
  commands []CmdSample := []CmdSample{
    CmdSample{
      Name: "创建helloworld容器",
      Code: "docker run -d --name helloworld hello-world",
      Logging: true,
    },
    CmdSample{
      Name: "移除容器",
      Code: "docker rm helloworld",
      Logging: true,
    },
  }
  if exor, err = excutor.Local().Default(commands, async) {
    panic(err)
  }
  
  if err = exor.Run(); err != nil {
    panic(err)
  }
}
```
or do it like this:

```go
package main

import (
	"regexp"

	"github.com/crsgyj/executor"
)

func main() {
	var (
		err   error
		ex    executor.Executor
		async bool = false
	)
	commands := []executor.Command{
		executor.Command{
			Name:       "创建helloworld容器",
			Code:       "docker run -tid --name=helloworld hello-world",
			Session:    executor.Sessions.Local(),
			AllowError: true,
			Done: func(c *executor.CmdController) {
				var (
					output      = c.GetOutput()
					containerID = ""
				)
				reg, _ := regexp.Compile("([a-z0-9]{64})")
				keys := reg.FindAllStringSubmatch(output, -1)
				if len(keys) >= 1 && len(keys[0]) >= 2 {
					containerID = keys[0][1]
				}
				c.SetState("containerID", containerID)
			},
		},
		executor.Command{
			Name:       "移除容器",
			Code:       "docker rm $container",
			Session:    executor.Sessions.Local(),
			AllowError: false,
			Logging:    true,
			Init: func(c *executor.CmdController) {
				var (
					containerID = c.GetState("containerID")
				)
				if containerID == nil || containerID == "" {
					c.Abandon()
					return
				}
				c.ReplaceCode("$container", containerID.(string), -1)
			},
		},
	}
	if ex, err = executor.New(commands, async); err != nil {
		panic(err)
	}

	if err = ex.Run(); err != nil {
		panic(err)
	}
}

```
