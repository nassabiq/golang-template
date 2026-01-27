package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nassabiq/golang-template/internal/modules/user/dto"
	"github.com/nassabiq/golang-template/internal/modules/user/usecase"
	"github.com/nassabiq/golang-template/internal/shared/common/response"
	"github.com/nassabiq/golang-template/internal/shared/helper"
	middleware "github.com/nassabiq/golang-template/internal/shared/middleware/auth"
	commonpb "github.com/nassabiq/golang-template/proto/common"
	proto "github.com/nassabiq/golang-template/proto/user"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (handler *UserHandler) GetMe(ctx context.Context, _ *proto.Empty) (*proto.UserResponse, error) {
	userID, _, ok := middleware.FromContext(ctx)
	if !ok {
		return &proto.UserResponse{
			Metadata: response.Unauthorized(),
		}, nil
	}

	user, err := handler.usecase.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.UserResponse{
				Metadata: response.NotFound("user not found"),
			}, nil
		}
		return &proto.UserResponse{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.UserResponse{
		Metadata: response.Success(200, "success"),
		Data: &proto.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.RoleID,
		},
	}, nil
}

func (handler *UserHandler) GetByID(ctx context.Context, req *proto.GetByIDRequest) (*proto.UserResponse, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.UserResponse{
			Metadata: response.Forbidden(),
		}, nil
	}

	user, err := handler.usecase.GetByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.UserResponse{
				Metadata: response.NotFound("user not found"),
			}, nil
		}

		return &proto.UserResponse{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.UserResponse{
		Metadata: response.Success(200, "success"),
		Data: &proto.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.RoleID,
		},
	}, nil
}

func (handler *UserHandler) List(ctx context.Context, req *proto.ListUserRequest) (*proto.ListUserResponse, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.ListUserResponse{
			Metadata: response.Forbidden(),
			Users:    []*proto.User{},
		}, nil
	}

	users, total, err := handler.usecase.List(ctx, int(req.Limit), int(req.Offset))

	if err != nil {
		return &proto.ListUserResponse{
			Metadata: response.Internal(),
			Users:    []*proto.User{},
		}, nil
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	resp := &proto.ListUserResponse{
		Metadata: response.Success(200, "success"),
		Pagination: &commonpb.Pagination{
			Limit:  limit,
			Offset: req.Offset,
			Total:  total,
		},
	}

	for _, u := range users {
		resp.Users = append(resp.Users, &proto.User{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.RoleID,
		})
	}

	return resp, nil
}

func (handler *UserHandler) Create(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.UserResponse{
			Metadata: response.Forbidden(),
		}, nil
	}

	request := &dto.CreateUserDto{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleId,
	}

	if err := helper.Validate.Struct(request); err != nil {
		return &proto.UserResponse{
			Metadata: response.Validation(err.Error()),
		}, nil
	}

	user, err := handler.usecase.Create(ctx, request)
	if err != nil {
		return &proto.UserResponse{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.UserResponse{
		Metadata: response.Success(200, "success"),
		Data: &proto.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.RoleID,
		},
	}, nil
}

func (handler *UserHandler) Update(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.UserResponse{
			Metadata: response.Forbidden(),
		}, nil
	}

	updateDto := &dto.UpdateUserDto{
		ID: req.GetId(),
	}

	if req.Name != nil {
		name := req.GetName()
		updateDto.Name = &name
	}

	if req.Email != nil {
		email := req.GetEmail()
		updateDto.Email = &email
	}

	if req.RoleId != nil {
		roleID := req.GetRoleId()
		updateDto.RoleID = &roleID
	}

	if err := helper.Validate.Struct(updateDto); err != nil {
		return &proto.UserResponse{
			Metadata: response.Validation(err.Error()),
		}, nil
	}
	user, err := handler.usecase.Update(ctx, updateDto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.UserResponse{
				Metadata: response.NotFound("user not found"),
			}, nil
		}
		return &proto.UserResponse{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.UserResponse{
		Metadata: response.Success(200, "success"),
		Data: &proto.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.RoleID,
		},
	}, nil
}

func (handler *UserHandler) Delete(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.DeleteUserResponse{
			Metadata: response.Forbidden(),
		}, nil
	}

	user, err := handler.usecase.GetByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.DeleteUserResponse{
				Metadata: response.NotFound("user not found"),
			}, nil
		}
		return &proto.DeleteUserResponse{
			Metadata: response.Internal(),
		}, nil
	}

	err = handler.usecase.Delete(ctx, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.DeleteUserResponse{
				Metadata: response.NotFound("user not found"),
			}, nil
		}
		return &proto.DeleteUserResponse{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.DeleteUserResponse{
		Metadata: response.Success(200, "success"),
	}, nil
}
