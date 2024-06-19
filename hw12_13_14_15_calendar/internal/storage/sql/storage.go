package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib" // pgx
)

// Config модель конфига для хранилища.
type Config struct {
	DBName          string
	User            string
	Pass            string
	PoolSize        int
	MaxConnLifeTime int
}

// Strorage модель хранилища.
type Storage struct {
	conf Config
	pool *sql.DB
	ctx  context.Context
}

// New конструктор хранилища.
func New(config Config) (*Storage, error) {
	ctx := context.Background()
	st := &Storage{conf: config, ctx: ctx}
	err := st.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connnection to DB: %w", err)
	}

	return st, nil
}

// Connect подключиться к БД.
func (s *Storage) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s", s.conf.User, s.conf.DBName, s.conf.Pass)
	db, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	db.SetConnMaxLifetime(time.Duration(s.conf.MaxConnLifeTime) * time.Second)
	db.SetMaxOpenConns(s.conf.PoolSize)

	s.pool = db
	return nil
}

// Close закрыть подключение к БД.
func (s *Storage) Close() error {
	err := s.pool.Close()
	if err != nil {
		return fmt.Errorf("failed to close connect to db: %w", err)
	}

	return nil
}

// Add добавить событие.
func (s *Storage) Add(event storage.Event) error {
	query := `INSERT INTO events(id, title, start_stamp, end_stamp, description, author_id, hours_before_to_notify)
	 VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.pool.ExecContext(
		s.ctx, query, event.ID,
		event.Title, event.StartDate,
		event.EndDate, event.Description,
		event.AuthorID, event.HoursBeforeToNotify)

	var pgxError pgx.PgError
	if errors.As(err, &pgxError) {
		if pgxError.Code == "23505" {
			return storage.ErrEventIDIsAlreadyExists
		}
	}
	if err != nil {
		return fmt.Errorf("добавление события: %w", err)
	}
	return nil
}

// Update обновить событие.
func (s *Storage) Update(event storage.Event) error {
	query := `UPDATE events SET (title, start_stamp, end_stamp, description, hours_before_to_notify)
	 = ($1, $2, $3, $4, $5) WHERE id = $6`
	res, err := s.pool.ExecContext(s.ctx, query, event.Title, event.StartDate,
		event.EndDate, event.Description, event.HoursBeforeToNotify, event.ID)
	if err != nil {
		return fmt.Errorf("обновление события: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("кол-во измененных строк: %w", err)
	}

	if n == 0 {
		return storage.ErrNoEvent
	}
	return nil
}

// Delete удалить событие.
func (s *Storage) Delete(eventID storage.EventID) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := s.pool.ExecContext(s.ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("удаление события: %w", err)
	}
	return nil
}

// GetAll получить все события.
func (s *Storage) GetAll() ([]storage.Event, error) {
	result := make([]storage.Event, 0)

	query := `SELECT * FROM events`
	rows, err := s.pool.QueryContext(s.ctx, query)
	if err != nil {
		return result, fmt.Errorf("получение всех событий: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event); err != nil {
			return result, fmt.Errorf("сканирование строки: %w", err)
		}
		result = append(result, event)
	}
	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("обработка получения событий: %w", err)
	}

	return result, nil
}

// Get получить событие.
func (s *Storage) Get(eventID storage.EventID) (storage.Event, error) {
	query := "SELECT * FROM events WHERE id = $1"
	row := s.pool.QueryRowContext(s.ctx, query, eventID)
	var event storage.Event
	err := row.Scan(&event.ID, &event.Title, &event.StartDate, &event.EndDate,
		&event.Description, &event.AuthorID, &event.HoursBeforeToNotify)
	if err == sql.ErrNoRows {
		return event, storage.ErrNoEvent
	} else if err != nil {
		return event, err
	}

	return event, nil
}
