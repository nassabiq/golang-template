package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nassabiq/golang-template/internal/modules/{{MODULE}}/dto"
	"github.com/nassabiq/golang-template/internal/modules/{{MODULE}}/usecase"
	"github.com/nassabiq/golang-template/internal/shared/common/response"
	"github.com/nassabiq/golang-template/internal/shared/helper"
	middleware "github.com/nassabiq/golang-template/internal/shared/middleware/auth"
	proto "github.com/nassabiq/golang-template/proto/{{MODULE}}"
)

type {{MODULE}}Handler struct {
	proto.Unimplemented{{MODULE}}ServiceServer
	usecase usecase.{{MODULE}}Usecase
}

func New{{MODULE}}Handler(usecase usecase.{{MODULE}}Usecase) *{{MODULE}}Handler {
	return &{{MODULE}}Handler{
		usecase: usecase,
	}
}

func (handler *{{MODULE}}Handler) GetByID(ctx context.Context, req *proto.GetByIDRequest) (*proto.{{MODULE}}Response, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.{{MODULE}}Response{
			Metadata: response.Forbidden(),
		}, nil
	}

	{{MODULE|lower}}, err := handler.usecase.GetByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.{{MODULE}}Response{
				Metadata: response.NotFound("{{MODULE|lower}} not found"),
			}, nil
		}

		return &proto.{{MODULE}}Response{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.{{MODULE}}Response{
		Metadata: response.Success(200, "success"),
		Data: &proto.{{MODULE}}{
			Id:   {{MODULE|lower}}.ID,
			Name: {{MODULE|lower}}.Name,
		},
	}, nil
}

func (handler *{{MODULE}}Handler) List(ctx context.Context, req *proto.List{{MODULE}}Request) (*proto.List{{MODULE}}Response, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.List{{MODULE}}Response{
			Metadata: response.Forbidden(),
			Data:     []*proto.{{MODULE}}{},
		}, nil
	}

	{{MODULE|lower}}s, err := handler.usecase.List(ctx, int(req.Limit), int(req.Offset))

	if err != nil {
		return &proto.List{{MODULE}}Response{
			Metadata: response.Internal(),
			Data:     []*proto.{{MODULE}}{},
		}, nil
	}

	resp := &proto.List{{MODULE}}Response{
		Metadata: response.Success(200, "success"),
	}

	for _, item := range {{MODULE|lower}}s {
		resp.Data = append(resp.Data, &proto.{{MODULE}}{
			Id:   item.ID,
			Name: item.Name,
		})
	}

	return resp, nil
}

func (handler *{{MODULE}}Handler) Create(ctx context.Context, req *proto.Create{{MODULE}}Request) (*proto.{{MODULE}}Response, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.{{MODULE}}Response{
			Metadata: response.Forbidden(),
		}, nil
	}

	request := &dto.Create{{MODULE}}Dto{
		Name: req.Name,
	}

	if err := helper.Validate.Struct(request); err != nil {
		return &proto.{{MODULE}}Response{
			Metadata: response.Validation(err.Error()),
		}, nil
	}

	{{MODULE|lower}}, err := handler.usecase.Create(ctx, request)
	if err != nil {
		return &proto.{{MODULE}}Response{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.{{MODULE}}Response{
		Metadata: response.Success(200, "success"),
		Data: &proto.{{MODULE}}{
			Id:   {{MODULE|lower}}.ID,
			Name: {{MODULE|lower}}.Name,
		},
	}, nil
}

func (handler *{{MODULE}}Handler) Update(ctx context.Context, req *proto.Update{{MODULE}}Request) (*proto.{{MODULE}}Response, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.{{MODULE}}Response{
			Metadata: response.Forbidden(),
		}, nil
	}

	updateDto := &dto.Update{{MODULE}}Dto{
		ID: req.GetId(),
	}

	if req.Name != nil {
		name := req.GetName()
		updateDto.Name = &name
	}

	if err := helper.Validate.Struct(updateDto); err != nil {
		return &proto.{{MODULE}}Response{
			Metadata: response.Validation(err.Error()),
		}, nil
	}

	{{MODULE|lower}}, err := handler.usecase.Update(ctx, updateDto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.{{MODULE}}Response{
				Metadata: response.NotFound("{{MODULE|lower}} not found"),
			}, nil
		}
		return &proto.{{MODULE}}Response{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.{{MODULE}}Response{
		Metadata: response.Success(200, "success"),
		Data: &proto.{{MODULE}}{
			Id:   {{MODULE|lower}}.ID,
			Name: {{MODULE|lower}}.Name,
		},
	}, nil
}

func (handler *{{MODULE}}Handler) Delete(ctx context.Context, req *proto.Delete{{MODULE}}Request) (*proto.Delete{{MODULE}}Response, error) {
	if err := middleware.RequireRole("admin", "super_admin")(ctx); err != nil {
		return &proto.Delete{{MODULE}}Response{
			Metadata: response.Forbidden(),
		}, nil
	}

	{{MODULE|lower}}, err := handler.usecase.GetByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.Delete{{MODULE}}Response{
				Metadata: response.NotFound("{{MODULE|lower}} not found"),
			}, nil
		}
		return &proto.Delete{{MODULE}}Response{
			Metadata: response.Internal(),
		}, nil
	}

	err = handler.usecase.Delete(ctx, {{MODULE|lower}})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.Delete{{MODULE}}Response{
				Metadata: response.NotFound("{{MODULE|lower}} not found"),
			}, nil
		}
		return &proto.Delete{{MODULE}}Response{
			Metadata: response.Internal(),
		}, nil
	}

	return &proto.Delete{{MODULE}}Response{
		Metadata: response.Success(200, "success"),
	}, nil
}