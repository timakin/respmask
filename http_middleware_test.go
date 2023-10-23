package respmask_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/timakin/respmask"
)

func TestMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/api/", respmask.NewMaskingMiddleware(dynamicKeysToMask, http.StripPrefix("/api", apiHandler())))
	server := httptest.NewServer(mux)
	defer server.Close()

	tests := []struct {
		name         string
		endpoint     string
		expectedBody string
	}{
		{
			name:         "masking email and password for /api/data endpoint",
			endpoint:     "/api/data",
			expectedBody: `{"email":"t***@example.com","password":"**********"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(server.URL + tt.endpoint)
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}
}

func apiHandler() http.Handler {
	apiMux := http.NewServeMux()
	apiMux.Handle("/data", http.HandlerFunc(handleSuccessCase))
	return apiMux
}

func handleSuccessCase(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"email":    "test@example.com",
		"password": "supersecret",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func dynamicKeysToMask(r *http.Request) map[string]respmask.MaskingFunc {
	switch r.URL.Path {
	case "/api/data":
		return map[string]respmask.MaskingFunc{
			"email":    respmask.DefaultMaskingRules[respmask.EmailMasking],
			"password": respmask.DefaultMaskingRules[respmask.PasswordMasking],
		}
	default:
		return map[string]respmask.MaskingFunc{}
	}
}
