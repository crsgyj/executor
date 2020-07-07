# Executor
 A go package for running a series of local or remote(SSH) commands.


## Expamle
execute a serial of  commands:
```go
import "github.com/crsgyj/executor"

func main() {
  var (
    err error
    exor Excutor
    commands []CmdSample = []CmdSample{
      CmdSample{
        Name: "Hello_Docker",
        Code: "docker run -d --name helloworld helloworld",
        Logging: true,
      },
      CmdSample{
        Name: "Remove_Container",
        Code: "docker rm helloworld",
        Logging: true,
      },
    }
    async bool          = false
  )
  if exor, err = excutor.Local().Default(commands, async) {
    panic(err)
  }
  
  if err = exor.Run(); err != nil {
    panic(err)
  }
}
```
or do it in depth:

```go
import "github.com/crsgyj/executor"

func main() {
  var (
    err error
    exor Excutor
    async   bool  = false
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
  if exor, err = excutor.New(commands, async) {
    panic(err)
  }
  
  if err = exor.Run(); err != nil {
    panic(err)
  }
}
```
