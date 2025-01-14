package authenticationprovider

import (
	"net/http"
	"os"
)

type AuthenticationProvider interface {
	RetrieveJWT(username string) (string, error)
}

type AuthProvider struct {
}

func (a *AuthProvider) RetrieveJWT(username string) (string, error) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	// call remote service to retrieve JWT
	req, err := http.NewRequest("GET", authServiceUrl, nil)
	if err != nil {
		return "",
			err
	}
	req.Header.Set("user_id", username)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "",
			err
	}
	defer resp.Body.Close()
	return resp.Header.Get("Authorization"), nil

}
