package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nexus-firewall-for-athens/athens"
	"nexus-firewall-for-athens/validate"

	"strconv"
)

func handleConvert(writer http.ResponseWriter, request *http.Request) {
	if "POST" != request.Method {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	hookMessage, err := decodeMessage(request)

	if err != nil {
		fmt.Println(err)
		return
	}

	if validate.Validate(hookMessage) {
		writer.WriteHeader(http.StatusOK)
	}
	writer.WriteHeader(http.StatusForbidden)
}

func decodeMessage(request *http.Request) (athens.Request, error) {
	hookMessage := athens.Request{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&hookMessage)
	return hookMessage, err
}

func HandleRequests(port int) {
	http.HandleFunc("/", handleConvert)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
