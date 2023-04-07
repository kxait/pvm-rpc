package main

import (
	"fmt"
	"os"
	"pvm_rpc/pvm"
	"pvm_rpc/pvmrpc"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic(":(")
	}
	pvm.Initsend(pvm.DataDefault)

	mode := os.Args[1]
	if mode == "client" {
		hostname, _ := os.Hostname()

		fmt.Printf("hello from child!\n")
		parentId, err := pvm.Parent()
		pvm.Initsend(pvm.DataDefault)

		if err != nil {
			panic(err)
		}

		target := pvmrpc.Target{TaskId: parentId}

		for {
			time.Sleep(1 * time.Millisecond)

			result := <-target.Call("ping", hostname)

			if result.Err != nil {
				panic(result.Err)
			}

			fmt.Printf("server responded: %s\n", result.Response.Content)
		}

	} else if mode == "server" {
		pvm.CatchoutStdout()
		result, err := pvm.Spawn("pvm_rpc", []string{"client"}, pvm.TaskDefault, "", 3)
		if err != nil {
			panic(err)
		}
		defer pvm.Kill(result.TIds[0])

		fmt.Printf("child spawned with id %d\n", result.TIds[0])

		server := pvmrpc.RpcServer{Handlers: make(map[pvmrpc.MessageType]pvmrpc.RpcHandler)}
		server.Handlers["ping"] = func(m *pvmrpc.Message) (*pvmrpc.Message, error) {
			fmt.Printf("%s\n", m.Content)
			return m.CreateResponse("pong"), nil
		}

		fmt.Printf("registered handler\n")

		var erro error
		for erro == nil {
			time.Sleep(1 * time.Millisecond)
			erro = server.StepEventLoop()
		}
		fmt.Printf("server stopped working: %s\n", erro)
	} else {
	}
}
