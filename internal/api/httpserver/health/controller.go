package health

import "net/http"

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

// @Summary Health check
// @Tags    misc
// @Produce plain
// @Success 200 {string} string "healthy"
// @Router  /health [get]
func (c *Controller) HealthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("Healthy as a clam"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
