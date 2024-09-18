package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	eventpb "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var serverPort int

func init() {
	flag.IntVar(&serverPort, "port", 0, "grpc server port")
}

func main() {
	flag.Parse()

	if serverPort == 0 {
		log.Println("No server port in args")
		os.Exit(1)
	}

	run(serverPort)
}

// run запуск тестов.
func run(port int) {
	conn, err := grpc.NewClient(fmt.Sprintf("calendar:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := eventpb.NewEventServiceClient(conn)
	ctx := metadata.NewOutgoingContext(context.Background(), nil)

	// -- add event
	fmt.Println("add event")
	firstEventID := "1"
	newEventTitle := "simple title"
	newEvent := createEvent(firstEventID, newEventTitle)
	req := &eventpb.AddRequest{
		Event: newEvent,
	}

	res, err := client.Add(ctx, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if res.EventID != firstEventID {
		fmt.Println("Неверный идентификатор созданного события")
		os.Exit(1)
	}

	if !res.Result {
		fmt.Println("Ошибка при создании события")
		os.Exit(1)
	}

	allEventsResponse, err := client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	eventExists := false
	for _, event := range allEventsResponse.Events {
		if event.Id == firstEventID && event.Title == newEventTitle {
			eventExists = true
			break
		}
	}

	if !eventExists {
		fmt.Println("Событие не добавлено в хранилище")
		os.Exit(1)
	}

	// -- -- --

	// -- add event with error duplicate id
	fmt.Println("add event with error duplicate id")
	newEventTitle = "simple title2"
	newEvent = createEvent(firstEventID, newEventTitle)
	req = &eventpb.AddRequest{
		Event: newEvent,
	}

	_, err = client.Add(ctx, req)
	if err == nil {
		fmt.Println("Прошло добавление события с дублирующимся ID")
		os.Exit(1)
	}
	// -- -- --

	// -- add second event on tomorrow
	fmt.Println("add second event on tomorrow")
	secondEventID := "2"
	newEventTitle = "simple title 2"
	newEvent = createTomorrowEvent(secondEventID, newEventTitle)
	req = &eventpb.AddRequest{
		Event: newEvent,
	}

	res, err = client.Add(ctx, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if res.EventID != secondEventID {
		fmt.Println("Неверный идентификатор созданного второго события")
		os.Exit(1)
	}

	if !res.Result {
		fmt.Println("Ошибка при создании второго события")
		os.Exit(1)
	}

	allEventsResponse, err = client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	eventExists = false
	for _, event := range allEventsResponse.Events {
		if event.Id == secondEventID && event.Title == newEventTitle {
			eventExists = true
			break
		}
	}

	if !eventExists {
		fmt.Println("Событие не добавлено в хранилище")
		os.Exit(1)
	}

	if len(allEventsResponse.Events) != 2 {
		fmt.Println("Событий в хранилище должно быть 2 на этот момент")
		os.Exit(1)
	}
	// -- -- --

	getEventsTesting(ctx, client, firstEventID, secondEventID)

	os.Exit(0)
}

// getEventsTesting тестирование получения событий.
func getEventsTesting(
	ctx context.Context,
	client eventpb.EventServiceClient,
	firstEventID string,
	secondEventID string,
) {
	// -- get today events
	fmt.Println("get today events")

	firstEventDate, err := time.Parse(time.RFC1123, "Mon, 16 Sep 2024 00:00:00 UTC")
	if err != nil {
		fmt.Println("Ошибка парсинга", err.Error())
	}

	todayEvents, err := client.GetDayEvents(ctx, &eventpb.GetEventsByRangeRequest{From: timestamppb.New(firstEventDate)})
	if err != nil {
		fmt.Println("ошибка получения событий за день")
		os.Exit(1)
	}

	if len(todayEvents.GetEvents()) != 1 {
		fmt.Println("должно быть одно событие")
		os.Exit(1)
	}

	for _, te := range todayEvents.GetEvents() {
		if te.Id != firstEventID {
			fmt.Println("не получено первое событие")
			os.Exit(1)
		}
	}
	// -- -- --

	// -- get week events
	fmt.Println("get week events")

	weekEvents, err := client.GetWeekEvents(ctx, &eventpb.GetEventsByRangeRequest{From: timestamppb.New(firstEventDate)})
	if err != nil {
		fmt.Println("ошибка получения событий за день")
		os.Exit(1)
	}

	if len(weekEvents.GetEvents()) != 2 {
		fmt.Println("должно быть два события")
		os.Exit(1)
	}

	for _, te := range weekEvents.GetEvents() {
		if te.Id != firstEventID && te.Id != secondEventID {
			fmt.Println("не получено недельное событие")
			os.Exit(1)
		}
	}
	// -- -- --

	// -- get month events
	fmt.Println("get month events")

	monthEvents, err := client.GetMonthEvents(ctx, &eventpb.GetEventsByRangeRequest{From: timestamppb.New(firstEventDate)})
	if err != nil {
		fmt.Println("ошибка получения событий за день")
		os.Exit(1)
	}

	if len(monthEvents.GetEvents()) != 2 {
		fmt.Println("должно быть два события")
		os.Exit(1)
	}

	for _, te := range monthEvents.GetEvents() {
		if te.Id != firstEventID && te.Id != secondEventID {
			fmt.Println("не получено месячное событие")
			os.Exit(1)
		}
	}
	// -- -- --
}

// createEvent дать событие.
func createEvent(ID string, title string) *eventpb.Event {
	startDate, err := time.Parse(time.RFC1123, "Mon, 16 Sep 2024 10:00:00 UTC")
	if err != nil {
		fmt.Println("Ошибка парсинга", err.Error())
	}

	return &eventpb.Event{
		Id:                  ID,
		Title:               title,
		StartDate:           timestamppb.New(startDate),
		EndDate:             timestamppb.New(startDate.Add(time.Hour * time.Duration(3))),
		Description:         "et Excepteur",
		AuthorId:            "1",
		HoursBeforeToNotify: 3,
	}
}

// createTomorrowEvent создать событие на завтра.
func createTomorrowEvent(ID string, title string) *eventpb.Event {
	startDate, err := time.Parse(time.RFC1123, "Tue, 17 Sep 2024 10:00:00 UTC")
	if err != nil {
		fmt.Println("Ошибка парсинга", err.Error())
	}
	return &eventpb.Event{
		Id:                  ID,
		Title:               title,
		StartDate:           timestamppb.New(startDate),
		EndDate:             timestamppb.New(startDate.Add(time.Hour * time.Duration(3))),
		Description:         "et Excepteur",
		AuthorId:            "1",
		HoursBeforeToNotify: 3,
	}
}
