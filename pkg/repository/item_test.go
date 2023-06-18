package repository

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
	todo "todo-app"
)

func TestItem_CreateItem(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemRepository(db)

	type args struct {
		listId int
		item   todo.TodoItem
	}
	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").
					WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Insert Error",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").
					WithArgs(args.listId, id).WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.CreateItem(testCase.args.listId, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestItem_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemRepository(db)

	type args struct {
		userId int
		listId int
	}

	testTable := []struct {
		name         string
		mockBehavior func()
		input        args
		want         []todo.TodoItem
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true).
					AddRow(2, "title2", "description2", false).
					AddRow(3, "title3", "description3", false)
				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li ON (.+) INNER JOIN users_lists ul ON (.+) WHERE (.+) AND (.+)").
					WillReturnRows(rows)
			},
			input: args{
				userId: 1,
				listId: 2,
			},
			want: []todo.TodoItem{
				{1, "title1", "description1", true},
				{2, "title2", "description2", false},
				{3, "title3", "description3", false},
			},
		},
		{
			name: "No Records",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li ON (.+) " +
					"INNER JOIN users_lists ul ON (.+) WHERE (.+)").WillReturnRows(rows)
			},
			input: args{
				userId: 1,
				listId: 2,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err := r.GetAll(testCase.input.userId, testCase.input.listId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}

func TestItem_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("an error on testing TestItem_getById: %v", err)
	}
	defer db.Close()

	r := NewTodoItemRepository(db)

	type args struct {
		userId int
		itemId int
	}

	testTable := []struct {
		name         string
		mockBehavior func()
		input        args
		want         todo.TodoItem
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li ON (.+) " +
					"INNER JOIN users_lists ul ON (.+) WHERE (.+)").WillReturnRows(rows)
			},
			input: args{
				userId: 1,
				itemId: 5,
			},
			want: todo.TodoItem{1, "title1", "description1", true},
		},
		{
			name: "Not Found",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li ON (.+) " +
					"INNER JOIN users_lists ul ON (.+) WHERE (.+)").WillReturnRows(rows)
			},
			input: args{
				userId: 1,
				itemId: 5,
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err := r.GetById(testCase.input.userId, testCase.input.itemId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}

func TestItem_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("an error on testing TestItem_Update: %v", err)
	}
	defer db.Close()

	r := NewTodoItemRepository(db)

	type args struct {
		userId int
		itemId int
		input  todo.UpdateItemInput
	}

	testTable := []struct {
		name         string
		mockBehavior func()
		args         args
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId: 1,
				itemId: 2,
				input: todo.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
					Done:        boolPointer(true),
				},
			},
			mockBehavior: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", true, 1, 2).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Without Done",
			args: args{
				userId: 1,
				itemId: 2,
				input: todo.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
			mockBehavior: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", 1, 2).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Without Done And Description",
			args: args{
				userId: 1,
				itemId: 2,
				input: todo.UpdateItemInput{
					Title: stringPointer("new title"),
				},
			},
			mockBehavior: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", 1, 2).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "No Input Fields",
			args: args{
				userId: 1,
				itemId: 2,
			},
			mockBehavior: func() {
				mock.ExpectExec("UPDATE todo_items ti SET FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err := r.Update(testCase.args.userId, testCase.args.itemId, testCase.args.input)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}

func TestItem_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("error on testing TestItem_Delete func: %v", err)
	}
	defer db.Close()

	r := NewTodoItemRepository(db)

	type args struct {
		userId int
		itemId int
	}

	testTable := []struct {
		name         string
		args         args
		mockBehavior func()
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId: 2,
				itemId: 7,
			},
			mockBehavior: func() {
				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_lists ul WHERE (.+)").
					WithArgs(2, 7).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Not Found",
			args: args{
				userId: 2,
				itemId: 7,
			},
			mockBehavior: func() {
				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_lists ul WHERE (.+)").
					WithArgs(2, 7).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err := r.Delete(testCase.args.userId, testCase.args.itemId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
