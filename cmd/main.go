package main

import (
	"fmt"
	"log"
	"os"

	info "github.com/christophercampbell/bridge-connector"
	"github.com/christophercampbell/bridge-connector/config"
	"github.com/christophercampbell/bridge-connector/db"
	"github.com/christophercampbell/bridge-connector/indexer"
	"github.com/urfave/cli/v2"
)

const appName = "bridge-connector"

var (
	cfgFlag = cli.StringFlag{
		Name:     config.ConfigFileFlagName,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: true,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = info.Version
	app.Commands = []*cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   fmt.Sprintf("Run the %v", appName),
			Action:  run,
			Flags:   []cli.Flag{&cfgFlag},
		},
		{
			Name:    "version",
			Aliases: []string{},
			Usage:   "Show version",
			Action: func(c *cli.Context) error {
				info.PrintVersion(os.Stderr)
				return nil
			},
			Flags: []cli.Flag{},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(cliCtx *cli.Context) error {
	cfg, err := config.Load(cliCtx)
	if err != nil {
		panic(err)
	}

	store, err := db.NewStorage(cfg.DB.File)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	lxService, err := indexer.New(cfg.LX, store)
	if err != nil {
		panic(err)
	}
	defer lxService.Stop()
	lxService.Start()

	lyService, err := indexer.New(cfg.LY, store)
	if err != nil {
		panic(err)
	}
	defer lyService.Stop()
	lyService.Start()

	return nil
}
