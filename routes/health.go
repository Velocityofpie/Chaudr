package routes

import "net/http"

func addHealthCheck(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.Write([]byte("healthy"))
	})
}
