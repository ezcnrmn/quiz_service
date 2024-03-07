package controller

import (
	"encoding/json"
	"io"

	"github.com/ezcnrmn/quiz_service/internal/app/model"
	"github.com/ezcnrmn/quiz_service/internal/app/utils/consts"
)

type Service interface {
	GetQuizzesList() ([]string, error)
	LoadAllQuizzes() ([]model.Quiz, error)
	LoadQuiz(fileName string) (quiz model.Quiz, err error)
	GetNewQuizId() (id string, err error)
	SaveQuiz(quiz model.Quiz) (err error)
	DeleteQuiz(id string) (err error)
	CheckIfQuizExist(id string) (res bool)
}

type QuizController struct {
	service Service
}

func New(service Service) *QuizController {
	return &QuizController{
		service: service,
	}
}

func (qc *QuizController) GetQuiz(id string) (bytes []byte, err error) {
	quiz, err := qc.service.LoadQuiz(id + consts.FILE_EXTENSION)
	if err != nil {
		return []byte{}, err
	}

	bytes, err = json.Marshal(quiz)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

type QuizError struct {
	RightAnswer   string `json:"rightAnswer"`
	QuestionIndex int    `json:"questionIndex"`
}
type Result struct {
	IsCorrect bool        `json:"isCorrect"`
	Errors    []QuizError `json:"errors"`
}

func (qc *QuizController) TryQuiz(id string, body io.ReadCloser) (res []byte, err error) {
	decoder := json.NewDecoder(body)

	var answers []string
	err = decoder.Decode(&answers)
	if err != nil {
		return []byte{}, err
	}

	quiz, err := qc.service.LoadQuiz(id + consts.FILE_EXTENSION)
	if err != nil {
		return []byte{}, err
	}

	results := Result{IsCorrect: true, Errors: []QuizError{}}
	for i, answer := range quiz.Answers {
		if answers[i] != answer {
			results.IsCorrect = false
			results.Errors = append(results.Errors, QuizError{QuestionIndex: i, RightAnswer: answer})
		}
	}

	res, err = json.Marshal(results)
	if err != nil {
		return []byte{}, err
	}

	return res, nil
}

func (qc *QuizController) GetAllQuizzes() (bytes []byte, err error) {
	// quizzes, err := loadAllQuizzes()
	quizzes, err := qc.service.GetQuizzesList()
	if err != nil {
		return []byte{}, err
	}

	bytes, err = json.Marshal(quizzes)

	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func (qc *QuizController) CreateQuiz(body io.ReadCloser) (bytes []byte, err error) {
	decoder := json.NewDecoder(body)

	var quiz model.Quiz
	err = decoder.Decode(&quiz)
	if err != nil {
		return []byte{}, err
	}

	newId, err := qc.service.GetNewQuizId()
	if err != nil {
		return []byte{}, err
	}

	quiz.Id = newId
	err = qc.service.SaveQuiz(quiz)
	if err != nil {
		return []byte{}, err
	}

	bytes, err = json.Marshal(quiz)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

type QuizIsNotExistError struct{}

func (e *QuizIsNotExistError) Error() string {
	return "quiz is not exist"
}

func (qc *QuizController) EditQuiz(body io.ReadCloser) (bytes []byte, err error) {
	decoder := json.NewDecoder(body)

	var quiz model.Quiz
	err = decoder.Decode(&quiz)
	if err != nil {
		return []byte{}, err
	}

	if !qc.service.CheckIfQuizExist(quiz.Id) {
		return []byte{}, &QuizIsNotExistError{}
	}

	err = qc.service.SaveQuiz(quiz)
	if err != nil {
		return []byte{}, err
	}

	bytes, err = json.Marshal(quiz)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func (qc *QuizController) DeleteQuiz(body io.ReadCloser) (bytes []byte, err error) {
	decoder := json.NewDecoder(body)

	var parsedId string
	err = decoder.Decode(&parsedId)
	if err != nil {
		return []byte{}, err
	}

	err = qc.service.DeleteQuiz(parsedId)
	if err != nil {
		return []byte{}, err
	}

	bytes, err = json.Marshal(parsedId)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
