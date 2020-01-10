package athens

type Request struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

type Response struct {
	Success bool
	Message string
}
