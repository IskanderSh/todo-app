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

func TestList_CreateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId int, list todo.TodoList)

	testTable := []struct {
		name                string
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
		{
			name:                "Empty Fields",
			headerValue:         "2",
			userId:              2,
			inputBody:           `{"title": "", "description": "test description"}`,
			mockBehavior:        func(s *mock_service.MockTodoList, userId int, list todo.TodoList) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "No Header",
			inputBody: `{"title": "test title", "description": "test description"}`,
			list: todo.TodoList{
				Title:       "test title",
				Description: "test description",
			},
			mockBehavior:        func(s *mock_service.MockTodoList, userId int, list todo.TodoList) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"unauthorized user"}`,
		},
		{
			name:        "Service Failure",
			headerValue: "2",
			userId:      2,
			inputBody:   `{"title": "test title", "description": "test description"}`,
			list: todo.TodoList{
				Title:       "test title",
				Description: "test description",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId int, list todo.TodoList) {
				s.EXPECT().CreateList(userId, list).Return(0, errors.New("service failure"))
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
			req.AddCookie(&http.Cookie{Name: userCtx, Value: testCase.headerValue})

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestList_GetAllLists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId int, output []todo.TodoList)

	testTable := []struct {
		name                string
		headerValue         string
		userId              int
		output              []todo.TodoList
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerValue: "4",
			userId:      4,
			output: []todo.TodoList{
				{
					Id:          1,
					Title:       "title1",
					Description: "description1",
				},
				{
					Id:          2,
					Title:       "title2",
					Description: "description2",
				},
				{
					Id:          3,
					Title:       "title3",
					Description: "description3",
				},
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId int, output []todo.TodoList) {
				s.EXPECT().GetAll(userId).Return(output, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":1,"title":"title1","description":"description1"},{"id":2,"title":"title2","description":"description2"},{"id":3,"title":"title3","description":"description3"}]}`,
		},
		{
			name:                "No Header",
			mockBehavior:        func(s *mock_service.MockTodoList, userId int, output []todo.TodoList) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"user unauthorized"}`,
		},
		{
			name:        "Service Failure",
			headerValue: "4",
			userId:      4,
			output:      []todo.TodoList{},
			mockBehavior: func(s *mock_service.MockTodoList, userId int, output []todo.TodoList) {
				s.EXPECT().GetAll(userId).Return(output, errors.New("service failure"))
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

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, testCase.userId, testCase.output)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/api/lists", handler.getAllLists)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/lists", nil)
			req.AddCookie(&http.Cookie{Name: userCtx, Value: testCase.headerValue})

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestList_GetListById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId int, output todo.TodoList)

	testTable := []struct {
		name                string
		headerValue         string
		userId              int
		listId              int
		output              todo.TodoList
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerValue: "3",
			userId:      3,
			listId:      2,
			output: todo.TodoList{
				Id:          2,
				Title:       "test title",
				Description: "test description",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, output todo.TodoList) {
				s.EXPECT().GetById(userId, listId).Return(output, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":2,"title":"test title","description":"test description"}`,
		},
		{
			name:                "No Header",
			mockBehavior:        func(s *mock_service.MockTodoList, userId, listId int, output todo.TodoList) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"user unauthorized"}`,
		},
		{
			name:        "Service Failure",
			headerValue: "3",
			userId:      3,
			listId:      2,
			output:      todo.TodoList{},
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, output todo.TodoList) {
				s.EXPECT().GetById(userId, listId).Return(output, errors.New("service failure"))
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

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, testCase.userId, testCase.listId, testCase.output)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/api/lists/:id", handler.getListById)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%d", testCase.listId), nil)
			req.AddCookie(&http.Cookie{Name: userCtx, Value: testCase.headerValue})

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestList_UpdateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId int, input todo.UpdateListInput)

	testTable := []struct {
		name                string
		headerValue         string
		userId              int
		listId              int
		input               todo.UpdateListInput
		inputString         string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK No Empty Fields",
			headerValue: "7",
			userId:      7,
			listId:      2,
			input: todo.UpdateListInput{
				Title:       stringPointer("new title"),
				Description: stringPointer("new description"),
			},
			inputString: `{"title":"new title","description":"new description"}`,
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, input todo.UpdateListInput) {
				s.EXPECT().Update(userId, listId, input).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:                "No Header",
			mockBehavior:        func(s *mock_service.MockTodoList, userId, listId int, input todo.UpdateListInput) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"user unauthorized"}`,
		},
		{
			name:        "OK One Field Empty",
			headerValue: "7",
			userId:      7,
			listId:      2,
			input: todo.UpdateListInput{
				Title:       stringPointer(""),
				Description: stringPointer("new description"),
			},
			inputString: `{"title":"","description":"new description"}`,
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, input todo.UpdateListInput) {
				s.EXPECT().Update(userId, listId, input).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:        "OK All Empty",
			headerValue: "7",
			userId:      7,
			listId:      2,
			input: todo.UpdateListInput{
				Title:       stringPointer(""),
				Description: stringPointer(""),
			},
			inputString: `{"title":"","description":""}`,
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, input todo.UpdateListInput) {
				s.EXPECT().Update(userId, listId, input).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, testCase.userId, testCase.listId, testCase.input)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.PUT("/api/lists/:id", handler.updateList)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%d", testCase.listId), bytes.NewBufferString(testCase.inputString))
			req.AddCookie(&http.Cookie{Name: userCtx, Value: testCase.headerValue})

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func TestList_DeleteList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId int)

	testTable := []struct {
		name                string
		headerValue         string
		userId              int
		listId              int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerValue: "2",
			userId:      2,
			listId:      4,
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().Delete(userId, listId).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:                "No Header",
			mockBehavior:        func(s *mock_service.MockTodoList, userId, listId int) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"user unauthorized"}`,
		},
		{
			name:        "Service Failure",
			headerValue: "2",
			userId:      2,
			listId:      4,
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().Delete(userId, listId).Return(errors.New("service failure"))
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

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, testCase.userId, testCase.listId)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.DELETE("/api/lists/:id", handler.deleteList)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/lists/%d", testCase.listId), nil)
			req.AddCookie(&http.Cookie{Name: userCtx, Value: testCase.headerValue})

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
