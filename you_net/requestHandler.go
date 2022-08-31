package you_net

type RequestHandler interface {
	PreHandle(*YouRequest)
	Handle(*YouRequest)
	PostHandle(*YouRequest)
}
