package user

import (
	"context"
	"fmt"

	"auth/security"
	"github.com/eminetto/talk-microservices-go-grpc/pkg/proto"
)

type UseCase interface {
	ValidateUser(email, password string) error
}

type Service struct{}


func NewService() *Service {
	return &Service{}
}
func (s *Service) ValidateUser(email, password string) error {
	//@TODO create validation rules, using databases or something else
	if email == "eminetto@gmail.com" && password != "1234567" {
		return fmt.Errorf("Invalid user")
	}
	return nil
}

func (s *Service) Validate(ctx context.Context, token *proto.Token) (*proto.User, error) {
	t, err := security.ParseToken(token.Token)
	if err != nil {
		return &proto.User{}, err
	}
	tData, err := security.GetClaims(t)
	if err != nil {
		return &proto.User{}, err
	}
	return &proto.User{
		IsValid: true,
		Email:   tData["email"].(string),
	}, nil
}
func (s *Service) mustEmbedUnimplementedTokenServiceServer() {
	panic("implement me")
}
