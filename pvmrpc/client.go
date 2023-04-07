package pvmrpc

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"pvm_rpc/pvm"
	"time"
)

func Sus() {}

func (t *Target) ResetSendBuffer() {
	pvm.Initsend(pvm.DataDefault)
}

func (t *Target) Kill() error {
	return pvm.Kill(t.TaskId)
}

func (t *Target) Call(msgType MessageType, content string) <-chan *ReceiveResult {
	r := make(chan *ReceiveResult)

	go (func() {
		me, err := pvm.Mytid()

		if err != nil {
			r <- &ReceiveResult{Err: err}
		}

		id := rand.Intn(math.MaxInt32)

		msg := &Message{
			Id:           id,
			CallerTaskId: me,
			Type:         msgType,
			Content:      content,
		}

		err = t.send(msg)
		if err != nil {
			r <- &ReceiveResult{Err: err}
		}

		result := <-t.receiveContinuously(msg.Id)

		if result.Err != nil {
			r <- &ReceiveResult{Err: err}
		}

		r <- result
	})()

	return r
}

func (t *Target) send(msg *Message) error {
	t.ResetSendBuffer()

	serialized, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = pvm.PackfString("%s", string(serialized))
	if err != nil {
		return err
	}
	err = pvm.Send(t.TaskId, msg.Id)
	if err != nil {
		return err
	}

	return nil
}

type ReceiveResult struct {
	Response *MessageResponse
	Err      error
}

func (t *Target) receiveContinuously(id int) <-chan *ReceiveResult {
	r := make(chan *ReceiveResult)

	fin := time.Now().Add(time.Duration(RequestTimeoutSeconds * int(time.Second)))

	go (func() {
		for {
			bufId, err := pvm.Nrecv(t.TaskId, id)
			if err != nil {
				r <- &ReceiveResult{Err: err}
			}

			if bufId == 0 {
				if time.Now().After(fin) {
					r <- &ReceiveResult{Err: fmt.Errorf("request timeout")}
					break
				}
				continue
			}

			content, err := pvm.UnpackfString("%s", MaxPacketSize)

			if err != nil {
				r <- &ReceiveResult{Err: err}
			}

			deserialized := &MessageResponse{}

			err = json.Unmarshal([]byte(content), deserialized)

			if err != nil {
				r <- &ReceiveResult{Err: err}
			}

			r <- &ReceiveResult{Response: deserialized}
			break
		}
	})()

	return r
}
