package fh

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
				Price: 0,
			},
			mockBehavior: func(s *mock_food.MockRepository, f fooddto.RequestDTO) {
				s.EXPECT().FindOne(context.TODO(), f.Name).Return(food.Food{
					UUID:  "d41b9758-f344-447f-b512-cc35b89c23e9",
					Name:  "Пицца",
					Price: 7.85,
				})
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
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_food.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.inputFood)

			router := chi.NewRouter()
			router.Get("/api/food", GetOne(context.TODO(), logging.GetLogger(), repo))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/food", bytes.NewBufferString(testCase.inputFood.Name))

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
