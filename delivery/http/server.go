package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandler,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/miners").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetMiners)
	router.Path("/miners").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateMiner)

	router.Path("/equipments").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetEquipments)
	router.Path("/equipments").Methods("POST").HandlerFunc(s.httpHandlers.HandleBuyEquipment)

	router.Path("/enterprise").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetEnterpriseStats)
	router.Path("/enterprise").Methods("POST").HandlerFunc(s.httpHandlers.HandleShutdownGame)

	fmt.Printf("Server starting")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Server error: ", err.Error())
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
