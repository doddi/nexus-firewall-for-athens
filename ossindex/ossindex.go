package ossindex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexus-firewall-for-athens/athens"
	purl "nexus-firewall-for-athens/validate"
)

type Vulnerability struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CvssScore   float64 `json:"cvssScore"`
	Cve         string  `json:"cve"`
	Reference   string  `json:"reference"`
}

type Report struct {
	Coordinates     string          `json:"coordinates"`
	Reference       string          `json:"reference"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type OssIndex struct {
}

func (o OssIndex) Validate(request athens.Request) athens.Response {
	coord := purl.Coordinate{Type: "golang", Name: request.Module, Version: request.Version}
	purl := purl.ConvertToPurlString(coord)

	report := o.checkComponent(purl)

	response := athens.Response{Success: true}
	if len(report.Vulnerabilities) > 0 {
		response.Success = false
		response.Message = fmt.Sprintf("Found %d vulnerabilities, go to %s for more information", len(report.Vulnerabilities), report.Reference)
	}
	return response
}

func (o OssIndex) decodeMessage(response *http.Response) (Report, error) {
	report := Report{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&report)
	return report, err
}

func (o OssIndex) checkComponent(purl string) Report {
	const baseUrl = "https://ossindex.sonatype.org"
	const endpoint = "/api/v3/component-report/"

	resp, err := http.Get(baseUrl + endpoint + purl)
	if err != nil {
		fmt.Println(err)
		return Report{}
	}

	defer resp.Body.Close()

	report, err := o.decodeMessage(resp)
	if err != nil {
		fmt.Println(err)
		return Report{}
	}

	return report
}
