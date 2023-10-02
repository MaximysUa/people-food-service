package fh

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
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
				ResponseStatus: "ok",
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
			require.Equal(t, testCase.expectedResponseBody, resp)
		})
	}
}
