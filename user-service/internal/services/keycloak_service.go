package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/user-service/config"
	"github.com/user-service/internal/models"
	"io"
	"net/http"
	"strings"
)

type keycloakService interface {
	CreateUserInKeycloak(userInput *models.UserCreationSchema) (string, error)
}

type KeycloakService struct {
	Config *config.Config
}

func NewKeycloakService(cfg *config.Config) *KeycloakService {
	return &KeycloakService{Config: cfg}
}

func (s *KeycloakService) CreateUserInKeycloak(userInput *models.UserCreationSchema) (string, error) {
	url := fmt.Sprintf("%s/admin/realms/%s/users", s.Config.KeycloakURL, s.Config.KeycloakRealm)

	keycloakUser := map[string]interface{}{
		"username":  userInput.Username,
		"email":     userInput.Email,
		"enabled":   true,
		"firstName": userInput.FirstName,
		"lastName":  userInput.LastName,
	}

	jsonData, err := json.Marshal(keycloakUser)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.getAccessToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("failed to create user in Keycloak")
	}

	location := resp.Header.Get("Location")
	keycloakID := location[strings.LastIndex(location, "/")+1:]

	return keycloakID, nil
}

func (s *KeycloakService) getAccessToken() string {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.Config.KeycloakURL, s.Config.KeycloakRealm)
	data := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials",
		s.Config.KeycloakClientID, s.Config.KeycloakClientSecret)

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return ""
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return ""
	}

	return result["access_token"].(string)
}
