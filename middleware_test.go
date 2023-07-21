package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Authorized"))
	})
	t.Run("returns unauthorized when user unknown", func(t *testing.T) {
		t.Parallel()
		request, err := http.NewRequest(http.MethodGet, "/test", nil)
		request.SetBasicAuth("unknown", "")
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		basicAuthMiddleware(handler).ServeHTTP(response, request)

		wantCode := http.StatusUnauthorized
		gotCode := response.Code
		if wantCode != gotCode {
			t.Errorf("want status code %d, got %d", wantCode, gotCode)
		}
		wantHeader := `Basic realm="Restricted"`
		gotHeader := response.Header().Get("WWW-Authenticate")
		if wantHeader != gotHeader {
			t.Errorf("want header %s, got %s", wantHeader, gotHeader)
		}
	})

	for _, user := range []string{"testA", "testB"} {
		test := fmt.Sprintf("authorizes '%s' username", user)
		t.Run(test, func(t *testing.T) {
			t.Parallel()
			request, err := http.NewRequest(http.MethodGet, "/test", nil)
			request.SetBasicAuth(user, "")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()

			basicAuthMiddleware(handler).ServeHTTP(response, request)

			want := http.StatusOK
			got := response.Code
			if want != got {
				t.Errorf("want status code %d, got %d", want, got)
			}
		})
	}
}
