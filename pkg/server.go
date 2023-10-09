package pkg

import (
	config "github.com/denovo/permission/configration"
)

type OpsLinkServer struct {
	config *config.Config
}

func NewOpsLinkServer(config *config.Config) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		config: config,
	}
	return
}
