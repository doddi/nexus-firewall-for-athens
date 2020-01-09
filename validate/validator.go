package validate

import "nexus-firewall-for-athens/athens"

type Validator interface {
	Validate(request athens.Request) athens.Response
}
