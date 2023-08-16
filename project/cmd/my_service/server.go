package my_service

type Server struct {
	Service *ItemService
}

func NewServer(service *ItemService) *Server {
	return &Server{Service: service}
}
