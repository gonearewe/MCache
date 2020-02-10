package client

const (
	RequestSet RequestType = "set"
	RequestGet RequestType = "get"
)

type RequestType string

type Request struct {
	Type  RequestType
	Key   string
	Val   string
	Error string // storing error response
}

type Client interface {
	Run(req *Request)
	PipelineRun(reqs []*Request)
}

func New() Client {
	return newTCPClient()
}
