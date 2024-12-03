package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middleWareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	serverHits := cfg.fileserverHits.Load()
	htmlTemplate := `<html>
										<body>
											<h1>Welcome, Chirpy Admin</h1>
											<p>Chirpy has been visited %d times!</p>
										</body>
									</html>`
	w.Write([]byte(fmt.Sprintf(htmlTemplate, serverHits)))
}
