package handler

import (
	"context"
	"database/sql"
	"errors"

	"com.tom-ludwig/go-server-template/internal/api/users"
	"com.tom-ludwig/go-server-template/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// compile-time check
var _ users.StrictServerInterface = (*UserHandler)(nil)

type UserHandler struct {
	Queries *repository.Queries
}

func NewUserHandler(queries *repository.Queries) *UserHandler {
	return &UserHandler{
		Queries: queries,
	}
}

func (u *UserHandler) GetUser(ctx context.Context, request users.GetUserRequestObject) (users.GetUserResponseObject, error) {
	userUUID, err := uuid.Parse(request.Params.UserId)
	if err != nil {
		return users.GetUser400JSONResponse{}, nil
	}
	user, err := u.Queries.GetUser(ctx, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users.GetUser404JSONResponse{}, nil
		}

		return users.GetUser404JSONResponse{}, nil
		// return users.GetUser500JSONResponse{}, nil
	}
	return users.GetUser200JSONResponse{
		UserId:    user.UserID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email.String,
	}, nil
}

func (u *UserHandler) CreateUser(ctx context.Context, request users.CreateUserRequestObject) (users.CreateUserResponseObject, error) {
	newUser, err := u.Queries.CreateUser(ctx, repository.CreateUserParams{
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
		Email:     pgtype.Text{String: request.Body.Email, Valid: true},
	})

	if err != nil {
		return users.CreateUser500JSONResponse{}, nil
	}

	return users.CreateUser201JSONResponse{
		UserId:    newUser.UserID.String(),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email.String,
	}, nil
}
