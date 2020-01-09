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
	FailBuild bool
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

	response := server.Validator.Validate(hookMessage)
	if response.Success || !server.FailBuild {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusForbidden)
	}
	message, _ := json.Marshal(&response)
	writer.Write(message)
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
