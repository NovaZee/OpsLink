package main

import (
	"fmt"
	"github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg"
	"github.com/denovo/permission/pkg/router"
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
	//定义config文件进行配置
	cfg, error := getCfg(c)
	if error != nil {
		return error
	}
	//初始化日志
	config.InitLoggerFromConfig(cfg.Logging)
	//初始化程序
	server, err := pkg.InitializeServer(cfg)
	if err != nil {
		return err
	}
	error = router.InitRouter(server.Casbin, cfg)
	if error != nil {
		return error
	}
	return nil
}

func getCfg(c *cli.Context) (*config.Config, error) {
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
