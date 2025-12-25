package handler

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

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
		slog.Error(
			"An error occured while trying to create a user",
			"Error: ", err,
		)
		return users.CreateUser500JSONResponse{}, nil
	}

	return users.CreateUser201JSONResponse{
		UserId:    newUser.UserID.String(),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email.String,
	}, nil
}

func (u *UserHandler) GetUsers(ctx context.Context, request users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	// Set default values for pagination parameters
	page := int32(1)
	if request.Params.Page != nil {
		page = int32(*request.Params.Page)
	}

	limit := int32(20) // Default limit
	if request.Params.Limit != nil {
		limit = int32(*request.Params.Limit)
	}

	// Validate pagination parameters
	if page < 1 || limit < 1 || limit > 100 {
		return users.GetUsers400JSONResponse{
			BadRequestJSONResponse: users.BadRequestJSONResponse{
				Message: "Invalid pagination parameters: page must be >= 1, limit must be between 1 and 100",
			},
		}, nil
	}

	// Get total count of users
	totalRecords, err := u.Queries.CountUsers(ctx)
	if err != nil {
		slog.Error(
			"An error occurred while trying to count users",
			"error", err,
		)
		return users.GetUsers500JSONResponse{
			InternalServerErrorJSONResponse: users.InternalServerErrorJSONResponse{
				Message: "An internal server error occurred",
			},
		}, nil
	}

	// Calculate pagination metadata
	totalPages := int32((totalRecords + int64(limit) - 1) / int64(limit)) // Ceiling division
	if totalPages == 0 {
		totalPages = 1 // At least 1 page even if empty
	}

	var nextPage *int
	if page < totalPages {
		next := int(page + 1)
		nextPage = &next
	}

	var prevPage *int
	if page > 1 {
		prev := int(page - 1)
		prevPage = &prev
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch users with pagination
	queryParameters := repository.GetUsersParams{
		Limit:  limit,
		Offset: offset,
	}

	dbUsers, err := u.Queries.GetUsers(ctx, queryParameters)
	if err != nil {
		slog.Error(
			"An error occurred while trying to get users",
			"error", err,
		)
		return users.GetUsers500JSONResponse{
			InternalServerErrorJSONResponse: users.InternalServerErrorJSONResponse{
				Message: "An internal server error occurred",
			},
		}, nil
	}

	// Convert database users to API users
	var apiUsers []users.User
	for _, dbUser := range dbUsers {
		apiUsers = append(apiUsers, users.User{
			UserId:    dbUser.UserID.String(),
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Email:     dbUser.Email.String,
		})
	}

	return users.GetUsers200JSONResponse{
		Data: apiUsers,
		Pagination: users.PaginationMetadata{
			CurrentPage:  int(page),
			Limit:        int(limit),
			NextPage:     nextPage,
			PrevPage:     prevPage,
			TotalPages:   int(totalPages),
			TotalRecords: int(totalRecords),
		},
	}, nil
}
