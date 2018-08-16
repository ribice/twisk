package mw

import "net/http"

// CORS adds support for Cross-Origin Resource Sharing
func CORS(hf http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "POST")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Add("Access-Control-Max-Age", "86400")
			w.WriteHeader(200)
			return
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")

		hf.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
