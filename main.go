package main

import (
	"fmt"
	"github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg"
	"github.com/denovo/permission/pkg/kubeclient"
	"github.com/denovo/permission/pkg/router"
	"github.com/oppslink/protocol/logger"
	"github.com/urfave/cli/v2"
	"os"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "config",
		Usage: "path to OpsLink config file",
	},
}

func main() {
	app := &cli.App{
		Name:    "OpsLink",
		Usage:   "Permission Control",
		Flags:   flags,
		Action:  start,
		Version: "v1.0.0",
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func start(c *cli.Context) error {
	var err error
	//load config file
	cfg, error := getCfg(c)
	if error != nil {
		return error
	}
	//init logger
	config.InitLoggerFromConfig(cfg.Logging)
	//init oppslink server
	server, err := pkg.InitializeServer(cfg)
	if err != nil {
		return err
	}
	client, err := kubeclient.NewK8sClient(cfg, "kubernetes")
	client.Get()

	logger.Infow("start server ", "port", cfg.Server.HttpPort)

	//init http router
	_, err = router.InitRouter(server)
	if err != nil {
		return err
	}
	return nil
}

func getCfg(c *cli.Context) (*config.OpsLinkConfig, error) {
	confString, err := getConfigString(c.String("config"), c.String("config-body"))
	if err != nil {
		return nil, err
	}

	strictMode := true
	if c.Bool("disable-strict-config") {
		strictMode = false
	}
	conf, err := config.NewConfig(confString, strictMode)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func getConfigString(configFile string, inConfigBody string) (string, error) {
	if inConfigBody != "" || configFile == "" {
		return inConfigBody, nil
	}

	outConfigBody, err := os.ReadFile(configFile)
	if err != nil {
		return "", err
	}

	return string(outConfigBody), nil
}
