package account

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/account"

	"connectrpc.com/connect"
	accountv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
)

var _ accountv1connect.AccountServiceHandler = (*AccountServer)(nil)

type AccountServer struct {
	accountv1connect.UnimplementedAccountServiceHandler
	service *account.AccountService
}

func NewAccountServer(service *account.AccountService) *AccountServer {
	return &AccountServer{service: service}
}

func (s *AccountServer) LoginUser(ctx context.Context, req *connect.Request[accountv1.LoginUserRequest]) (*connect.Response[accountv1.LoginUserResponse], error) {
	token, err := s.service.Login(ctx, account.LoginParams{
		Role:     model.RoleUser,
		Username: req.Msg.Username,
		Email:    req.Msg.Email,
		Phone:    req.Msg.Phone,
		Password: req.Msg.Password,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.LoginUserResponse{Token: token}), nil
}

func (s *AccountServer) LoginAdmin(ctx context.Context, req *connect.Request[accountv1.LoginAdminRequest]) (*connect.Response[accountv1.LoginAdminResponse], error) {
	token, err := s.service.Login(ctx, account.LoginParams{
		Role:     model.RoleAdmin,
		Username: &req.Msg.Username,
		Password: req.Msg.Password,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.LoginAdminResponse{Token: token}), nil
}

func (s *AccountServer) RegisterUser(ctx context.Context, req *connect.Request[accountv1.RegisterUserRequest]) (*connect.Response[accountv1.RegisterUserResponse], error) {
	hashedPassword, err := s.service.CreateHash(req.Msg.GetPassword())
	if err != nil {
		return nil, err
	}

	token, err := s.service.Register(ctx, model.AccountUser{
		AccountBase: model.AccountBase{
			Username: req.Msg.GetUsername(),
			Password: hashedPassword,
		},
		Email:    req.Msg.GetEmail(),
		Phone:    req.Msg.GetPhone(),
		Gender:   convertGender(req.Msg.GetGender()),
		FullName: req.Msg.GetFullName(),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.RegisterUserResponse{Token: token}), nil
}

func (s *AccountServer) RegisterAdmin(ctx context.Context, req *connect.Request[accountv1.RegisterAdminRequest]) (*connect.Response[accountv1.RegisterAdminResponse], error) {
	hashedPassword, err := s.service.CreateHash(req.Msg.GetPassword())
	if err != nil {
		return nil, err
	}

	token, err := s.service.Register(ctx, model.AccountAdmin{
		AccountBase: model.AccountBase{
			Username: req.Msg.GetUsername(),
			Password: hashedPassword,
		},
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.RegisterAdminResponse{Token: token}), nil
}

func convertGender(protoGender accountv1.Gender) model.Gender {
	switch protoGender {
	case accountv1.Gender_GENDER_MALE:
		return model.GenderMale
	case accountv1.Gender_GENDER_FEMALE:
		return model.GenderFemale
	case accountv1.Gender_GENDER_OTHER:
		return model.GenderOther
	case accountv1.Gender_GENDER_UNSPECIFIED:
		panic("gender is unspecified")
	default:
		panic("unknown gender")
	}
}
