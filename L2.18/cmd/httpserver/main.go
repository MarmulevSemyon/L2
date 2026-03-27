package main

import (
	"flag"
	"log"
	"net/http"

	"httpserver/internal/calendar"
	myhttp "httpserver/internal/http"
	"httpserver/internal/middleware"
)

func main() {
	port := flag.String("port", "8080", "server port")
	flag.Parse()

	cal := calendar.NewCalendar()
	handler := myhttp.NewHandler(cal)

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", handler.CreateEvent)
	mux.HandleFunc("/update_event", handler.UpdateEvent)
	mux.HandleFunc("/delete_event", handler.DeleteEvent)
	mux.HandleFunc("/events_for_day", handler.EventsForDay)
	mux.HandleFunc("/events_for_week", handler.EventsForWeek)
	mux.HandleFunc("/events_for_month", handler.EventsForMonth)

	loggedMux := middleware.LoggingMiddleware(mux)

	log.Println("server started on :" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, loggedMux))
}
