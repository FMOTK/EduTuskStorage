package myApi

import (
	"context"
	"encoding/json"
	"fmt"
	storage "httpServer/TuskStorage"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	storage *storage.TuskStorage
	ctx     context.Context
}

func NewHandler(ctx context.Context, storage *storage.TuskStorage) *Handler {
	return &Handler{
		storage: storage,
		ctx:     ctx,
	}
}

type Response struct {
	Status string `json:"status"`
	Id     string `json:"id"`
}

func (h *Handler) CreateTuskHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	stringDuration := query.Get("duration")
	duration, err := time.ParseDuration(stringDuration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Bad reqest")
		return
	}

	t := h.storage.CreateTusk(duration)

	resp := Response{
		Status: t.GetStatus(),
		Id:     t.GetUUID(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	go func() {
		h.storage.UpdateTuskById(t.GetUUID(), storage.RunningStatus)

		select {
		case <-h.ctx.Done():
			h.storage.UpdateTuskById(t.GetUUID(), storage.FailedStatus)
			return
		case <-time.After(duration):
			h.storage.UpdateTuskById(t.GetUUID(), storage.CompletedStatus)
			return
		}
	}()
}

func (h *Handler) GetTuskStatusHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	status, err := h.storage.GetTuskStatuById(id)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Tusk status is %s", status)
}
