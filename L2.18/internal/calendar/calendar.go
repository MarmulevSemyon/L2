package calendar

import (
	"errors"
	"sync"
	"time"
)

// Event описывает событие календаря.
type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Text   string    `json:"event"`
}

// Calendar хранит события календаря в памяти и предоставляет методы для работы с ними.
type Calendar struct {
	mu     sync.RWMutex
	events map[int][]Event
	nextID int
}

// ErrEventNotFound возвращается, если событие с указанным идентификатором не найдено.
var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDate   = errors.New("invalid date")
)

// NewCalendar создаёт и возвращает новый экземпляр календаря.
func NewCalendar() *Calendar {
	return &Calendar{
		events: make(map[int][]Event),
		nextID: 1,
	}
}

// CreateEvent создаёт новое событие, присваивает ему идентификатор и сохраняет в календаре.
func (c *Calendar) CreateEvent(e Event) (Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e.ID = c.nextID
	c.nextID++

	c.events[e.UserID] = append(c.events[e.UserID], e)
	return e, nil
}

// UpdateEvent обновляет существующее событие по его идентификатору.
func (c *Calendar) UpdateEvent(e Event) (Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	events := c.events[e.UserID]

	for ind := range events {
		if events[ind].ID == e.ID {
			// c.events[e.UserID][ind] = e
			events[ind] = e
			c.events[e.UserID] = events
			return e, nil
		}
	}
	return Event{}, ErrEventNotFound
}

// DeleteEvent удаляет событие пользователя по идентификатору события.
func (c *Calendar) DeleteEvent(userID, id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	events := c.events[userID]

	for ind := range events {
		if events[ind].ID == id {
			c.events[userID] = append(events[:ind], events[ind+1:]...)
			return nil
		}
	}

	return ErrEventNotFound
}

// EventsForDay возвращает все события пользователя за указанный день.
func (c *Calendar) EventsForDay(userID int, date time.Time) ([]Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var res []Event
	for _, ev := range c.events[userID] {
		if sameDay(date, ev.Date) {
			res = append(res, ev)
		}
	}

	return res, nil
}
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// EventsForWeek возвращает все события пользователя за неделю, в которую входит указанная дата.
func (c *Calendar) EventsForWeek(userID int, date time.Time) ([]Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var res []Event
	for _, ev := range c.events[userID] {
		if sameWeek(date, ev.Date) {
			res = append(res, ev)
		}
	}

	return res, nil
}
func sameWeek(a, b time.Time) bool {
	ay, aw := a.ISOWeek()
	by, bw := b.ISOWeek()
	return ay == by && aw == bw
}

// EventsForMonth возвращает все события пользователя за месяц, в который входит указанная дата.
func (c *Calendar) EventsForMonth(userID int, date time.Time) ([]Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var res []Event
	for _, ev := range c.events[userID] {
		if sameMonth(date, ev.Date) {
			res = append(res, ev)
		}
	}

	return res, nil
}
func sameMonth(a, b time.Time) bool {
	ay, am, _ := a.Date()
	by, bm, _ := b.Date()
	return ay == by && am == bm
}
