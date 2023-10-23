package respmask

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MaskingMiddleware struct {
	keysAndModeFunc func(r *http.Request) (map[string]MaskingFunc, MaskingMode)
	next            http.Handler
}

func NewMaskingMiddleware(keysAndModeFunc func(r *http.Request) (map[string]MaskingFunc, MaskingMode), next http.Handler) *MaskingMiddleware {
	return &MaskingMiddleware{
		keysAndModeFunc: keysAndModeFunc,
		next:            next,
	}
}

func (m *MaskingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keysFuncs, mode := m.keysAndModeFunc(r)

	recorder := &responseRecorder{
		ResponseWriter: w,
		body:           new(bytes.Buffer),
	}

	m.next.ServeHTTP(recorder, r)

	if recorder.statusCode < 200 || recorder.statusCode >= 300 {
		w.Write(recorder.body.Bytes())
		return
	}

	bodyBytes := recorder.body.Bytes()
	var data map[string]interface{}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		w.Write(bodyBytes)
		return
	}

	Mask(data, keysFuncs, mode)
	maskedBytes, _ := json.Marshal(data)
	w.Write(maskedBytes)
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
