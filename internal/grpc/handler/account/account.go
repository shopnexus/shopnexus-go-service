package account

import (
	"context"
	"shopnexus-go-service/internal/grpc/interceptor/auth"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/account"

	"connectrpc.com/connect"
	accountv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
)

type ImplementedAccountServiceHandler struct {
	accountv1connect.UnimplementedAccountServiceHandler
	service *account.AccountService
}

func NewAccountServiceHandler(service *account.AccountService) accountv1connect.AccountServiceHandler {
	return &ImplementedAccountServiceHandler{service: service}
}

func (s *ImplementedAccountServiceHandler) GetUser(ctx context.Context, req *connect.Request[accountv1.GetUserRequest]) (*connect.Response[accountv1.GetUserResponse], error) {
	claims, err := auth.GetAccount(req)
	if err != nil {
		return nil, err
	}

	account, err := s.service.FindAccount(ctx, account.FindAccountParams{
		Role:   model.RoleUser,
		UserID: &claims.UserID,
	})
	if err != nil {
		return nil, err
	}

	user, ok := account.(model.AccountUser)
	if !ok {
		return nil, model.ErrForbidden
	}

	return connect.NewResponse(&accountv1.GetUserResponse{
		Email:            user.Username,
		Phone:            user.Phone,
		Username:         user.Username,
		Gender:           convertGenderToProto(user.Gender),
		FullName:         user.FullName,
		DefaultAddressId: user.DefaultAddressID,
	}), nil
}

func (s *ImplementedAccountServiceHandler) LoginUser(ctx context.Context, req *connect.Request[accountv1.LoginUserRequest]) (*connect.Response[accountv1.LoginUserResponse], error) {
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

func (s *ImplementedAccountServiceHandler) LoginAdmin(ctx context.Context, req *connect.Request[accountv1.LoginAdminRequest]) (*connect.Response[accountv1.LoginAdminResponse], error) {
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

func (s *ImplementedAccountServiceHandler) RegisterUser(ctx context.Context, req *connect.Request[accountv1.RegisterUserRequest]) (*connect.Response[accountv1.RegisterUserResponse], error) {
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

func (s *ImplementedAccountServiceHandler) RegisterAdmin(ctx context.Context, req *connect.Request[accountv1.RegisterAdminRequest]) (*connect.Response[accountv1.RegisterAdminResponse], error) {
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

func convertGenderToProto(gender model.Gender) accountv1.Gender {
	switch gender {
	case model.GenderMale:
		return accountv1.Gender_GENDER_MALE
	case model.GenderFemale:
		return accountv1.Gender_GENDER_FEMALE
	case model.GenderOther:
		return accountv1.Gender_GENDER_OTHER
	default:
		return accountv1.Gender_GENDER_UNSPECIFIED
	}
}
