package jhttp

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

	handler := RequestWrap(func(r *http.Request) (*http.Request, *HTTPError) {
		r = r.WithContext(context.WithValue(r.Context(), contextKey, testVal))
		if testStatus != http.StatusOK {
			return nil, Error(testStatus, testErr)
		}
		return r, nil
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
		tc := tc
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
	handlerFunc := ResponseWrap(func(r *http.Request) (any, *HTTPError) {
		if r.Method != http.MethodGet {
			return nil, Error(http.StatusInternalServerError, fmt.Errorf("not a GET"))
		}
		if testStatus != http.StatusOK {
			return nil, Error(testStatus, testErr)
		}
		return testObj{Val: testVal}, nil
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			testVal, testStatus, testErr = tc.val, tc.status, tc.err

			resp := httptest.NewRecorder()
			handlerFunc(resp, httptest.NewRequest(tc.method, "/", nil))

			require.Equal(tc.status, resp.Code)
			require.Equal(tc.expectedBody, resp.Body.String())
		})
	}
}

func TestRespondError(t *testing.T) {
	require := require.New(t)

	resp := httptest.NewRecorder()
	httpError := Error(http.StatusInternalServerError, fmt.Errorf("my error"))

	RespondError(resp, httpError)

	require.Equal(httpError.Code, resp.Code)
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

func TestExtractJSON(t *testing.T) {
	testCases := []struct {
		name          string
		data          []byte
		extractTo     any
		expectedValue string
		expectedCode  int
		expectedError string
	}{
		{name: "normal", data: test.JSON(t, testObj{Val: "testy"}), extractTo: &testObj{}, expectedValue: "testy"},
		{name: "empty data", data: []byte("{}"), extractTo: &testObj{}, expectedValue: ""},
		{name: "not applicable data", data: []byte("{ \"ok\": \"me\" }"), extractTo: &testObj{}, expectedValue: ""},
		{name: "bad data", data: nil, extractTo: &testObj{}, expectedCode: http.StatusBadRequest,
			expectedError: "unexpected end of JSON input"},
		{name: "non-json", data: []byte("abc"), extractTo: &testObj{}, expectedCode: http.StatusBadRequest,
			expectedError: "invalid character 'a' looking for beginning of value"},
		{name: "empty bind", data: []byte("{}"), extractTo: nil, expectedCode: http.StatusBadRequest,
			expectedError: "json: Unmarshal(nil)"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			req := httptest.NewRequest("", "/", bytes.NewBuffer(tc.data))

			httpError := ExtractJSON(req, tc.extractTo)
			if tc.expectedError == "" {
				require.Nil(httpError)
				bindTo, ok := tc.extractTo.(*testObj)
				require.True(ok)
				require.Equal(tc.expectedValue, bindTo.Val)
			} else {
				require.Equal(tc.expectedCode, httpError.Code)
				require.Equal(tc.expectedError, httpError.Cause.Error())
			}
		})
	}
}

func TestReturnModelOr500(t *testing.T) {
	httpErr := &HTTPError{Code: http.StatusInternalServerError, Cause: fmt.Errorf("waka")}
	testCases := []struct {
		name  string
		model any
		err   error

		expectedModel any
		expectedErr   *HTTPError
	}{
		{name: "success", model: "model", expectedModel: "model"},
		{name: "err", err: httpErr.Cause, expectedErr: httpErr},
		{name: "http_error", err: httpErr, expectedErr: httpErr},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			model, err := ReturnModelOr500(func() (any, error) {
				return tc.model, tc.err
			})
			require.Equal(tc.expectedModel, model)
			require.Equal(tc.expectedErr, err)
		})
	}
}
