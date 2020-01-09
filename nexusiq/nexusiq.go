package nexusiq

import (
	"fmt"
	"nexus-firewall-for-athens/athens"
)

type NexusIq struct{}

func (n NexusIq) Validate(request athens.Request) athens.Response {
	fmt.Println("Unimplemented")
	return athens.Response{Success: true}
}
