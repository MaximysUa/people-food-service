package ph

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
	"people-food-service/iternal/person"

	mock_person "people-food-service/iternal/person/mock"

	persondto "people-food-service/iternal/person/dto"
	logging "people-food-service/pkg/client/logger"
	"testing"
)

func TestHandler_GetOne(t *testing.T) {
	type mockBehavior func(s *mock_person.MockRepository, p persondto.RequestDTO)

	testTable := []struct {
		name                 string
		inputPerson          persondto.RequestDTO
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody persondto.ResponseDTO
	}{
		{
			name: "OK",
			inputPerson: persondto.RequestDTO{
				UUID:       "",
				Name:       "Диман",
				FamilyName: "Рекрент",
				Food:       []food.Food{},
			},
			mockBehavior: func(s *mock_person.MockRepository, p persondto.RequestDTO) {
				s.EXPECT().FindOne(context.TODO(), p.Name, p.FamilyName).Return(person.Person{
					UUID:       "48775c0e-820b-4f95-a984-85aa68a88475",
					Name:       "Диман",
					FamilyName: "Рекрент",
					Food: []food.Food{
						{
							UUID:  "41b72d27-c250-4a3a-8c0b-8a7de570a564",
							Name:  "Бурито",
							Price: 9.55,
						},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: persondto.ResponseDTO{
				Person: []person.Person{
					{
						UUID:       "48775c0e-820b-4f95-a984-85aa68a88475",
						Name:       "Диман",
						FamilyName: "Рекрент",
						Food: []food.Food{
							{
								UUID:  "41b72d27-c250-4a3a-8c0b-8a7de570a564",
								Name:  "Бурито",
								Price: 9.55,
							},
						},
					},
				},
				ResponseStatus: "OK",
			},
		},
		{
			name: "Empty name",
			inputPerson: persondto.RequestDTO{
				UUID:       "",
				Name:       "",
				FamilyName: "Рекрент",
				Food:       []food.Food{},
			},
			mockBehavior: func(s *mock_person.MockRepository, p persondto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: persondto.ResponseDTO{

				ResponseStatus: "ERROR: field 'Name' and 'Family name' should be not empty and consists only alphabet characters",
			},
		},
		{
			name: "Empty family name",
			inputPerson: persondto.RequestDTO{
				UUID:       "",
				Name:       "Рекрент",
				FamilyName: "",
				Food:       []food.Food{},
			},
			mockBehavior: func(s *mock_person.MockRepository, p persondto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: persondto.ResponseDTO{

				ResponseStatus: "ERROR: field 'Name' and 'Family name' should be not empty and consists only alphabet characters",
			},
		},
		{
			name:        "Empty json",
			inputPerson: persondto.RequestDTO{},
			mockBehavior: func(s *mock_person.MockRepository, p persondto.RequestDTO) {

			},
			expectedStatusCode: 400,
			expectedResponseBody: persondto.ResponseDTO{

				ResponseStatus: "ERROR: field 'Name' and 'Family name' should be not empty and consists only alphabet characters",
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_person.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.inputPerson)
			l := logrus.New()
			level, _ := logrus.ParseLevel("trace")
			l.SetLevel(level)
			le := logrus.NewEntry(l)
			logger := logging.Logger{le}
			router := chi.NewRouter()

			router.Get("/api/person", GetOne(context.TODO(), &logger, repo))

			w := httptest.NewRecorder()
			marshal, err := json.Marshal(testCase.inputPerson)
			if err != nil {
				return
			}
			req := httptest.NewRequest("GET", "/api/person", bytes.NewReader(marshal))

			router.ServeHTTP(w, req)
			var resp persondto.ResponseDTO
			err = json.Unmarshal([]byte(w.Body.String()), &resp)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedResponseBody, resp)
		})
	}
}
