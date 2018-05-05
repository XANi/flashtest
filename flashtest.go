package main

import (
	"github.com/op/go-logging"
	"github.com/urfave/cli"
	"os"
	"sort"
)

var version string
var log = logging.MustGetLogger("main")
var stdout_log_format = logging.MustStringFormatter("%{color:bold}%{time:2006-01-02T15:04:05.0000Z-07:00}%{color:reset}%{color} [%{level:.1s}] %{color:reset}%{shortpkg}[%{longfunc}] %{message}")

func main() {
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stderrFormatter := logging.NewBackendFormatter(stderrBackend, stdout_log_format)
	logging.SetBackend(stderrFormatter)
	logging.SetFormatter(stdout_log_format)
	app := cli.NewApp()
	app.Name = "foobar"
	app.Description = "do foo to bar"
	app.Version = version
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "help, h", Usage: "show help"},
		cli.StringFlag{
			Name: "device,d",
			//Value:  "",
			Usage: "Device, or file, to write to",
			//EnvVar: "_DEVICE",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("help") {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		log.Infof("Starting app version: %s", version)
		log.Infof("var example %s", c.GlobalString("url"))
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "write",
			Aliases: []string{"w"},
			Usage:   "Write test pattern into the flash",
			Action: func(c *cli.Context) error {
				log.Warning("running example cmd")
				return nil
			},
		},
		{
			Name:    "verify",
			Aliases: []string{"v"},
			Usage:   "Verify already written flash",
			Action: func(c *cli.Context) error {
				log.Warning("running example cmd")
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "write and then verify flash content",
			Action: func(c *cli.Context) error {
				log.Warning("running example cmd")
				return nil
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
