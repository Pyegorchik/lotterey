package handler

import (
	"encoding/json"
	"lottery/internal/domain"
	"lottery/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// LotteryHandler обработчик HTTP запросов
type LotteryHandler struct {
	service service.LotteryService
}

func NewLotteryHandler(service service.LotteryService) *LotteryHandler {
	return &LotteryHandler{service: service}
}

func (h *LotteryHandler) CreateDraw(w http.ResponseWriter, r *http.Request) {
	draw, err := h.service.CreateDraw(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, draw)
}

func (h *LotteryHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req domain.TicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ticket, err := h.service.CreateTicket(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, ticket)
}

func (h *LotteryHandler) CloseDraw(w http.ResponseWriter, r *http.Request) {
	drawID, err := h.getDrawIDFromURL(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid draw ID")
		return
	}

	draw, err := h.service.CloseDraw(r.Context(), drawID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, draw)
}

func (h *LotteryHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	drawID, err := h.getDrawIDFromURL(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid draw ID")
		return
	}

	results, err := h.service.GetDrawResults(r.Context(), drawID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, results)
}

func (h *LotteryHandler) GetDraw(w http.ResponseWriter, r *http.Request) {
	drawID, err := h.getDrawIDFromURL(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid draw ID")
		return
	}

	draw, err := h.service.GetDraw(r.Context(), drawID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, draw)
}

func (h *LotteryHandler) GetDraws(w http.ResponseWriter, r *http.Request) {
	draws, err := h.service.GetAllDraws()
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, draws)
}

// Вспомогательные методы handler'а
func (h *LotteryHandler) getDrawIDFromURL(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["draw_id"])
}

func (h *LotteryHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *LotteryHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *LotteryHandler) handleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *service.BusinessError:
		h.writeError(w, http.StatusConflict, err.Error())
	case *service.ValidationError:
		h.writeError(w, http.StatusBadRequest, err.Error())
	case *service.NotFoundError:
		h.writeError(w, http.StatusNotFound, err.Error())
	default:
		h.writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}
