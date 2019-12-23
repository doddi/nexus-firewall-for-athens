package ossindex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Coordinate struct {
	Type   string `json:"type"`
	Namespace string `json:"namespace,omitempty"`
	Name string `json:"name"`
	Version string `json:"version"`
	Qualifiers string `json:"qualifiers,omitempty"`
	Subpath string `json:"subpath, omitempty"`
}

type Vulnerability struct{

}

type Report struct {
	Coordinates string `json:"coordinates"`
	Reference string `json:"reference"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`

}
func Validate(module string, version string) bool {
	coord := Coordinate{Type:"golang", Name:module, Version:version}

	purl := convertToPurl(coord)

	return checkComponent(purl)
}

func decodeMessage(response *http.Response) (Report, error) {
	report := Report{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&report)
	return report, err
}

func checkComponent(purl string) bool {
	resp, err := http.Get("https://ossindex.sonatype.org" + "/api/v3/component-report/" + purl)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer resp.Body.Close()

	report, err := decodeMessage(resp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// TODO Validate types of vulnerabilities
	if len(report.Vulnerabilities) > 0 {
		return false
	}
	return true
}

func convertToPurl(coord Coordinate) string {
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