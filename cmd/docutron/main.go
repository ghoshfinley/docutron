package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/minimalistsoftware/docutron"
)

func runNew(w http.ResponseWriter, r *http.Request) {

	var request docutron.UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Printf("Error decoding: %s", err)
		// return HTTP 400 bad request
	}
	fmt.Printf("Initialising project: %s\n", request.Project)

	// If the project doesn't exist - create it
	docutron.InitProject(request.Project)

	// Create the JSON Skeleton
	num := docutron.NextNumber(request)
	inv := docutron.NewJSONFile(request, fmt.Sprintf("%s%d", request.Config.Invoice.Prefix, num))

	// Write HTML Invoice from skeleton
	docutron.WriteHTML(request, inv, request.Config.Invoice.Template)

	// Write PDF Invoice from skeleton
	docutron.WritePDFChrome(request, inv)
}

func main() {

	http.HandleFunc("/new", runNew)

	fmt.Println(("Running server on :8081"))
	log.Fatal(http.ListenAndServe(":8081", nil))

}
