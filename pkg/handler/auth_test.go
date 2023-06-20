package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
	todo "todo-app"
	"todo-app/pkg/service"
	mock_service "todo-app/pkg/service/mocks"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           todo.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"username": "test", "password": "qwerty"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("service failure"))
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

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/auth/sign-up", handler.signUp)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedRequestBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, input signInInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           signInInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"qwerty"}`,
			inputUser: signInInput{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"token"}`,
		},
		{
			name:      "Incorrect password",
			inputBody: `{"username":"test", "password":"qwertyu"}`,
			inputUser: signInInput{
				Username: "test",
				Password: "qwertyu",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).
					Return("", errors.New("incorrect login or password"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"incorrect login or password"}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"password":"qwerty"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, input signInInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		//{
		//	name:      "Service Failure",
		//	inputBody: `{"username":"test", "password":"qwerty"}`,
		//	inputUser: signInInput{
		//		Username: "test",
		//		Password: "qwerty",
		//	},
		//	mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
		//		s.EXPECT().GenerateToken(input.Username, input.Password).
		//			Return("", errors.New("service failure"))
		//	},
		//	expectedStatusCode:  500,
		//	expectedRequestBody: `{"message":"service failure"}`,
		//},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedRequestBody)
		})
	}
}
