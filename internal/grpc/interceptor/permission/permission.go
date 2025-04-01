package permission

import (
	"context"
	"errors"
	"shopnexus-go-service/internal/grpc/interceptor/auth"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/account"

	"slices"

	"connectrpc.com/connect"
)

type Options struct {
	Roles       []model.Role
	Permissions []model.Permission
}

type Option func(*Options)

func NeedRoles(roles ...model.Role) Option {
	return func(ro *Options) {
		ro.Roles = roles
	}
}

func NeedPermissions(permissions ...model.Permission) Option {
	return func(ro *Options) {
		ro.Permissions = permissions
	}
}

func UseOptions(opts ...Option) Options {
	var ro Options
	for _, opt := range opts {
		opt(&ro)
	}
	return ro
}

func NewPermissionInterceptor(
	accountSvc *account.AccountService,
	routes map[string]Options,
) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			method := req.Spec().Procedure

			// Check if method requires permission validation
			opts, exists := routes[method]
			if !exists || len(opts.Roles) == 0 {
				return next(ctx, req)
			}

			// Perform permission validation
			if err := validatePermissions(ctx, req, accountSvc, opts); err != nil {
				return nil, err
			}

			return next(ctx, req)
		}
	}
}

// validatePermissions performs the complete permission validation process
func validatePermissions(ctx context.Context, req connect.AnyRequest, accountSvc *account.AccountService, opts Options) error {
	claims, err := auth.GetAccount(req)
	if err != nil {
		return err
	}

	// Check roles
	if err := checkUserRole(claims, opts.Roles); err != nil {
		return err
	}

	// Check permissions if needed
	if len(opts.Permissions) > 0 {
		if err := checkUserPermissions(ctx, accountSvc, claims, opts.Permissions); err != nil {
			return err
		}
	}

	return nil
}

// checkUserRole verifies if the user has one of the required roles
func checkUserRole(claims model.Claims, requiredRoles []model.Role) error {
	if !slices.Contains(requiredRoles, claims.Role) {
		return connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}
	return nil
}

// checkUserPermissions verifies if the user has all required permissions
func checkUserPermissions(
	ctx context.Context,
	accountSvc *account.AccountService,
	claims model.Claims,
	requiredPermissions []model.Permission,
) error {
	permissions, err := accountSvc.GetPermissions(ctx, account.GetPermissionsParams{
		AccountID: claims.UserID,
		Role:      &claims.Role,
	})
	if err != nil {
		return err
	}

	for _, permission := range requiredPermissions {
		if !slices.Contains(permissions, permission) {
			return connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
		}
	}

	return nil
}
