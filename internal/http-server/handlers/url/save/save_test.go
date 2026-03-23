package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/logging"
	"url-shortener/internal/storage"
)

const aliasLen = 6

func Test_SaveHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        any
		mockAlias      string
		mockID         int64
		mockErr        error
		expectedStatus int
		callMock       bool
		checkAlias     bool
	}{
		{
			name:           "saved successfully with cur alias",
			request:        Request{URL: "https://github.com/mkaascs", Alias: "mkaascs"},
			mockAlias:      "mkaascs",
			mockID:         int64(5),
			mockErr:        nil,
			expectedStatus: http.StatusCreated,
			callMock:       true,
			checkAlias:     true,
		},
		{
			name:           "saved successfully with gen alias",
			request:        Request{URL: "https://github.com/mkaascs"},
			mockAlias:      "",
			mockID:         int64(5),
			mockErr:        nil,
			expectedStatus: http.StatusCreated,
			callMock:       true,
			checkAlias:     false,
		},
		{
			name:           "already exists",
			request:        Request{URL: "https://docs.docker.com/get-started/introduction/develop-with-containers/", Alias: "docker-docs"},
			mockAlias:      "docker-docs",
			mockID:         int64(0),
			mockErr:        storage.ErrURLExists,
			expectedStatus: http.StatusConflict,
			callMock:       true,
			checkAlias:     true,
		},
		{
			name:           "internal server error",
			request:        Request{URL: "https://github.com/mkaascs", Alias: "mkaascs"},
			mockAlias:      "mkaascs",
			mockID:         int64(0),
			mockErr:        errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
			callMock:       true,
			checkAlias:     false,
		},
		{
			name:           "invalid request body",
			request:        `{url:some`,
			expectedStatus: http.StatusBadRequest,
			callMock:       false,
		},
		{
			name:           "invalid url",
			request:        Request{URL: "not-a-url", Alias: "mkaascs"},
			expectedStatus: http.StatusUnprocessableEntity,
			callMock:       false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUrlSaver := mocks.NewMockURLSaver(ctrl)

			if test.callMock {
				aliasMatcher := gomock.Any()
				if test.mockAlias != "" {
					aliasMatcher = gomock.Eq(test.mockAlias)
				}

				mockUrlSaver.EXPECT().
					SaveURL(gomock.Any(), aliasMatcher).
					Return(test.mockID, test.mockErr).
					Times(1)

			} else {
				mockUrlSaver.EXPECT().
					SaveURL(gomock.Any(), gomock.Any()).
					Times(0)
			}

			handler := New(logging.NewPlugLogger(), mockUrlSaver, aliasLen)

			var bodyBytes []byte
			switch v := test.request.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, test.expectedStatus, rr.Code)
			if test.expectedStatus == http.StatusCreated {
				var response Response
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)

				if test.checkAlias {
					require.Equal(t, test.mockAlias, response.Alias)
				} else {
					require.Len(t, response.Alias, aliasLen)
				}
			}
		})
	}
}
