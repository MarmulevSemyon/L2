package http

import (
	"encoding/json"
	"httpserver/internal/calendar"
	"net/http"
	"strconv"
	"time"
)

// Handler обрабатывает HTTP-запросы календаря. торчит через NewHandler
type Handler struct {
	cal *calendar.Calendar
}

type response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type createEventRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"` // YYYY-MM-DD
	Event  string `json:"event"`
}

type updateEventRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Event  string `json:"event"`
}

// NewHandler создаёт новый HTTP-обработчик для работы с календарём.
func NewHandler(cal *calendar.Calendar) *Handler {
	return &Handler{cal: cal}
}

// CreateEvent обрабатывает HTTP-запрос на создание нового события.
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// проверка, что это POST
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusBadRequest, response{Error: "bad request"})
		return
	}

	// декодируем r (response) в req
	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid json"})
		return
	}

	// парсим дату из string в time.Time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid date"})
		return
	}

	// создаём событие (CreateEvent кладет его в мапу)
	event, err := h.cal.CreateEvent(calendar.Event{
		UserID: req.UserID,
		Date:   date,
		Text:   req.Event,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	// пишем в w http.ResponseWriter
	writeJSON(w, http.StatusOK, response{Result: event})
}
func writeJSON(w http.ResponseWriter, status int, resp response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// UpdateEvent обрабатывает HTTP-запрос на обновление существующего события.
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// проверка, что это POST
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid method"})
		return
	}

	// декодируем r (response) в req
	var req updateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid json"})
		return
	}

	// парсим дату из string в time.Time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid date"})
		return
	}

	// обновляем событие (UpdateEvent изменяет мапу событий)
	event, err := h.cal.UpdateEvent(calendar.Event{
		ID:     req.ID,
		UserID: req.UserID,
		Date:   date,
		Text:   req.Event,
	})

	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, response{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, response{Result: event})
}

type deleteEventRequest struct {
	UserID int `json:"user_id"`
	ID     int `json:"id"`
}

// DeleteEvent обрабатывает HTTP-запрос на удаление события.
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	// проверка, что это POST
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid method"})
		return
	}
	// декодируем r (response) в req
	var req deleteEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid json"})
		return
	}
	// удаляем событие
	err := h.cal.DeleteEvent(req.UserID, req.ID)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, response{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, response{Result: "deleted"})
}

// EventsForDay обрабатывает HTTP-запрос на получение событий за день.
func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	// проверка, что это GET
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid method"})
		return
	}

	// забираем из запроса дату и ID пользователя
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	// проверка на валидность
	if userID == "" || dateStr == "" {
		writeJSON(w, http.StatusBadRequest, response{Error: "missing params, user_id: " + userID + "\tdate: " + dateStr})
		return
	}

	// перевод ID пользователя в int
	uid, err := strconv.Atoi(userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid user_id"})
		return
	}

	// парсим дату из string в time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid date"})
		return
	}

	//  получаем все события в нужный день
	events, err := h.cal.EventsForDay(uid, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, response{Result: events})
}

// EventsForWeek обрабатывает HTTP-запрос на получение событий за неделю.
func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	// проверка, что это GET
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid method"})
		return
	}

	// забираем из запроса дату и ID пользователя
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	// проверка на валидность
	if userID == "" || dateStr == "" {
		writeJSON(w, http.StatusBadRequest, response{Error: "missing params, user_id: " + userID + "\tdate: " + dateStr})
		return
	}

	// перевод ID пользователя в int
	uid, err := strconv.Atoi(userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid user_id"})
		return
	}

	// парсим дату из string в time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid date"})
		return
	}

	//  получаем все события в нужный день
	events, err := h.cal.EventsForWeek(uid, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, response{Result: events})
}

// EventsForMonth обрабатывает HTTP-запрос на получение событий за месяц.
func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	// проверка, что это GET
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid method"})
		return
	}

	// забираем из запроса дату и ID пользователя
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	// проверка на валидность
	if userID == "" || dateStr == "" {
		writeJSON(w, http.StatusBadRequest, response{Error: "missing params, user_id: " + userID + "\tdate: " + dateStr})
		return
	}

	// перевод ID пользователя в int
	uid, err := strconv.Atoi(userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid user_id"})
		return
	}

	// парсим дату из string в time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid date"})
		return
	}

	//  получаем все события в нужный день
	events, err := h.cal.EventsForMonth(uid, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, response{Result: events})
}
