package routes

import (
	"net/http"

	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/data"
	"github.com/sirupsen/logrus"
)

type ManagerHandler struct {
	log *logrus.Logger
	db  *data.DatabaseService
}

func NewManagerHandler(log *logrus.Logger, db *data.DatabaseService) *ManagerHandler {
	return &ManagerHandler{
		log: log,
		db:  db,
	}
}

func (h *ManagerHandler) GetHello(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("{ message: 'it works' }"))
}
