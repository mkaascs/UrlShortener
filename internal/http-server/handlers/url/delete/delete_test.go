package delete

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/logging"
	"url-shortener/internal/storage"
)

func Test_DeleteHandler(t *testing.T) {
	tests := []struct {
		name           string
		alias          string
		mockErr        error
		expectedStatus int
		callMock       bool
	}{
		{
			name:           "success delete",
			alias:          "mkaascs",
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			callMock:       true,
		},
		{
			name:           "alias does not exist",
			alias:          "mkaascs",
			mockErr:        storage.ErrURLNotFound,
			expectedStatus: http.StatusNotFound,
			callMock:       true,
		},
		{
			name:           "internal server error",
			alias:          "mkaascs",
			mockErr:        errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
			callMock:       true,
		},
		{
			name:           "empty alias",
			alias:          "",
			mockErr:        nil,
			expectedStatus: http.StatusBadRequest,
			callMock:       false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUrlDeleter := mocks.NewMockURLDeleter(ctrl)

			if test.callMock {
				mockUrlDeleter.EXPECT().
					DeleteURL(test.alias).
					Return(test.mockErr).
					Times(1)
			} else {
				mockUrlDeleter.EXPECT().
					DeleteURL(gomock.Any()).
					Times(0)
			}

			handler := New(logging.NewPlugLogger(), mockUrlDeleter)

			target := fmt.Sprintf("/url/%s", test.alias)
			req, err := http.NewRequest(http.MethodDelete, target, nil)
			require.NoError(t, err)

			if test.alias != "" {
				routeCtx := chi.NewRouteContext()
				routeCtx.URLParams.Add("alias", test.alias)
				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, test.expectedStatus, rr.Code)
		})
	}
}
