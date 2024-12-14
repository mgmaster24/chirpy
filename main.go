package main

import "net/http"

func main() {
	httpServerMux := http.NewServeMux()
	httpServerMux.Handle("/", http.FileServer(http.Dir(".")))
	httpServer := http.Server{
		Addr:    ":8080",
		Handler: httpServerMux,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
