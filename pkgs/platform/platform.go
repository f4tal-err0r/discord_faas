package platform

import (
	spec "k8s.io/api/core/v1"
)

type PlatformOpts struct {
	GuildID string
	
}

type Platform interface {
	CreateRunner
}
