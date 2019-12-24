package ossindex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexus-firewall-for-athens/athens"
	"strings"
)

type Coordinate struct {
	Type       string `json:"type"`
	Namespace  string `json:"namespace,omitempty"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Qualifiers string `json:"qualifiers,omitempty"`
	Subpath    string `json:"subpath, omitempty"`
}

type Vulnerability struct{}

type Report struct {
	Coordinates     string          `json:"coordinates"`
	Reference       string          `json:"reference"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type OssIndex struct {
}

func (o OssIndex) Validate(request athens.Request) bool {
	coord := Coordinate{Type: "golang", Name: request.Module, Version: request.Version}

	purl := o.convertToPurl(coord)

	return o.checkComponent(purl)
}

func (o OssIndex) decodeMessage(response *http.Response) (Report, error) {
	report := Report{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&report)
	return report, err
}

func (o OssIndex) checkComponent(purl string) bool {
	const baseUrl = "https://ossindex.sonatype.org"
	const endpoint = "/api/v3/component-report/"

	resp, err := http.Get(baseUrl + endpoint + purl)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer resp.Body.Close()

	report, err := o.decodeMessage(resp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// TODO Validator types of vulnerabilities
	if len(report.Vulnerabilities) > 0 {
		return false
	}
	return true
}

func (o OssIndex) convertToPurl(coord Coordinate) string {
	builder := strings.Builder{}
	builder.WriteString("pkg:" + coord.Type + "/")
	if coord.Namespace != "" {
		builder.WriteString(coord.Namespace + "/")
	}
	builder.WriteString(coord.Name)
	if coord.Version != "" {
		builder.WriteString("@" + coord.Version)
	}
	if coord.Qualifiers != "" {
		builder.WriteString("?" + coord.Qualifiers)
	}
	if coord.Subpath != "" {
		builder.WriteString("#" + coord.Subpath)
	}
	return builder.String()
}
