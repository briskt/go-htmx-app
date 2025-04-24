package data

import (
	"time"
)

func (s *Suite) TestUserGetEmail() {
	var user User
	user.Email = "x@y.com"
	s.Equal("x@y.com", user.GetEmail())
}

func (s *Suite) TestUserDisplayName() {
	var user User
	user.FirstName = "a"
	s.Equal("a", user.GetDisplayName())

	user.LastName = "b"
	s.Equal("a b", user.GetDisplayName())

	user.User.DisplayName = "c"
	s.Equal("c", user.GetDisplayName())
}

func (s *Suite) TestFindUserByUsernameOrEmail() {
	_, err := FindUserByUsernameOrEmail(s.ctx, s.db, "foo")
	s.Error(err)
	user := insertUser(s.db)
	got, err := FindUserByUsernameOrEmail(s.ctx, s.db, user.Email)
	s.NoError(err)
	s.Equal(user.ID, got.ID)
	got, err = FindUserByUsernameOrEmail(s.ctx, s.db, user.Username)
	s.NoError(err)
	s.Equal(user.ID, got.ID)
}

func (s *Suite) TestFindUserByEmployeeID() {
	user := insertUser(s.db)

	tests := []struct {
		name       string
		employeeID string
		wantEmail  string
		wantErr    bool
	}{
		{
			name:       "return err when not found",
			employeeID: "1",
			wantErr:    true,
		},
		{
			name:       "return data when found",
			employeeID: user.EmployeeID,
			wantEmail:  user.Email,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := FindUserByEmployeeID(s.ctx, s.db, tt.employeeID)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.Equal(tt.wantEmail, got.Email)
		})
	}
}

func (s *Suite) TestGetUser() {
	user := insertUser(s.db)

	tests := []struct {
		name      string
		id        int
		wantEmail string
		wantErr   bool
	}{
		{
			name:    "return err when not found",
			id:      0,
			wantErr: true,
		},
		{
			name:      "return data when found",
			id:        int(user.ID),
			wantEmail: user.Email,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := GetUser(s.ctx, s.db, tt.id)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.Equal(tt.wantEmail, got.Email)
		})
	}
}

func (s *Suite) TestUpdateUserLastLoggedIn() {
	user := insertUser(s.db)

	got, err := UpdateUserLastLoggedIn(s.ctx, s.db, User{User: user})
	s.NoError(err)
	s.WithinDuration(time.Now().UTC(), got.LastLoginAt, time.Second)
}

func (s *Suite) TestUserUpdate() {
	user := insertUser(s.db)
	user.FirstName = "Frances"

	err := User{User: user}.Update(s.ctx, s.db)
	s.NoError(err)

	got, err := GetUser(s.ctx, s.db, int(user.ID))
	s.NoError(err)
	s.Equal(user.FirstName, got.FirstName)
}

func (s *Suite) TestLoadUserRelations() {
	user := insertUser(s.db)

	got, err := loadUserRelations(s.ctx, s.db, User{User: user})
	s.NoError(err)
	s.Equal(user.ID, got.ID)
}
