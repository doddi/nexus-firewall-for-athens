package nexusiq

import (
	"fmt"
	"nexus-firewall-for-athens/athens"
)

type NexusIq struct{}

func (n NexusIq) Validate(request athens.Request) bool {
	fmt.Println("Unimplemented")
	return false
}
