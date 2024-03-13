package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ezcnrmn/quiz_service/internal/app/utils"
)

type Controller interface {
	GetQuiz(id string) (bytes []byte, err error)
	TryQuiz(id string, body io.ReadCloser) (res []byte, err error)
	GetAllQuizzes() (bytes []byte, err error)
	CreateQuiz(body io.ReadCloser) (bytes []byte, err error)
	EditQuiz(body io.ReadCloser) (bytes []byte, err error)
	DeleteQuiz(body io.ReadCloser) (bytes []byte, err error)
}

type Loggers struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type QuizRouter struct {
	loggers    *Loggers
	controller Controller
}

func New(loggers *Loggers, controller Controller) *QuizRouter {
	return &QuizRouter{
		loggers:    loggers,
		controller: controller,
	}
}

func (qr *QuizRouter) MainRoute(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(writer, req)
		return
	}

	writer.Write([]byte("quiz server is working"))
}

func (qr *QuizRouter) QuizListRoute(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		quiz, err := qr.controller.GetAllQuizzes()
		if err != nil {
			qr.loggers.ErrorLog.Println(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(quiz)
	case http.MethodPost:
		quiz, err := qr.controller.CreateQuiz(req.Body)
		if err != nil {
			qr.loggers.ErrorLog.Println(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(quiz)
	case http.MethodPatch:
		quiz, err := qr.controller.EditQuiz(req.Body)
		if err != nil {
			qr.loggers.ErrorLog.Println(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(quiz)
	case http.MethodDelete:
		id, err := qr.controller.DeleteQuiz(req.Body)
		if err != nil {
			qr.loggers.ErrorLog.Println(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(id)
	default:
		writer.Header().Set("Allow", utils.GetAllowMethodsString(http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete))
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		fmt.Fprintf(writer, "for quiz list %s method is not processed", req.Method)
	}
}

func (qr *QuizRouter) QuizRoute(writer http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/quiz/")

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		quiz, err := qr.controller.GetQuiz(id)
		if err != nil {
			qr.loggers.ErrorLog.Println(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(quiz)
	case http.MethodPost:
		result, err := qr.controller.TryQuiz(id, req.Body)
		if err != nil {
			qr.loggers.ErrorLog.Println(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(result)
	default:
		writer.Header().Set("Allow", utils.GetAllowMethodsString(http.MethodGet, http.MethodPost))
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		fmt.Fprintf(writer, "for quiz %s method is not processed", req.Method)
	}
}
