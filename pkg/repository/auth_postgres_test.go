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

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("error on testing TestAuthPostgres_GetUser func: %v", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type args struct {
		username string
		password string
	}

	testTable := []struct {
		name         string
		args         args
		mockBehavior func(args args)
		want         todo.User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				username: "test",
				password: "qwerty",
			},
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"}).
					AddRow("1", "Test", "test", "qwerty")
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").WithArgs(args.username, args.password).
					WillReturnRows(rows)
			},
			want: todo.User{
				Id:       1,
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
		},
		{
			name: "Empty Fields",
			args: args{
				username: "",
				password: "qwerty",
			},
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"})

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").WithArgs(args.username, args.password).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetUser(testCase.args.username, testCase.args.password)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}
