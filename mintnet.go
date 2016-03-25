package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const ValSetAnon = "anon"

var (
	machFlag = cli.StringFlag{
		Name:  "machines",
		Value: "mach[1-4]",
		Usage: "Comma separated list of machine names",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "mintnet"
	app.Usage = "mintnet [command] [args...]"
	app.Commands = []cli.Command{
		{
			Name:      "info",
			Usage:     "Information about running containers",
			ArgsUsage: "[appName]",
			Action: func(c *cli.Context) {
				cmdInfo(c)
			},
			Flags: []cli.Flag{machFlag},
			Subcommands: []cli.Command{
				{
					Name:      "port",
					Usage:     "Print container port mapping",
					ArgsUsage: "[appName]",
					Action: func(c *cli.Context) {
						cmdPorts(c)
					},
				},
			},
		},

		{
			Name:      "init",
			Usage:     "Initialize node configuration directories",
			ArgsUsage: "[baseDir]",
			Flags: []cli.Flag{
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdInit(c)
			},
			Subcommands: []cli.Command{
				{
					Name:      "chain",
					Usage:     "Initialize a new blockchain",
					ArgsUsage: "[baseDir]",
					Action: func(c *cli.Context) {
						cmdChainInit(c)
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "validator-set",
							Value: "",
							Usage: "Specify a path to the validator set for the new chain",
						},
						cli.StringFlag{
							Name:  "app",
							Value: "",
							Usage: "Specify an init.sh file for the app to run",
						},
						cli.StringFlag{
							Name:  "app-hash",
							Value: "",
							Usage: "Specify the app's initial hash. Prefix with 0x if it's hex",
						},
					},
				},
				{
					Name:      "validator-set",
					Usage:     "Initialize a new validator set",
					ArgsUsage: "[baseDir]",
					Action: func(c *cli.Context) {
						cmdValidatorsInit(c)
					},
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "N",
							Value: 4,
							Usage: "Size of the validator set",
						},
					},
				},
			},
		},

		{
			Name:      "create",
			Usage:     "Create a new Tendermint network with newly provisioned machines. Use -- to pass args through to docker-machine",
			ArgsUsage: "",
			Flags: []cli.Flag{
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdCreate(c)
			},
		},

		{
			Name:      "provision",
			Usage:     "Provision already created machines (useful if the create failed)",
			ArgsUsage: "",
			Flags: []cli.Flag{
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdProvision(c)
			},
		},

		{
			Name:      "destroy",
			Usage:     "Destroy a Tendermint network",
			ArgsUsage: "",
			Flags: []cli.Flag{
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdDestroy(c)
			},
		},

		{
			Name:      "start",
			Usage:     "Start blockchain application",
			ArgsUsage: "[appName] [baseDir]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "seed-machines",
					Value: "",
					Usage: "Comma separated list of machine names for seed, defaults to --machines",
				},
				cli.BoolFlag{
					Name:  "publish-all,P",
					Usage: "Publish all exposed ports to random ports",
				}, // or should we make random be default, and let users attempt to force the port?
				cli.BoolFlag{
					Name:  "no-tmsp",
					Usage: "Use a null, in-process app",
				},
				cli.StringFlag{
					Name:  "tmcore-image",
					Usage: "Docker image to use for tendermint core",
					Value: "tendermint/tmbase",
				},
				cli.StringFlag{
					Name:  "tmapp-image",
					Usage: "Docker image to use for tendermint app",
					Value: "tendermint/tmbase",
				},
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdStart(c)
			},
		},

		{
			Name:      "restart",
			Usage:     "Re start a stopped blockchain application",
			ArgsUsage: "[appName]",
			Flags: []cli.Flag{
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdRestart(c)
			},
		},

		{
			Name:      "stop",
			Usage:     "Stop blockchain application",
			ArgsUsage: "[appName]",
			Flags: []cli.Flag{
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdStop(c)
			},
		},

		{
			Name:  "rm",
			Usage: "Remove blockchain application",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force",
					Usage: "Force stop app if already running",
				},
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdRm(c)
			},
		},

		{
			Name:  "docker",
			Usage: "Execute a docker command on all machines",
			Flags: []cli.Flag{
				machFlag,
			},
			Action: func(c *cli.Context) {
				cmdDocker(c)
			},
		},
	}
	app.Run(os.Args)

}

//--------------------------------------------------------------------------------
