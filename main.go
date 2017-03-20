package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/glesica/farnsworth/project"
	_ "github.com/glesica/farnsworth/proxy/golang"
	_ "github.com/glesica/farnsworth/proxy/java"
	"github.com/glesica/farnsworth/proxy"
	"github.com/glesica/farnsworth/ignore"
)

func main() {
	app := cli.NewApp()

	app.Name = "Farnsworth"
	app.HelpName = "farnsworth"
	app.Usage = "Create and evaluate programming assignments"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		{
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

				projectPath := c.String("project")

				proj, projErr := project.Load(projectPath, ignore.Get, proxy.Get)
				if projErr != nil {
					return cli.NewExitError("Failed to load project.", 10)
				}

				proj.Zip(c.Args().Get(0), c.Bool("public"))
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "public",
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
		{
			Name:    "merge",
			Aliases: []string{"m"},
			Usage:   "Merge another project into the specified project.",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return cli.NewExitError("No project to merged specified.", 10)
				}

				projectPath := c.String("project")

				proj, projErr := project.Load(projectPath, ignore.Get, proxy.Get)
				if projErr != nil {
					return cli.NewExitError("Failed to load project.", 10)
				}

				mergeProjectPath := c.Args().Get(0)

				mergeProj, mergeProjErr := project.Load(mergeProjectPath, ignore.Get, proxy.Get)
				if mergeProjErr != nil {
					return cli.NewExitError("Failed to load merge target project.", 10)
				}

				mergeErr := proj.MergeFrom(*mergeProj)
				if mergeErr != nil {
					return mergeErr
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "project",
					Usage: "Path to the base project `ROOT`",
					Value: ".",
				},
			},
			ArgsUsage: "[merge source]",
		},
	}

	app.Run(os.Args)
}
