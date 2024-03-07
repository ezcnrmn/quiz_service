package service

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ezcnrmn/quiz_service/internal/app/model"
	"github.com/ezcnrmn/quiz_service/internal/app/utils/consts"
)

type Service struct{}

func New() *Service {
	service := &Service{}
	service.CreateFolder()

	return service
}

func (s *Service) CreateFolder() {
	_ = os.Mkdir(consts.FOLDER_NAME, 0666)
}

func (s *Service) GetQuizzesList() ([]string, error) {
	entries, err := os.ReadDir(consts.FOLDER_NAME)

	quizzes := make([]string, len(entries))
	for i, entry := range entries {
		quizzes[i] = entry.Name()
	}

	return quizzes, err
}

func (s *Service) LoadQuiz(fileName string) (quiz model.Quiz, err error) {
	path := path.Join(consts.FOLDER_NAME, fileName)

	bytes, err := os.ReadFile(path)
	if err != nil {
		return model.Quiz{}, err
	}

	reader := strings.NewReader(string(bytes))
	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&quiz)
	if err != nil {
		return model.Quiz{}, err
	}

	return quiz, nil
}

func (s *Service) LoadAllQuizzes() ([]model.Quiz, error) {
	entries, err := os.ReadDir(consts.FOLDER_NAME)
	if err != nil {
		return []model.Quiz{}, err
	}

	quizzes := make([]model.Quiz, len(entries))

	for i, entry := range entries {
		quiz, err := s.LoadQuiz(entry.Name())
		if err != nil {
			return []model.Quiz{}, err
		}

		quizzes[i] = quiz
	}

	return quizzes, nil
}

func (s *Service) GetNewQuizId() (id string, err error) {
	bytes := make([]byte, 16)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", bytes), nil
}

func (s *Service) SaveQuiz(quiz model.Quiz) (err error) {
	path := path.Join(consts.FOLDER_NAME, quiz.Id+consts.FILE_EXTENSION)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(quiz)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteQuiz(id string) (err error) {
	path := path.Join(consts.FOLDER_NAME, id+consts.FILE_EXTENSION)

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckIfQuizExist(id string) (res bool) {
	path := path.Join(consts.FOLDER_NAME, id+consts.FILE_EXTENSION)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
