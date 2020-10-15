package scheduler

type Response struct {
	Status bool `json:"status"`
}

func NewResponse(status bool) Response {
	return Response{
		Status: status,
	}
}
