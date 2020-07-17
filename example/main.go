// +build !codeanalysis

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alkemics/goflow/example/graphs"
)

func httpPlayground(pg graphs.Playground) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var in json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err,
			})
			return
		}

		name := r.URL.Query().Get("name")
		out, err := pg.Run(r.Context(), name, &in)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err,
			})
			return
		}

		_ = json.NewEncoder(w).Encode(out)
	}
}

func main() {
	playground := graphs.NewPlayground(false)

	addr := "127.0.0.1:8080"
	fmt.Printf("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, httpPlayground(playground)))
}
