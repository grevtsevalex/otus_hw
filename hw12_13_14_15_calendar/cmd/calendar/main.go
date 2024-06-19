package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		err = fmt.Errorf("config initialization: %w", err)
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	logg := logger.New(config.Logger.Level, os.Stdout)

	eventStorage, err := getStorage(config)
	if err != nil {
		err = fmt.Errorf("storage initialization: %w", err)
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	logg.Info(fmt.Sprintf("Получили объект хранилища, тип: %s", config.Storage.Type))

	calendar := app.New(logg, eventStorage)

	server := internalhttp.NewServer(logg, calendar, internalhttp.Config{
		Port:            config.Server.Port,
		HandlerTimeoutS: config.Server.HandlerTimeoutS,
		WriteTimeoutMS:  config.Server.WriteTimeoutMS,
		ReadTimeoutMS:   config.Server.ReadTimeoutMS,
	})

	logg.Info("calendar is running...")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
}

// getStorage получить объект хранилища.
func getStorage(conf Config) (storage.EventStorage, error) {
	var storage storage.EventStorage
	var err error
	switch conf.Storage.Type {
	case "mem":
		storage = memorystorage.New()
	case "sql":
		dbConf := sqlstorage.Config{DBName: conf.DB.Name, User: conf.DB.User, Pass: conf.DB.Pass}
		storage, err = sqlstorage.New(dbConf)
		if err != nil {
			return nil, fmt.Errorf("получение хранилища: %w", err)
		}
	}

	return storage, nil
}
