package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
	"todo-app/pkg/service"
	mock_service "todo-app/pkg/service/mocks"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `1`,
		},
		{
			name:                "No Header",
			headerName:          "",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"empty auth header"}`,
		},
		{
			name:                "Invalid Bearer",
			headerName:          "Authorization",
			headerValue:         "Bear token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                "Invalid Token",
			headerName:          "Authorization",
			headerValue:         "Bearer ",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, errors.New("failed to parse token"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"failed to parse token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/protected", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedRequestBody)
		})
	}
}
