package cloud_database

type handler struct {
	service Service
}

func InitHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ReadDB() {
	// 1. search s3 store
	// 2. return file to client
}
