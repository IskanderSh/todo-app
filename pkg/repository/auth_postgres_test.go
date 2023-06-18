package repository

import (
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
	todo "todo-app"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("error on testing TestAuthPostgres_CreateUser func: %v", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	testTable := []struct {
		name         string
		user         todo.User
		id           int
		mockBehavior func(user todo.User, id int)
		wantErr      bool
	}{
		{
			name: "OK",
			user: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			id: 4,
			mockBehavior: func(user todo.User, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO users").WithArgs(user.Name, user.Username, user.Password).
					WillReturnRows(rows)
			},
		},
		{
			name: "Empty Fields",
			user: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "",
			},
			id: 4,
			mockBehavior: func(user todo.User, id int) {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO users").WithArgs(user.Name, user.Username, user.Password).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.user, testCase.id)

			got, err := r.CreateUser(testCase.user)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}
