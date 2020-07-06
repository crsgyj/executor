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
    commands []executor.Command
    async bool          = false
  )
  commands = = []executor.Command{
       executor.Command{
         Name: "Hello_Docker",
         Code: "docker run -d --name helloworld helloworld",
         Logging: true,
         done: func(c *executor.CmdController) {
            var (
              output = c.GetOutPut()
              containerID = ""      
            )
            reg, _ = regexp.Compile("sha256:(\\w+)")
            keys := reg.FindAllStringSubmatch(output, -1)
            if len(keys) >= 1 && len(keys[0]) >= 2 {
                containerID = keys[0][1]
            }
            t.payload["containerID"] = oldImageID
            return t  
         }
       },
       executor.Command{
         Name: "Remove Helloworld Container",
         Code: "docker rm $containerID",
         init: func (c *executor.CmdController) {
            var (
              containerID = c.GetState("containerID")
            )
            if containerID == nil {
               c.Abandon()
               return        
            }
            c.ReplaceCode("$containerID", containerID.(string))
            return
         },
         Logging: true,
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
