package internalgrpc

import (
	"context"
	"fmt"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/app"
	eventpb "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EventService тип сервиса событий.
type EventService struct {
	eventpb.UnimplementedEventServiceServer
	app *app.App
}

// newEventService конструктор сервиса.
func newEventService(app *app.App) *EventService {
	return &EventService{app: app}
}

// Add добавление нового события.
func (s *EventService) Add(_ context.Context, req *eventpb.AddRequest) (*eventpb.AddResponse, error) {
	event := storage.Event{
		ID:                  storage.EventID(req.Event.GetId()),
		Title:               req.Event.GetTitle(),
		Description:         req.Event.GetDescription(),
		AuthorID:            req.Event.GetAuthorId(),
		StartDate:           req.Event.GetStartDate().AsTime(),
		EndDate:             req.Event.GetEndDate().AsTime(),
		HoursBeforeToNotify: int(req.Event.GetHoursBeforeToNotify()),
	}

	err := s.app.RegisterNewEvent(event)
	if err != nil {
		s.app.Logger.Error(fmt.Sprintf("обработка запроса на создание события: %s", err.Error()))

		return &eventpb.AddResponse{Result: false}, err
	}

	return &eventpb.AddResponse{Result: true, EventID: string(event.ID)}, nil
}

// Update обновление события.
func (s *EventService) Update(_ context.Context, req *eventpb.UpdateRequest) (*eventpb.UpdateResponse, error) {
	event := storage.Event{
		ID:                  storage.EventID(req.Event.GetId()),
		Title:               req.Event.GetTitle(),
		Description:         req.Event.GetDescription(),
		AuthorID:            req.Event.GetAuthorId(),
		StartDate:           req.Event.GetStartDate().AsTime(),
		EndDate:             req.Event.GetEndDate().AsTime(),
		HoursBeforeToNotify: int(req.Event.GetHoursBeforeToNotify()),
	}

	err := s.app.ChangeEvent(event)
	if err != nil {
		s.app.Logger.Error(fmt.Sprintf("обработка запроса на изменение события: %s", err.Error()))

		return &eventpb.UpdateResponse{Result: false}, err
	}

	return &eventpb.UpdateResponse{Result: true}, nil
}

// Delete удаление события.
func (s *EventService) Delete(_ context.Context, req *eventpb.DeleteRequest) (*eventpb.DeleteResponse, error) {
	err := s.app.DeleteEvent(storage.EventID(req.EventID))
	if err != nil {
		s.app.Logger.Error(fmt.Sprintf("обработка запроса на удаление события: %s", err.Error()))

		return &eventpb.DeleteResponse{Result: false}, err
	}

	return &eventpb.DeleteResponse{Result: true}, nil
}

// GetAll получение списка событий.
func (s *EventService) GetAll(_ context.Context, _ *emptypb.Empty) (*eventpb.GetAllResponse, error) {
	events, err := s.app.GetAllEvents()
	if err != nil {
		s.app.Logger.Error(fmt.Sprintf("обработка запроса на получение списка всех событий: %s", err.Error()))

		return &eventpb.GetAllResponse{Events: []*eventpb.Event{}}, err
	}

	result := make([]*eventpb.Event, 0, len(events))
	for _, e := range events {
		result = append(result, &eventpb.Event{
			Id:                  string(e.ID),
			Title:               e.Title,
			Description:         e.Description,
			HoursBeforeToNotify: int32(e.HoursBeforeToNotify),
			StartDate:           timestamppb.New(e.StartDate),
			EndDate:             timestamppb.New(e.EndDate),
			AuthorId:            e.AuthorID,
		})
	}

	return &eventpb.GetAllResponse{Events: result}, nil
}

// GetDayEvents получение списка событий за день.
func (s *EventService) GetDayEvents(_ context.Context,
	req *eventpb.GetEventsByRangeRequest,
) (*eventpb.GetDayEventsResponse, error) {
	events, err := s.app.GetEventsForDay(req.GetFrom().AsTime())
	if err != nil {
		s.app.Logger.Error(fmt.Sprintf("обработка запроса на получение списка событий за день: %s", err.Error()))

		return &eventpb.GetDayEventsResponse{Events: []*eventpb.Event{}}, err
	}
	return &eventpb.GetDayEventsResponse{Events: transformEventsToResponse(events)}, nil
}

// GetWeekEvents получение списка событий за неделю.
func (s *EventService) GetWeekEvents(_ context.Context,
	req *eventpb.GetEventsByRangeRequest,
) (*eventpb.GetWeekEventsResponse, error) {
	events, err := s.app.GetEventsForWeek(req.GetFrom().AsTime())
	if err != nil {
		s.app.Logger.Error(fmt.Sprintf("обработка запроса на получение списка событий за неделю: %s", err.Error()))

		return &eventpb.GetWeekEventsResponse{Events: []*eventpb.Event{}}, err
	}
	return &eventpb.GetWeekEventsResponse{Events: transformEventsToResponse(events)}, nil
}

// GetMonthEvents получение списка событий за месяц.
func (s *EventService) GetMonthEvents(_ context.Context,
	req *eventpb.GetEventsByRangeRequest,
) (*eventpb.GetMonthEventsResponse, error) {
	events, err := s.app.GetEventsForMonth(req.GetFrom().AsTime())
	if err != nil {
		s.app.Logger.Error(fmt.Sprintf("обработка запроса на получение списка событий за месяц: %s", err.Error()))

		return &eventpb.GetMonthEventsResponse{Events: []*eventpb.Event{}}, err
	}
	return &eventpb.GetMonthEventsResponse{Events: transformEventsToResponse(events)}, nil
}

// transformEventsToResponse приведение соботий к формату ответа.
func transformEventsToResponse(events storage.Events) []*eventpb.Event {
	result := make([]*eventpb.Event, 0, len(events))
	for _, e := range events {
		result = append(result, &eventpb.Event{
			Id:                  string(e.ID),
			Title:               e.Title,
			Description:         e.Description,
			HoursBeforeToNotify: int32(e.HoursBeforeToNotify),
			StartDate:           timestamppb.New(e.StartDate),
			EndDate:             timestamppb.New(e.EndDate),
			AuthorId:            e.AuthorID,
		})
	}
	return result
}
