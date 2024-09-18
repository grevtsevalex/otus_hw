package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/logger"
	rabbitqueue "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/queue/rabbitQueue"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewConfig(configFile)
	if err != nil {
		err = fmt.Errorf("config initialization: %w", err)
		fmt.Println(err.Error())
		return
	}

	logg := logger.New(config.Logger.Level, os.Stdout)

	st, err := getStorage(config)
	if err != nil {
		err = fmt.Errorf("storage initialization: %w", err)
		fmt.Println(err.Error())
		return
	}
	logg.Info(fmt.Sprintf("Получили объект хранилища, тип: %s", config.Storage.Type))

	// create queue.
	q := rabbitqueue.NewQueue(config.Queue.Address, logg)

	// create scheduler.
	daemon := scheduler.NewScheduler(st, logg, q, time.Second*time.Duration(config.Cron.Period))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		daemon.Start(ctx)
		logg.Info("scheduler is running...")
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	daemon.Stop(ctx)
}

// getStorage получить объект хранилища.
func getStorage(conf Config) (storage.EventStorage, error) {
	var storage storage.EventStorage
	var err error
	switch conf.Storage.Type {
	case "mem":
		storage = memorystorage.New()
	case "sql":
		dbConf := sqlstorage.Config{DBName: conf.DB.Name, User: conf.DB.User, Pass: conf.DB.Pass, Host: conf.DB.Host}
		storage, err = sqlstorage.New(dbConf)
		if err != nil {
			return nil, fmt.Errorf("получение хранилища: %w", err)
		}
	}

	return storage, nil
}
