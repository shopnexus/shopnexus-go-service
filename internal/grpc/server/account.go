package server

import (
	"context"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service"
)

type AccountServer struct {
	pb.UnimplementedAccountServer
	service *service.AccountService
}

func NewAccountServer(service *service.AccountService) *AccountServer {
	return &AccountServer{service: service}
}

func (s *AccountServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.TokenResponse, error) {
	token, err := s.service.Login(ctx, service.LoginParams{
		Role:     model.RoleUser,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.TokenResponse{Token: token}, nil
}

func (s *AccountServer) LoginAdmin(ctx context.Context, req *pb.LoginAdminRequest) (*pb.TokenResponse, error) {
	token, err := s.service.Login(ctx, service.LoginParams{
		Role:     model.RoleAdmin,
		Username: &req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.TokenResponse{Token: token}, nil
}

func (s *AccountServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.TokenResponse, error) {
	hashedPassword, err := s.service.CreateHash(req.GetPassword())
	if err != nil {
		return nil, err
	}

	token, err := s.service.Register(ctx, model.AccountUser{
		AccountBase: model.AccountBase{
			Username: req.GetUsername(),
			Password: hashedPassword,
		},
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Gender:   convertGender(req.GetGender()),
		FullName: req.GetFullName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.TokenResponse{Token: token}, nil
}

func (s *AccountServer) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest) (*pb.TokenResponse, error) {
	hashedPassword, err := s.service.CreateHash(req.GetPassword())
	if err != nil {
		return nil, err
	}

	token, err := s.service.Register(ctx, model.AccountAdmin{
		AccountBase: model.AccountBase{
			Username: req.GetUsername(),
			Password: hashedPassword,
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.TokenResponse{Token: token}, nil
}

func convertGender(protoGender pb.Gender) model.Gender {
	switch protoGender {
	case pb.Gender_MALE:
		return model.GenderMale
	case pb.Gender_FEMALE:
		return model.GenderFemale
	case pb.Gender_GENDER_UNSPECIFIED:
		panic("gender is unspecified")
	default:
		panic("unknown gender")
	}
}
