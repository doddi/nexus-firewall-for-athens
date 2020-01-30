package validate

import "strings"

type Coordinate struct {
	Type       string `json:"type"`
	Namespace  string `json:"namespace,omitempty"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Qualifiers string `json:"qualifiers,omitempty"`
	Subpath    string `json:"subpath, omitempty"`
}

func ConvertToPurlString(coord Coordinate) string {
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
