package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nexus-firewall-for-athens/athens"
	purl "nexus-firewall-for-athens/validate"
)

type Component struct {
	PackageUrl string `json:"packageUrl"`
}

type ComponentEvaluationRequest struct {
	Components []Component `json:"components"`
}

type ComponentEvaluationResponse struct {
	ResultId      string `json:"resultId"`
	SubmittedDate string `json:"submittedDate"`
	ApplicationId string `json:"applicationId"`
	ResultsUrl    string `json:"resultsUrl"`
}

type Reason struct {
	Reason string `json:"reason"`
}

type ConstraintViolation struct {
	ConstraintId   string   `json:"constraintId"`
	ConstraintName string   `json:"constraintName"`
	Reasons        []Reason `json:"reasons"`
}

type PolicyViolation struct {
	PolicyId   string `json:"policyId"`
	PolicyName string `json:"policyName"`

	ConstraintViolations []ConstraintViolation `json:"constraintViolations"`
}

type PolicyData struct {
	PolicyViolations []PolicyViolation `json:"policyViolations"`
}

type SecurityIssue struct {
	Source    string `json:"source"`
	Reference string `json:"reference"`
	Url       string `json:"url"`
}

type SecurityData struct {
	SecurityIssues []SecurityIssue `json:"securityIssues"`
}

type License struct {
	LicenseId   string `json:"licenseDd"`
	LicenseName string `json:"licenseName"`
}

type LicenseData struct {
	DeclaredLicenses []License `json:"declaredLicenses"`
	ObservedLicenses []License `json:"observedLicenses"`
}

type Result struct {
	LicenseData  LicenseData  `json:"licenseData"`
	SecurityData SecurityData `json:"securityData"`
	PolicyData   PolicyData   `json:"policyData"`
}

type ComponentEvaluationResult struct {
	SubmittedDate  string   `json:"submittedDate"`
	EvaluationDate string   `json:"evaluationDate"`
	ApplicationId  string   `json:"applicationId"`
	Results        []Result `json:"results"`
	IsError        bool     `json:"isError"`
	ErrorMessage   string   `json:"errorMessage"`
}

type Authentication struct {
	Username string
	Password string
}
type NexusIq struct {
	BaseUrl        string
	ApplicationId  string
	Authentication Authentication
}

func (n NexusIq) Validate(request athens.Request) athens.Response {
	coord := purl.Coordinate{Type: "golang", Name: request.Module, Version: request.Version}
	purlString := purl.ConvertToPurlString(coord)

	component := Component{PackageUrl: purlString}
	submitRequest := ComponentEvaluationRequest{[]Component{component}}
	submitResponse := n.submitComponent(submitRequest)

	response := athens.Response{Success: true}
	if submitResponse.ResultId == "" {
		response.Success = false
		return response
	}

	report, err := n.checkComponent(submitResponse)
	if err != nil {
		response.Success = false
		return response
	}

	if len(report.Results) > 0 && len(report.Results[0].PolicyData.PolicyViolations) > 0 {
		response.Success = false
		response.Message = n.buildViolationResponse(report)
	}
	return response
}

func (n NexusIq) buildViolationResponse(report ComponentEvaluationResult) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Component violates %d policies: \n", len(report.Results[0].PolicyData.PolicyViolations)))
	for _, violation := range report.Results[0].PolicyData.PolicyViolations {
		buffer.WriteString(fmt.Sprintf(" Policy %s, constrain name %s\n", violation.PolicyName, violation.ConstraintViolations[0].ConstraintName))
	}

	buffer.WriteString("\nSecurity Issues: ")
	for _, securityIssue := range report.Results[0].SecurityData.SecurityIssues {
		buffer.WriteString(fmt.Sprintf("%s (%s), ", securityIssue.Reference, securityIssue.Url))
	}
	buffer.WriteString("\n")

	buffer.WriteString("\nLicense Data: ")
	for _, declaredLicenses := range report.Results[0].LicenseData.DeclaredLicenses {
		buffer.WriteString(fmt.Sprintf("%s, ", declaredLicenses.LicenseName))
	}
	buffer.WriteString("\n")

	return buffer.String()
}

func (n NexusIq) decodeSubmitResponse(response *http.Response) (ComponentEvaluationResponse, error) {
	report := ComponentEvaluationResponse{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&report)
	return report, err
}

func (n NexusIq) submitComponent(purl ComponentEvaluationRequest) ComponentEvaluationResponse {
	const endpoint = "/api/v2/evaluation/applications/"

	data, _ := json.Marshal(purl)

	client := &http.Client{}
	req, err := http.NewRequest("POST", n.BaseUrl+endpoint+n.ApplicationId, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(n.Authentication.Username, n.Authentication.Password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ComponentEvaluationResponse{}
	}

	defer resp.Body.Close()

	report, err := n.decodeSubmitResponse(resp)
	if err != nil {
		fmt.Println(err)
		return ComponentEvaluationResponse{}
	}

	return report
}

func (n NexusIq) checkComponent(evaluationResponse ComponentEvaluationResponse) (ComponentEvaluationResult, error) {
	const endpoint = "/api/v2/evaluation/applications/"

	client := &http.Client{}
	req, err := http.NewRequest("GET", n.BaseUrl+endpoint+n.ApplicationId+"/results/"+evaluationResponse.ResultId, nil)
	req.SetBasicAuth(n.Authentication.Username, n.Authentication.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return n.decodeEvaluationResult(resp)
}

func (n NexusIq) decodeEvaluationResult(response *http.Response) (ComponentEvaluationResult, error) {
	report := ComponentEvaluationResult{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&report)
	return report, err
}
