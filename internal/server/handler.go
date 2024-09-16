package server

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	Router  *mux.Router
	Service TenderService
	Server  *http.Server
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(service TenderService) *Handler {
	h := &Handler{Service: service}

	h.Router = mux.NewRouter()
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

// mapRoutes задает маршруты для API.
func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")
	h.Router.HandleFunc("/api/tenders/new", h.CreateTender).Methods("POST")
	h.Router.HandleFunc("/api/tenders", h.GetTenders).Methods("GET")
	h.Router.HandleFunc("/api/tenders/my", h.GetUserTenders).Methods("GET")
	h.Router.HandleFunc("/api/tenders/{tenderId}/status", h.GetTenderStatus).Methods("GET")
	h.Router.HandleFunc("/api/tenders/{tenderId}/status", h.UpdateTenderStatus).Methods("PUT")
	h.Router.HandleFunc("/api/tenders/{tenderId}/edit", h.EditTender).Methods("PATCH")
	h.Router.HandleFunc("/api/bids/new", h.CreateBid).Methods("POST")
	h.Router.HandleFunc("/api/bids/my", h.GetUserBids).Methods("GET")
	h.Router.HandleFunc("/api/bids/{tenderId}/list", h.GetBidsForTender).Methods("GET")
	h.Router.HandleFunc("/api/bids/{bidId}/status", h.GetBidStatus).Methods("GET")
	h.Router.HandleFunc("/api/bids/{bidId}/status", h.UpdateBidStatus).Methods("PUT")
	h.Router.HandleFunc("/api/bids/{bidId}/edit", h.EditBid).Methods("PATCH")
}

// Serve запускает HTTP-сервер и обрабатывает остановку сервера.
func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)
	log.Println("shut down gracefully")
	return nil
}
