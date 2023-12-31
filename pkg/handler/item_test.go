package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	todo "todo-app"
	"todo-app/pkg/service"
	mock_service "todo-app/pkg/service/mocks"
)

func TestItem_CreateItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, userId, listId int, item todo.TodoItem)

	testTable := []struct {
		name                string
		headerValue         string
		userId              int
		listId              int
		item                todo.TodoItem
		itemBody            string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerValue: "12",
			userId:      12,
			listId:      14,
			item: todo.TodoItem{
				Title:       "test title",
				Description: "test description",
				Done:        false,
			},
			itemBody: `{"title":"test title","description":"test description","done":false}`,
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, item todo.TodoItem) {
				s.EXPECT().CreateItem(userId, listId, item).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "No Header",
			mockBehavior:        func(s *mock_service.MockTodoItem, userId, listId int, item todo.TodoItem) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"unauthorized user"}`,
		},
		{
			name:                "Empty Fields",
			headerValue:         "12",
			userId:              12,
			listId:              14,
			itemBody:            `{"title":"","description":"test description","done":false}`,
			mockBehavior:        func(s *mock_service.MockTodoItem, userId, listId int, item todo.TodoItem) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:        "Service Failure",
			headerValue: "12",
			userId:      12,
			listId:      14,
			item: todo.TodoItem{
				Title:       "test title",
				Description: "test description",
				Done:        false,
			},
			itemBody: `{"title":"test title","description":"test description","done":false}`,
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, item todo.TodoItem) {
				s.EXPECT().CreateItem(userId, listId, item).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mock_service.NewMockTodoItem(c)
			testCase.mockBehavior(todoItem, testCase.userId, testCase.listId, testCase.item)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/api/lists/:id/items", handler.createItem)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/lists/%d/items", testCase.listId), bytes.NewBufferString(testCase.itemBody))
			req.AddCookie(&http.Cookie{Name: userCtx, Value: testCase.headerValue})

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestItem_GetAllItems(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, userId, listId int, items []todo.TodoItem)

	testTable := []struct {
		name                string
		headerValue         string
		userId              int
		listId              int
		output              []todo.TodoItem
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerValue: "19",
			userId:      19,
			listId:      3,
			output: []todo.TodoItem{
				{1, "title1", "description1", true},
				{2, "title2", "description2", false},
				{3, "title3", "description3", true},
			},
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, output []todo.TodoItem) {
				s.EXPECT().GetAll(userId, listId).Return(output, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"id":1,"title":"title1","description":"description1","done":true},{"id":2,"title":"title2","description":"description2","done":false},{"id":3,"title":"title3","description":"description3","done":true}]`,
		},
		{
			name:                "No Header",
			mockBehavior:        func(s *mock_service.MockTodoItem, userId, listId int, output []todo.TodoItem) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"unauthorized user"}`,
		},
		{
			name:        "Service Failure",
			headerValue: "19",
			userId:      19,
			listId:      3,
			output:      []todo.TodoItem{},
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, output []todo.TodoItem) {
				s.EXPECT().GetAll(userId, listId).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mock_service.NewMockTodoItem(c)
			testCase.mockBehavior(todoItem, testCase.userId, testCase.listId, testCase.output)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/api/lists/:id/items", handler.getAllItems)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%d/items", testCase.listId), nil)
			req.AddCookie(&http.Cookie{Name: userCtx, Value: testCase.headerValue})

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestItem_GetItemById(t *testing.T) {

}

func TestItem_UpdateItem(t *testing.T) {

}

func TestItem_DeleteItem(t *testing.T) {

}
