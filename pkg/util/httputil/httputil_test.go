package httputil

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
)

type testContextKey string

const contextKey testContextKey = "key"

func TestRequestWrap(t *testing.T) {
	var testVal string
	var testStatus int
	var testErr error

	handler := RequestWrap(func(r *http.Request) (*http.Request, int, error) {
		r = r.WithContext(context.WithValue(r.Context(), contextKey, testVal))
		return r, testStatus, testErr
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, ok := r.Context().Value(contextKey).(string)
		if !ok {
			http.Error(w, "bad string", http.StatusInternalServerError)
		}
		if _, err := fmt.Fprint(w, val); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	testCases := []struct {
		name         string
		val          string
		status       int
		err          error
		expectedBody string
	}{
		{name: "normal", val: "123", status: http.StatusOK, expectedBody: "123"},
		{name: "err", status: http.StatusUnprocessableEntity, err: fmt.Errorf("unprocessible"),
			expectedBody: "{\"error\":\"unprocessible\",\"code\":422,\"status_text\":\"Unprocessable Entity\"}\n"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			testVal, testStatus, testErr = tc.val, tc.status, tc.err

			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/", nil))

			require.Equal(tc.status, resp.Code)
			require.Equal(tc.expectedBody, resp.Body.String())
		})
	}
}

type testObj struct {
	Val string `json:"val,omitempty"`
}

func TestRespondJSONWrap(t *testing.T) {
	var testVal string
	var testStatus int
	var testErr error
	f := RespondJSONWrap(func(r *http.Request) (any, int, error) {
		if r.Method != http.MethodGet {
			return nil, http.StatusInternalServerError, fmt.Errorf("not a GET")
		}
		return testObj{Val: testVal}, testStatus, testErr
	})

	testCases := []struct {
		name         string
		method       string
		val          string
		status       int
		err          error
		expectedBody string
	}{
		{name: "normal", val: "123", status: http.StatusOK, expectedBody: "{\"val\":\"123\"}\n"},
		{name: "post", method: http.MethodPost, status: http.StatusInternalServerError,
			expectedBody: "{\"error\":\"not a GET\",\"code\":500,\"status_text\":\"Internal Server Error\"}\n"},
		{name: "err", status: http.StatusUnprocessableEntity, err: fmt.Errorf("unprocessible"),
			expectedBody: "{\"error\":\"unprocessible\",\"code\":422,\"status_text\":\"Unprocessable Entity\"}\n"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			testVal, testStatus, testErr = tc.val, tc.status, tc.err

			resp := httptest.NewRecorder()
			f(resp, httptest.NewRequest(tc.method, "/", nil))

			require.Equal(tc.status, resp.Code)
			require.Equal(tc.expectedBody, resp.Body.String())
		})
	}
}

func TestRespondError(t *testing.T) {
	require := require.New(t)

	resp := httptest.NewRecorder()
	status, err := http.StatusInternalServerError, fmt.Errorf("my error")

	RespondError(resp, status, err)

	require.Equal(status, resp.Code)
	require.Equal("{\"error\":\"my error\",\"code\":500,\"status_text\":\"Internal Server Error\"}\n", resp.Body.String())
	require.Equal(resp.Header().Get("Content-Type"), JSONContentType)
	require.Equal(resp.Header().Get("X-Content-Type-Options"), "nosniff")
}

func TestRespondJSON(t *testing.T) {
	require := require.New(t)

	resp := httptest.NewRecorder()

	RespondJSON(resp, testObj{Val: "my test"})

	require.Equal(http.StatusOK, resp.Code)
	require.Equal("{\"val\":\"my test\"}\n", resp.Body.String())
	require.Equal(resp.Header().Get("Content-Type"), JSONContentType)
}

func TestBindJSON(t *testing.T) {
	testCases := []struct {
		name           string
		data           []byte
		bindTo         any
		expectedValue  string
		expectedStatus int
		expectedError  string
	}{
		{name: "normal", data: test.JSON(t, testObj{Val: "testy"}), bindTo: &testObj{}, expectedValue: "testy"},
		{name: "empty data", data: []byte("{}"), bindTo: &testObj{}, expectedValue: ""},
		{name: "not applicable data", data: []byte("{ \"ok\": \"me\" }"), bindTo: &testObj{}, expectedValue: ""},
		{name: "bad data", data: nil, bindTo: &testObj{}, expectedStatus: http.StatusBadRequest,
			expectedError: "unexpected end of JSON input"},
		{name: "non-json", data: []byte("abc"), bindTo: &testObj{}, expectedStatus: http.StatusBadRequest,
			expectedError: "invalid character 'a' looking for beginning of value"},
		{name: "empty bind", data: []byte("{}"), bindTo: nil, expectedStatus: http.StatusBadRequest,
			expectedError: "json: Unmarshal(nil)"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			req := httptest.NewRequest("", "/", bytes.NewBuffer(tc.data))

			status, err := BindJSON(req, tc.bindTo)
			require.Equal(tc.expectedStatus, status)
			if tc.expectedError == "" {
				require.NoError(err)
				bindTo, ok := tc.bindTo.(*testObj)
				require.True(ok)
				require.Equal(tc.expectedValue, bindTo.Val)
			} else {
				require.Equal(tc.expectedError, err.Error())
			}
		})
	}
}
