package test

import (
	srvHandler "L0task/pkg/handler"
	"L0task/pkg/model"
	mockService "L0task/pkg/test/mocks"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetAllOrders(t *testing.T) {
	type mockBehavior func(s *mockService.MockIService)
	ordersOutput := make([]model.Order, 1, 1)

	jsonErr := json.Unmarshal([]byte(orderJson), &ordersOutput[0])
	if jsonErr != nil {
		return
	}

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockService.MockIService) {
				s.EXPECT().GetAllOrders().Return(ordersOutput, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: allOrdersJson,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mockService.NewMockIService(c)
			testCase.mockBehavior(service)

			handler := srvHandler.NewHandler(service)

			testSrv := gin.New()
			testSrv.GET("/", handler.GetAllOrders)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)

			testSrv.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, recorder.Body.String())
		})
	}
}

func TestHandler_GetAllCachedOrders(t *testing.T) {
	type mockBehavior func(s *mockService.MockIService)
	ordersOutput := make([]model.Order, 1, 1)

	jsonErr := json.Unmarshal([]byte(orderJson), &ordersOutput[0])
	if jsonErr != nil {
		return
	}

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockService.MockIService) {
				s.EXPECT().GetAllCachedOrders().Return(ordersOutput, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: allOrdersJson,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mockService.NewMockIService(c)
			testCase.mockBehavior(service)

			handler := srvHandler.NewHandler(service)

			testSrv := gin.New()
			testSrv.GET("/cache/", handler.GetAllCachedOrders)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/cache/", nil)

			testSrv.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, recorder.Body.String())
		})
	}
}

func TestHandler_GetOrderByUid(t *testing.T) {
	type mockBehavior func(s *mockService.MockIService, uid string)
	var orderOutput model.Order

	jsonErr := json.Unmarshal([]byte(orderJson), &orderOutput)
	if jsonErr != nil {
		return
	}

	tests := []struct {
		name                 string
		uid                  string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			uid:  "b563feb7b2b84b6test",
			mockBehavior: func(s *mockService.MockIService, uid string) {
				s.EXPECT().GetOrderByUid(uid).Return(orderOutput, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: orderJson,
		},
		{
			name: "WrongOrderUid",
			uid:  "test",
			mockBehavior: func(s *mockService.MockIService, uid string) {
				s.EXPECT().GetOrderByUid(uid).Return(model.Order{}, fmt.Errorf("sql: no rows in result set"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "sql: no rows in result set",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mockService.NewMockIService(c)
			testCase.mockBehavior(service, testCase.uid)

			handler := srvHandler.NewHandler(service)

			testSrv := gin.New()
			testSrv.GET("/:uid", handler.GetOrderByUid)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", fmt.Sprintf("/%s", testCase.uid), nil)

			testSrv.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, recorder.Body.String())
		})
	}
}

func TestHandler_GetCachedOrderByUid(t *testing.T) {
	type mockBehavior func(s *mockService.MockIService, uid string)
	var orderOutput model.Order

	jsonErr := json.Unmarshal([]byte(orderJson), &orderOutput)
	if jsonErr != nil {
		return
	}

	tests := []struct {
		name                 string
		uid                  string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			uid:  "b563feb7b2b84b6test",
			mockBehavior: func(s *mockService.MockIService, uid string) {
				s.EXPECT().GetCachedOrderByUid(uid).Return(orderOutput, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: orderJson,
		},
		{
			name: "WrongOrderUid",
			uid:  "test",
			mockBehavior: func(s *mockService.MockIService, uid string) {
				s.EXPECT().GetCachedOrderByUid(uid).Return(model.Order{},
					fmt.Errorf("failed to get value from cache by uid test"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "failed to get value from cache by uid test",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mockService.NewMockIService(c)
			testCase.mockBehavior(service, testCase.uid)

			handler := srvHandler.NewHandler(service)

			testSrv := gin.New()
			testSrv.GET("/cache/:uid", handler.GetCachedOrderByUid)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", fmt.Sprintf("/cache/%s", testCase.uid), nil)

			testSrv.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, recorder.Body.String())
		})
	}
}
