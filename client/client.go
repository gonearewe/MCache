package client

const (
	RequestSet RequestType = "set"
	RequestGet RequestType = "get"
	RequestDel RequestType = "del"
)

type RequestType string

type Request struct {
	Type  RequestType
	Key   string
	Val   []byte
	Error error // storing error response
}

type Client interface {
	Run(req *Request)
	PipelinedRun(reqs []*Request)
}

func New(type_, server string) Client {
	switch type_ {
	case "tcp":
		return newTCPClient(server)
	case "redis":
		return newRedisClient(server)
	}

	panic("TODO")
}
