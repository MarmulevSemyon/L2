package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()

	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("failed to parse date %q: %v", s, err)
	}
	return d
}

func TestCreateEvent(t *testing.T) {
	cal := NewCalendar()

	event, err := cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "first event",
	})

	require.NoError(t, err)
	assert.Equal(t, 1, event.UserID)
	assert.Equal(t, "first event", event.Text)

	events, err := cal.EventsForDay(1, mustDate(t, "2026-03-27"))

	require.NoError(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, events[0].ID, event.ID)
}

func TestCreateEventAssignsIncrementalIDs(t *testing.T) {
	cal := NewCalendar()

	first, err := cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "first",
	})
	require.NoError(t, err)

	second, err := cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-28"),
		Text:   "second",
	})
	require.NoError(t, err)

	assert.Equal(t, 1, first.ID)
	assert.Equal(t, 2, second.ID)
}

func TestUpdateEvent(t *testing.T) {
	cal := NewCalendar()

	created, err := cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "old text",
	})
	require.NoError(t, err)

	updatedDate := mustDate(t, "2026-03-28")
	updated, err := cal.UpdateEvent(Event{
		ID:     created.ID,
		UserID: 1,
		Date:   updatedDate,
		Text:   "new text",
	})

	require.NoError(t, err)

	assert.Equal(t, 1, updated.ID)
	assert.Equal(t, 1, created.ID)

	assert.Equal(t, "new text", updated.Text)
	assert.True(t, updated.Date.Equal(updatedDate))

	eventsOldDay, err := cal.EventsForDay(1, mustDate(t, "2026-03-27"))
	require.NoError(t, err)
	assert.Equal(t, 0, len(eventsOldDay))

	eventsNewDay, err := cal.EventsForDay(1, mustDate(t, "2026-03-28"))
	require.NoError(t, err)
	assert.Equal(t, 1, len(eventsNewDay))
	assert.Equal(t, "new text", eventsNewDay[0].Text)
}

func TestUpdateEventNotFound(t *testing.T) {
	cal := NewCalendar()

	_, err := cal.UpdateEvent(Event{
		ID:     999,
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "missing",
	})

	require.ErrorIs(t, err, ErrEventNotFound)
}

func TestDeleteEvent(t *testing.T) {
	cal := NewCalendar()

	created, err := cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "to delete",
	})
	require.NoError(t, err)

	err = cal.DeleteEvent(1, created.ID)
	require.NoError(t, err)

	events, err := cal.EventsForDay(1, mustDate(t, "2026-03-27"))
	require.NoError(t, err)
	assert.Equal(t, 0, len(events))

}

func TestDeleteEventNotFound(t *testing.T) {
	cal := NewCalendar()

	err := cal.DeleteEvent(1, 999)
	require.ErrorIs(t, err, ErrEventNotFound)
}

func TestEventsForDay(t *testing.T) {
	cal := NewCalendar()

	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "event 1",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "event 2",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-28"),
		Text:   "event 3",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 2,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "other user event",
	})

	events, err := cal.EventsForDay(1, mustDate(t, "2026-03-27"))
	require.NoError(t, err)
	assert.Equal(t, 2, len(events))
}

func TestEventsForWeek(t *testing.T) {
	cal := NewCalendar()

	//  даты в одной неделе.
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-23"),
		Text:   "monday",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "friday",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-29"),
		Text:   "sunday",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-30"),
		Text:   "next week",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 2,
		Date:   mustDate(t, "2026-03-27"),
		Text:   "other user",
	})

	events, err := cal.EventsForWeek(1, mustDate(t, "2026-03-27"))
	require.NoError(t, err)
	assert.Equal(t, 3, len(events))
}

func TestEventsForMonth(t *testing.T) {
	cal := NewCalendar()

	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-01"),
		Text:   "march 1",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-03-15"),
		Text:   "march 15",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 1,
		Date:   mustDate(t, "2026-04-01"),
		Text:   "april 1",
	})
	_, _ = cal.CreateEvent(Event{
		UserID: 2,
		Date:   mustDate(t, "2026-03-20"),
		Text:   "other user",
	})

	events, err := cal.EventsForMonth(1, mustDate(t, "2026-03-27"))
	require.NoError(t, err)
	assert.Equal(t, 2, len(events))
}
