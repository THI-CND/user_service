package authenticationprovider

import (
	"context"
	pb "github.com/BieggerM/userservice/proto/github.com/BieggerM/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

type AuthenticationProvider interface {
	RetrieveJWT(username string) (string, error)
}

type AuthProvider struct {
}

type GRPCAuthProvider struct {
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

func (a *GRPCAuthProvider) RetrieveJWT(username string) (string, error) {
	conn, err := grpc.NewClient(
		"localhost:50081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Call the Auth method
	authResponse, err := client.Auth(context.Background(), &pb.AuthRequest{UserId: "user123"})
	if err != nil {
		log.Fatalf("Could not authenticate: %v", err)
	}

	// Call the Protected method
	protectedResponse, err := client.Protected(context.Background(), &pb.ProtectedRequest{Token: authResponse.Token})
	if err != nil {
		log.Fatalf("Could not access protected resource: %v", err)
	}
	log.Printf("Protected message: %s", protectedResponse.Message)

	return authResponse.Token, nil
}
