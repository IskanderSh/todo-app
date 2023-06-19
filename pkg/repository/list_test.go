package repository

import (
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
	todo "todo-app"
)

func TestList_CreateList(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("error on testing TestList_CreateList func: %v", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		list   todo.TodoList
	}
	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId: 2,
				list: todo.TodoList{
					Id:          1,
					Title:       "title",
					Description: "description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_lists").WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_lists").WithArgs(args.userId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			id: 5,
		},
		{
			name: "Empty Fields",
			args: args{
				userId: 2,
				list: todo.TodoList{
					Id:          1,
					Title:       "",
					Description: "description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO todo_lists").WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			id:      5,
			wantErr: true,
		},
		{
			name: "2nd Insert Error",
			args: args{
				userId: 2,
				list: todo.TodoList{
					Id:          1,
					Title:       "title",
					Description: "description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_lists").WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_lists").WithArgs(args.userId, id).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			id:      5,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.CreateList(testCase.args.userId, testCase.args.list)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestList_GetAll(t *testing.T) {

}

func TestList_GetById(t *testing.T) {

}

func TestList_Update(t *testing.T) {

}

func TestList_Delete(t *testing.T) {

}