package handler

import (
	"context"

	"com.tom-ludwig/go-server-template/internal/api"
	"com.tom-ludwig/go-server-template/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	Queries *repository.Queries
}

func (u *UserHandler) GetUser(ctx context.Context, request api.GetUserRequestObject) (api.GetUserResponseObject, error) {
	userUUID, err := uuid.Parse(request.Params.UserId)
	if err != nil {
		return api.GetUser400JSONResponse{
			Message: "Invalid user ID",
		}, nil
	}
	user, err := u.Queries.GetUser(ctx, userUUID)
	if err != nil {
		return api.GetUser404JSONResponse{Message: "User not found"}, nil
	}
	return api.GetUser200JSONResponse{
		UserId:    user.UserID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email.String,
	}, nil
}

func (u *UserHandler) CreateUser(ctx context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	newUser, err := u.Queries.CreateUser(ctx, repository.CreateUserParams{
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
		Email:     pgtype.Text{String: request.Body.Email, Valid: true},
	})

	if err != nil {
		return api.CreateUser500JSONResponse{
			InternalServerErrrorJSONResponse: api.InternalServerErrrorJSONResponse{
				Message: "Failed to create user",
			},
		}, nil
	}

	return api.CreateUser201JSONResponse{
		UserId:    newUser.UserID.String(),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email.String,
	}, nil
}
