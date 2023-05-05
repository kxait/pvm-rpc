package pvm_rpc

type Target struct {
	TaskId int
}

func NewTarget(taskId int) *Target {
	targ := Target{TaskId: taskId}

	return &targ
}

type ReceiveResult struct {
	Response *MessageResponse
	Err      error
}

type MessageType string
type Message struct {
	Id           int
	CallerTaskId int
	Type         MessageType
	Content      string
}

func (m *Message) CreateResponse(content string) *Message {
	response := *m
	response.Content = content
	return &response
}

type MessageResponse struct {
	Message
	IsError bool
}

type RpcHandler func(*Message) (*Message, error)

type RpcServer struct {
	Handlers map[MessageType]RpcHandler
}
