# Executor
 A go package for running a series of local or remote(SSH) commands.


## Expamle
execute a serial of  commands:
```go
package main

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
or do it deeply:

```go
import "github.com/crsgyj/excutor"

func main() {
  var (
    err error
    exor Excutor
    commands []Command
    async bool          = false
  )
  commands = = []Command{
       Command{
         Name: "Hello_Docker",
         Code: "docker run -d --name helloworld helloworld",
         Logging: true,
         done: func(ctx ) {
        
         }
       },
       Command{
         Name: "Hello_Docker",
         Code: "docker rm helloworld",
         
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
