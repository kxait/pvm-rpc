package pvm_rpc

import (
	"fmt"
	"os"
	"time"

	"pvm_rpc/pvm"
)

func Demo() {
	if len(os.Args) < 2 {
		panic(":(")
	}
	_, err := pvm.Initsend(pvm.DataDefault)

	if err != nil {
		println("pvm.Initsend (pvm_initsend()) failed. this might indicate that you do not have pvm on your system or it is not running.")
		panic(err)
	}

	mode := os.Args[1]
	if mode == "client" {
		hostname, _ := os.Hostname()

		fmt.Printf("hello from child!\n")
		parentId, err := pvm.Parent()
		pvm.Initsend(pvm.DataDefault)

		if err != nil {
			panic(err)
		}

		target := Target{TaskId: parentId}

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
		result, err := pvm.Spawn("pvm-rpc", []string{"client"}, pvm.TaskDefault, "", 3)
		if err != nil {
			panic(err)
		}
		defer pvm.Kill(result.TIds[0])

		fmt.Printf("child spawned with id %d\n", result.TIds[0])

		server := RpcServer{Handlers: make(map[MessageType]RpcHandler)}
		server.Handlers["ping"] = func(m *Message) (*Message, error) {
			fmt.Printf("%s\n", m.Content)
			return m.CreateResponse("pong"), nil
		}

		fmt.Printf("registered handler\n")

		var erro error
		var hadMessage bool = false
		for erro == nil {
			if !hadMessage {
				time.Sleep(1 * time.Millisecond)
			}
			hadMessage, erro = server.StepEventLoop()
		}
		fmt.Printf("server stopped working: %s\n", erro)
	} else {
		panic(":(")
	}
}
