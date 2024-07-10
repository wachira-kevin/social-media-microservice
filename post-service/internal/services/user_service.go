package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"post-service/config"
	"post-service/internal/models"
)

type userService interface {
	GetUserById(userID uint) error
}

type UserService struct {
	Config *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{Config: cfg}
}

func (us *UserService) GetUserById(userID uint) (*models.User, error) {
	url := fmt.Sprintf("%s/users/%v", us.Config.UserServiceURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Log the response body
	log.Printf("Response Body: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("notfound")
		}
		return nil, errors.New("could not fetch user information")
	}

	var user models.User
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
