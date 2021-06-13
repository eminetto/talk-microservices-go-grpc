package main

import (
	"auth/security"
	"auth/user"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/eminetto/talk-microservices-go-grpc/pkg/proto"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"

)

func main() {
	uService := user.NewService()
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.NewLogger(),
	)
	r.Handle("/v1/auth", n.With(
		negroni.Wrap(userAuth(uService)),
	)).Methods("POST", "OPTIONS")

	http.Handle("/", r)
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":8081", //@TODO usar variável de ambiente
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
	netListener := getNetListener(7000)
	gRPCServer := grpc.NewServer()
	grpcService := user.NewService()
	proto.RegisterTokenServiceServer(gRPCServer, grpcService)
	if err := gRPCServer.Serve(netListener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

func getNetListener(port uint) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	return lis
}

func userAuth(uService user.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var param struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		err = uService.ValidateUser(param.Email, param.Password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		var result struct {
			Token string `json:"token"`
		}
		result.Token, err = security.NewToken(param.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		return
	})
}
