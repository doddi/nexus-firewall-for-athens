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

type LocalServer struct {
	Port      int
	Validator validate.Validator
}

func (server LocalServer) Handle() {
	server.handleRequests(server.Port)
}

func (server LocalServer) handleConvert(writer http.ResponseWriter, request *http.Request) {
	if "POST" != request.Method {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	hookMessage, err := server.decodeMessage(request)

	if err != nil {
		fmt.Println(err)
		return
	}

	if server.Validator.Validate(hookMessage) {
		writer.WriteHeader(http.StatusOK)
	}
	writer.WriteHeader(http.StatusForbidden)
}

func (server LocalServer) decodeMessage(request *http.Request) (athens.Request, error) {
	hookMessage := athens.Request{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&hookMessage)
	return hookMessage, err
}

func (server LocalServer) handleRequests(port int) {
	http.HandleFunc("/", server.handleConvert)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
