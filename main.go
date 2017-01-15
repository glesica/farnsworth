package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Farnsworth"
	app.HelpName = "farnsworth"
	app.Usage = "Create and evaluate programming assignments"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "George Lesica",
			Email: "george@lesica.com",
		},
	}
	app.Copyright = "(c) 2017 George Lesica"
	app.Commands = []cli.Command{
		{
			Name:    "archive",
			Aliases: []string{"a"},
			Usage:   "Bundle a project into an archive",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return cli.NewExitError("No archive path provided.", 10)
				}
				proj := LoadProj(c.String("project"))
				proj.Zip(c.Args().Get(0))
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "private",
					Usage: "Remove hidden content before archive is created",
				},
				cli.StringFlag{
					Name:  "project",
					Usage: "Path to the project `ROOT`",
					Value: ".",
				},
			},
			ArgsUsage: "[archive]",
		},
	}

	app.Run(os.Args)
}
