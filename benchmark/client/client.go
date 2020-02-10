package client

const (
	RequestSet  RequestType = "set"
	RequestGet  RequestType = "get"
	MissRequest RequestType = "miss" // not a actual request, used for identifying cache miss
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
	PipelineRun(reqs []*Request)
}

func New(type_, server string) Client {
	switch type_ {
	case "tcp":
		return newTCPClient(server)
	}
}
