package services

import (
	"auth-service/config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type KeycloakService struct {
	Config *config.Config
}

func NewKeycloakService(cfg *config.Config) *KeycloakService {
	return &KeycloakService{Config: cfg}
}

func (ks *KeycloakService) Login(username, password string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", ks.Config.KeycloakURL, ks.Config.KeycloakRealm)
	data := fmt.Sprintf("client_id=%s&client_secret=%s&username=%s&password=%s&grant_type=password",
		ks.Config.KeycloakClientID, ks.Config.KeycloakClientSecret, username, password)
	return ks.postRequest(url, data)
}

func (ks *KeycloakService) Logout(sessionId string) error {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", ks.Config.KeycloakURL, ks.Config.KeycloakRealm)
	data := fmt.Sprintf("client_id=%s&refresh_token=%s&client_secret=%s",
		ks.Config.KeycloakClientID, sessionId, ks.Config.KeycloakClientSecret)
	_, err := ks.postRequest(url, data)
	return err
}

func (ks *KeycloakService) RefreshToken(refreshToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", ks.Config.KeycloakURL, ks.Config.KeycloakRealm)
	data := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s",
		ks.Config.KeycloakClientID, ks.Config.KeycloakClientSecret, refreshToken)
	return ks.postRequest(url, data)
}

func (ks *KeycloakService) postRequest(url string, data string) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak returned error: %v", result)
	}

	return result, nil
}
