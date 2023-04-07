package pvmrpc

type Target struct {
	TaskId int
}

type MessageType string
type Message struct {
	// IMPORTANT: Generated randomly by the caller
	Id int
	// Will not change on response (will be client tId)
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
