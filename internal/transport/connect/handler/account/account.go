package account

import (
	"context"
	"net/http"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/account"
	"shopnexus-go-service/internal/transport/connect/interceptor/auth"
	"shopnexus-go-service/internal/utils/ptr"

	"connectrpc.com/connect"
	accountv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
)

type ImplementedAccountServiceHandler struct {
	accountv1connect.UnimplementedAccountServiceHandler
	service account.Service
}

func NewAccountServiceHandler(service account.Service, opts ...connect.HandlerOption) (string, http.Handler) {
	return accountv1connect.NewAccountServiceHandler(
		&ImplementedAccountServiceHandler{service: service},
		opts...,
	)
}

func (s *ImplementedAccountServiceHandler) GetUser(ctx context.Context, req *connect.Request[accountv1.GetUserRequest]) (*connect.Response[accountv1.GetUserResponse], error) {
	claims, err := auth.GetClaims(req)
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
		Id:               user.ID,
		Email:            user.Email,
		Phone:            user.Phone,
		Username:         user.Username,
		Gender:           convertGenderToProto(user.Gender),
		FullName:         user.FullName,
		DefaultAddressId: user.DefaultAddressID,
		Avatar:           user.Avatar,
	}), nil
}

func (s *ImplementedAccountServiceHandler) GetAdmin(ctx context.Context, req *connect.Request[accountv1.GetAdminRequest]) (*connect.Response[accountv1.GetAdminResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	account, err := s.service.FindAccount(ctx, account.FindAccountParams{
		Role:   model.RoleAdmin,
		UserID: &claims.UserID,
	})
	if err != nil {
		return nil, err
	}

	admin, ok := account.(model.AccountAdmin)
	if !ok {
		return nil, model.ErrForbidden
	}

	return connect.NewResponse(&accountv1.GetAdminResponse{
		Id:       admin.ID,
		Username: admin.Username,
		Role:     convertRoleToProto(admin.Role),
		Avatar:   admin.Avatar,
	}), nil
}

func (s *ImplementedAccountServiceHandler) GetUserPublic(ctx context.Context, req *connect.Request[accountv1.GetUserPublicRequest]) (*connect.Response[accountv1.GetUserPublicResponse], error) {
	account, err := s.service.FindAccount(ctx, account.FindAccountParams{
		Role:   model.RoleUser,
		UserID: &req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	user, ok := account.(model.AccountUser)
	if !ok {
		return nil, model.ErrForbidden
	}

	return connect.NewResponse(&accountv1.GetUserPublicResponse{
		Id:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Avatar:   user.Avatar,
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

func (s *ImplementedAccountServiceHandler) UpdateAccount(ctx context.Context, req *connect.Request[accountv1.UpdateAccountRequest]) (*connect.Response[accountv1.UpdateAccountResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	var (
		accountID int64
	)
	if req.Msg.Id != nil {
		accountID = *req.Msg.Id
	} else {
		accountID = claims.UserID
	}

	// TODO: check permission on update account & update user (add permission too)
	_, err = s.service.UpdateAccount(ctx, account.UpdateAccountParams{
		ID:                   accountID,
		Username:             req.Msg.Username,
		Password:             req.Msg.Password,
		NullCustomPermission: req.Msg.NullCustomPermission,
		CustomPermission:     req.Msg.CustomPermission,
		AvatarURL:            req.Msg.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.UpdateAccountResponse{}), nil
}

func (s *ImplementedAccountServiceHandler) UpdateUser(ctx context.Context, req *connect.Request[accountv1.UpdateUserRequest]) (*connect.Response[accountv1.UpdateUserResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	var (
		accountID int64
		gender    *model.Gender
	)
	if req.Msg.Gender != nil {
		gender = ptr.ToPtr(convertGender(*req.Msg.Gender))
	}
	if req.Msg.Id != nil {
		accountID = *req.Msg.Id
	} else {
		accountID = claims.UserID
	}

	_, err = s.service.UpdateAccountUser(ctx, account.UpdateAccountUserParams{
		ID:               accountID,
		Email:            req.Msg.Email,
		Phone:            req.Msg.Phone,
		Gender:           gender,
		FullName:         req.Msg.FullName,
		DefaultAddressID: req.Msg.DefaultAddressId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.UpdateUserResponse{}), nil
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

func convertRoleToProto(role model.Role) accountv1.Role {
	switch role {
	case model.RoleAdmin:
		return accountv1.Role_ROLE_ADMIN
	case model.RoleUser:
		return accountv1.Role_ROLE_USER
	case model.RoleStaff:
		return accountv1.Role_ROLE_STAFF
	}

	panic("unknown role")
}
