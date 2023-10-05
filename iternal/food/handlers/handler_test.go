package fh

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"people-food-service/iternal/food"
	fooddto "people-food-service/iternal/food/dto"
	mock_food "people-food-service/iternal/food/mock"
	logging "people-food-service/pkg/client/logger"
	"testing"
)

func TestHandler_GetOne(t *testing.T) {
	type mockBehavior func(s *mock_food.MockRepository, f fooddto.RequestDTO)

	testTable := []struct {
		name                 string
		inputFood            fooddto.RequestDTO
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody fooddto.ResponseDTO
	}{
		{
			name: "OK",
			inputFood: fooddto.RequestDTO{
				UUID:  "",
				Name:  "Пицца",
				Price: 7.85,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {
				s.EXPECT().FindOne(context.TODO(), f.Name).Return(food.Food{
					UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
					Name:  "Пицца",
					Price: 7.85,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fooddto.ResponseDTO{
				Food: []food.Food{
					{
						UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
						Name:  "Пицца",
						Price: 7.85,
					},
				},
				ResponseStatus: "OK",
			},
		},
		{
			name: "Empty name",
			inputFood: fooddto.RequestDTO{
				UUID:  "",
				Name:  "",
				Price: 0,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'Name' should be not empty and consists only alphabet characters",
			},
		},
		{
			name:      "Empty json",
			inputFood: fooddto.RequestDTO{},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'Name' should be not empty and consists only alphabet characters",
			},
		},
		{
			name: "Price less then 0",
			inputFood: fooddto.RequestDTO{
				UUID:  "",
				Name:  "Пицца",
				Price: -5.24,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: price have to be greater than 0",
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_food.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.inputFood)
			l := logrus.New()
			level, _ := logrus.ParseLevel("trace")
			l.SetLevel(level)
			le := logrus.NewEntry(l)
			logger := logging.Logger{le}
			router := chi.NewRouter()

			router.Get("/api/food", GetOne(context.TODO(), &logger, repo))

			w := httptest.NewRecorder()
			marshal, err := json.Marshal(testCase.inputFood)
			if err != nil {
				return
			}
			req := httptest.NewRequest("GET", "/api/food", bytes.NewReader(marshal))

			router.ServeHTTP(w, req)
			var resp fooddto.ResponseDTO
			err = json.Unmarshal([]byte(w.Body.String()), &resp)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedResponseBody, resp)
		})
	}
}
func TestHandler_Create(t *testing.T) {
	type mockBehavior func(s *mock_food.MockRepository, f fooddto.RequestDTO)

	testTable := []struct {
		name                 string
		inputFood            fooddto.RequestDTO
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody fooddto.ResponseDTO
	}{
		{
			name: "OK without UUID",
			inputFood: fooddto.RequestDTO{
				UUID:  "",
				Name:  "Пицца",
				Price: 7.85,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {
				s.EXPECT().Create(context.TODO(), food.Food(f)).Return("d41b9758-f344-447f-b512-cc35b89c23e9", nil)
			},
			expectedStatusCode: 201,
			expectedResponseBody: fooddto.ResponseDTO{
				Food: []food.Food{
					{
						UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
						Name:  "Пицца",
						Price: 7.85,
					},
				},
				ResponseStatus: "OK",
			},
		},
		{
			name: "OK with UUID",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc11b57c23e9",
				Name:  "Пицца",
				Price: 7.85,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {
				s.EXPECT().Create(context.TODO(), food.Food(f)).Return("d41b9758-f344-447f-b512-cc11b57c23e9", nil)
			},
			expectedStatusCode: 201,
			expectedResponseBody: fooddto.ResponseDTO{
				Food: []food.Food{
					{
						UUID:  "d41b9758-f344-447f-b512-cc11b57c23e9",
						Name:  "Пицца",
						Price: 7.85,
					},
				},
				ResponseStatus: "OK",
			},
		},
		{
			name: "Empty name",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc11b57c23e9",
				Name:  "",
				Price: 7.85,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'Name' should be not empty and consists only alphabet characters",
			},
		},
		{
			name:      "Empty json",
			inputFood: fooddto.RequestDTO{},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'Name' should be not empty and consists only alphabet characters",
			},
		},
		{
			name: "Price less then 0",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc11b57c23e9",
				Name:  "Пицца",
				Price: -5.24,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: price have to be greater than 0",
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_food.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.inputFood)
			l := logrus.New()
			level, _ := logrus.ParseLevel("trace")
			l.SetLevel(level)
			le := logrus.NewEntry(l)
			logger := logging.Logger{le}
			router := chi.NewRouter()

			router.Post("/api/food", Create(context.TODO(), &logger, repo))

			w := httptest.NewRecorder()
			marshal, err := json.Marshal(testCase.inputFood)
			if err != nil {
				return
			}
			req := httptest.NewRequest("POST", "/api/food", bytes.NewReader(marshal))

			router.ServeHTTP(w, req)
			var resp fooddto.ResponseDTO
			err = json.Unmarshal([]byte(w.Body.String()), &resp)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedResponseBody, resp)
		})
	}
}
func TestHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_food.MockRepository, f fooddto.RequestDTO)

	testTable := []struct {
		name                 string
		inputFood            fooddto.RequestDTO
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody fooddto.ResponseDTO
	}{
		{
			name: "OK",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
				Name:  "Пицца",
				Price: 7.45,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {
				s.EXPECT().Delete(context.TODO(), food.Food(f)).Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "OK",
			},
		},
		{
			name: "OK with empty body",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
				Name:  "",
				Price: 0,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {
				s.EXPECT().Delete(context.TODO(), food.Food(f)).Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "OK",
			},
		},
		{
			name: "Empty uuid",
			inputFood: fooddto.RequestDTO{
				UUID:  "",
				Name:  "",
				Price: 0,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'UUID' should be not empty and consists uuid",
			},
		},
		{
			name:      "Empty json",
			inputFood: fooddto.RequestDTO{},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'UUID' should be not empty and consists uuid",
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_food.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.inputFood)
			l := logrus.New()
			level, _ := logrus.ParseLevel("trace")
			l.SetLevel(level)
			le := logrus.NewEntry(l)
			logger := logging.Logger{le}
			router := chi.NewRouter()

			router.Delete("/api/food", Delete(context.TODO(), &logger, repo))

			w := httptest.NewRecorder()
			marshal, err := json.Marshal(testCase.inputFood)
			if err != nil {
				return
			}
			req := httptest.NewRequest("DELETE", "/api/food", bytes.NewReader(marshal))

			router.ServeHTTP(w, req)
			var resp fooddto.ResponseDTO
			err = json.Unmarshal([]byte(w.Body.String()), &resp)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedResponseBody, resp)
		})
	}
}
func TestHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_food.MockRepository, f fooddto.RequestDTO)

	testTable := []struct {
		name                 string
		inputFood            fooddto.RequestDTO
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody fooddto.ResponseDTO
	}{
		{
			name: "OK",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
				Name:  "Пицца",
				Price: 7.45,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {
				s.EXPECT().Update(context.TODO(), food.Food(f)).Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "OK",
			},
		},
		{
			name: "Empty uuid",
			inputFood: fooddto.RequestDTO{
				UUID:  "",
				Name:  "Пицца",
				Price: 7.45,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field ID is required",
			},
		},
		{
			name: "Empty name",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
				Name:  "",
				Price: 7.45,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'Name' should be not empty and consists only alphabet characters",
			},
		},
		{
			name: "Zero price",
			inputFood: fooddto.RequestDTO{
				UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
				Name:  "Пицца",
				Price: 0,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: price have to be greater than 0",
			},
		},
		{
			name:      "Empty json",
			inputFood: fooddto.RequestDTO{},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: fooddto.ResponseDTO{
				Food:           []food.Food(nil),
				ResponseStatus: "ERROR: field 'Name' should be not empty and consists only alphabet characters",
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_food.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.inputFood)
			l := logrus.New()
			level, _ := logrus.ParseLevel("trace")
			l.SetLevel(level)
			le := logrus.NewEntry(l)
			logger := logging.Logger{le}
			router := chi.NewRouter()

			router.Patch("/api/food", Update(context.TODO(), &logger, repo))

			w := httptest.NewRecorder()
			marshal, err := json.Marshal(testCase.inputFood)
			if err != nil {
				return
			}
			req := httptest.NewRequest("PATCH", "/api/food", bytes.NewReader(marshal))

			router.ServeHTTP(w, req)
			var resp fooddto.ResponseDTO
			err = json.Unmarshal([]byte(w.Body.String()), &resp)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedResponseBody, resp)
		})
	}
}
