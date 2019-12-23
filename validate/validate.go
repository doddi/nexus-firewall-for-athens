package validate

import (
	"nexus-firewall-for-athens/athens"
	"nexus-firewall-for-athens/ossindex"
)

func Validate(request athens.Request) bool {
	return ossindex.Validate(request.Module, request.Version)
}