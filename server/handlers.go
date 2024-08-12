package server

import (
	"L0/entities"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

func (s *Server) SayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
	w.Write([]byte("Don't forget to smile, älskling 😌"))
}

func (s *Server) GetOrderById(w http.ResponseWriter, r *http.Request) {
	orderReq := &entities.Order{}

	log := s.log.With(
		slog.String("method", "GetOrderById"))

	orderUID := chi.URLParam(r, "id")

	err := render.DecodeJSON(r.Body, orderReq)
	if err != nil {
		log.Error("Failed to decode JSON", err)
		render.JSON(w, r, map[string]string{"error": "Failed to decode JSON"})
		return
	}

	orderReq.OrderUID = orderUID

	response, err := s.db.GetOrderById(r.Context(), orderUID)
	if err != nil {
		log.Error("Failed to get order", err)
		render.JSON(w, r, map[string]string{"error": "Failed to get order"})
		return
	}

	log.Info("Got order", response)
	render.JSON(w, r, response)
}
