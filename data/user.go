package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/briskt/go-htmx-app/data/sqlc"
)

type User struct {
	sqlc.User 
}

type UserCreateInput struct {
	EmployeeID  string
	FirstName   string
	LastName    string
	DisplayName string
	Username    string
	Email       string
}

func (u User) GetEmail() string {
	return u.Email
}

// Delete a user
func (u User) Delete(ctx context.Context, tx sqlc.DBTX) error {
	return q(tx).DeleteUser(ctx, u.ID)
}

// GetDisplayName returns the DisplayName field if it is non-empty, otherwise the FirstName and LastName concatenated.
func (u User) GetDisplayName() string {
	if u.User.DisplayName != "" {
		return u.User.DisplayName
	}
	return strings.TrimSpace(u.User.FirstName + " " + u.User.LastName)
}

func FindUserByUsernameOrEmail(ctx context.Context, tx sqlc.DBTX, v string) (User, error) {
	user, err := q(tx).FindUserByUsername(ctx, v)
	if err != nil {
		user, err = q(tx).FindUserByEmail(ctx, v)
		if err != nil {
			return User{}, fmt.Errorf("no user found with username or email matching %q: %w", v, err)
		}
	}
	dataUser, err := loadUserRelations(ctx, tx, User{User: user})
	if err != nil {
		return User{}, fmt.Errorf("failed to load user relations %q: %w", user.ID, err)
	}
	return dataUser, nil
}

func FindUserByEmployeeID(ctx context.Context, tx sqlc.DBTX, employeeID string) (User, error) {
	user, err := q(tx).FindUserByEmployeeID(ctx, employeeID)
	if err != nil {
		return User{}, fmt.Errorf("no user found with employeeID %q: %w", employeeID, err)
	}
	dataUser, err := loadUserRelations(ctx, tx, User{User: user})
	if err != nil {
		return User{}, fmt.Errorf("failed to load user relations %q: %w", user.ID, err)
	}
	return dataUser, nil
}

func GetUser(ctx context.Context, tx sqlc.DBTX, id int) (User, error) {
	user, err := q(tx).GetUser(ctx, int32(id))
	if err != nil {
		return User{}, fmt.Errorf("no user found with id %q: %w", id, err)
	}
	dataUser, err := loadUserRelations(ctx, tx, User{User: user})
	if err != nil {
		return User{}, fmt.Errorf("failed to load user relations %q: %w", id, err)
	}
	return dataUser, nil
}

// ListActiveUnlockedUsers returns a list of users that are active but not locked
func ListActiveUnlockedUsers(ctx context.Context, tx sqlc.DBTX) ([]User, error) {
	users, err := q(tx).ListActiveUnlockedUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list active, unlocked users: %w", err)
	}
	return toDataUsers(ctx, tx, users, true)
}

// UpdateUserLastLoggedIn sets the user's last_login_utc timestamp to the current time
func UpdateUserLastLoggedIn(ctx context.Context, tx sqlc.DBTX, u User) (User, error) {
	if err := q(tx).UpdateUserLastLoggedIn(ctx, u.ID); err != nil {
		return User{}, fmt.Errorf("failed to update user LastLoggedIn, employeeID=%s: %w", u.EmployeeID, err)
	}
	return GetUser(ctx, tx, int(u.ID))
}

func (u User) Update(ctx context.Context, tx sqlc.DBTX) error {
	return q(tx).UpdateUser(ctx, sqlc.UpdateUserParams{
		EmployeeID:  u.EmployeeID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		DisplayName: u.GetDisplayName(),
		Username:    u.Username,
		Email:       u.Email,
		Active:      u.Active,
		Locked:      u.Locked,
		ID:          u.ID,
	})
}

func CreateUser(ctx context.Context, tx sqlc.DBTX, input UserCreateInput) (User, error) {
	user, err := q(tx).CreateUser(ctx, sqlc.CreateUserParams{
		EmployeeID:  input.EmployeeID,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		DisplayName: input.DisplayName,
		Username:    input.Username,
		Email:       input.Email,
	})
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return User{user}, nil
}

func toDataUsers(ctx context.Context, tx sqlc.DBTX, users []sqlc.User, loadRelations bool) ([]User, error) {
	out := make([]User, len(users))
	for i, u := range users {
		if !loadRelations {
			out[i] = User{User: u}
			continue
		}
		loaded, err := loadUserRelations(ctx, tx, User{User: u})
		if err != nil {
			return nil, fmt.Errorf("failed to load user relations %q: %w", u.ID, err)
		}
		out[i] = loaded

	}
	return out, nil
}

func loadUserRelations(ctx context.Context, tx sqlc.DBTX, user User) (User, error) {
	return user, nil
}
