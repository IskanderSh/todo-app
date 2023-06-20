package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	todo "todo-app"
	"todo-app/pkg/service"
	mock_service "todo-app/pkg/service/mocks"
)

func TestList_CreateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId int, list todo.TodoList)

	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		inputBody           string
		userId              int
		list                todo.TodoList
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerName:  "userId",
			headerValue: "2",
			userId:      2,
			inputBody:   `{"title": "test title", "description": "test description"}`,
			list: todo.TodoList{
				Title:       "test title",
				Description: "test description",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId int, list todo.TodoList) {
				s.EXPECT().CreateList(userId, list).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, testCase.userId, testCase.list)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/api/lists", handler.createList)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/lists", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(testCase.headerName, testCase.headerValue)

			//ctx := req.Context()
			//ctx = context.WithValue(ctx, testCase.headerName, testCase.headerValue)
			//req = req.WithContext(ctx)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestList_GetAllLists(t *testing.T) {

}

func TestList_GetListById(t *testing.T) {

}

func TestList_UpdateList(t *testing.T) {

}

func TestList_DeleteList(t *testing.T) {

}
