package app

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/ezcnrmn/quiz_service/internal/app/controller"
	"github.com/ezcnrmn/quiz_service/internal/app/router"
	"github.com/ezcnrmn/quiz_service/internal/app/service"
)

type App struct {
	addr    string
	server  *http.Server
	loggers *router.Loggers
}

func New(addr string) (*App, error) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	service := service.New()

	controller := controller.New(service)

	loggers := &router.Loggers{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	router := router.New(loggers, controller)

	mux := http.NewServeMux()
	mux.HandleFunc("/", router.MainRoute)
	mux.HandleFunc("/quiz", router.QuizListRoute)
	mux.HandleFunc("/quiz/", router.QuizRoute)

	server := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	return &App{
		addr:    addr,
		server:  server,
		loggers: loggers,
	}, nil
}

func (a *App) Run() {
	a.loggers.InfoLog.Printf("server is running on %s\n", a.addr)
	err := a.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		a.loggers.InfoLog.Println("server closed")
	} else if err != nil {
		a.loggers.ErrorLog.Fatal(err)
	}
}
