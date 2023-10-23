package respmask

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MaskingMiddleware struct {
	keysAndMaskFuncs func(r *http.Request) map[string]MaskingFunc
	next             http.Handler
}

func NewMaskingMiddleware(keysAndMaskFuncs func(r *http.Request) map[string]MaskingFunc, next http.Handler) *MaskingMiddleware {
	return &MaskingMiddleware{
		keysAndMaskFuncs: keysAndMaskFuncs,
		next:             next,
	}
}

func (m *MaskingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keysFuncs := m.keysAndMaskFuncs(r)

	recorder := &responseRecorder{
		ResponseWriter: w,
		body:           new(bytes.Buffer),
	}

	m.next.ServeHTTP(recorder, r)

	if recorder.statusCode >= 200 && recorder.statusCode < 300 {
		bodyBytes := recorder.body.Bytes()
		var data map[string]interface{}

		if err := json.Unmarshal(bodyBytes, &data); err == nil {
			MaskData(data, keysFuncs)
			maskedBytes, _ := json.Marshal(data)
			w.Write(maskedBytes)
		} else {
			w.Write(bodyBytes)
		}
	} else {
		w.Write(recorder.body.Bytes())
	}
}

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
