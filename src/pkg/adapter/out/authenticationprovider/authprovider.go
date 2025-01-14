package authenticationprovider

import "net/http"

type AuthenticationProvider interface {
	RetrieveJWT(username string) (string, error)
}

type AuthProvider struct {
}

func (a *AuthProvider) RetrieveJWT(username string) (string, error) {
	// call remote service to retrieve JWT
	req, err := http.NewRequest("GET", "http://localhost:8080/auth", nil)
	if err != nil {
		return "",
			err
	}
	req.Header.Set("username", username)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "",
			err
	}
	defer resp.Body.Close()
	return resp.Header.Get("Authorization"), nil

}
