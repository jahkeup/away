package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "away"
	app.Author = "Jacob Vallejo"
	app.Email = "jake@jahkeup.com"
	app.Usage = "Put away files from one directory to another"
	app.ArgsUsage = "[source...]"
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "dry-run,n", Usage: "Print out operations without running"},
		cli.BoolFlag{Name: "deep-link,l", Usage: "Create parent directories and link only files"},
		cli.BoolFlag{Name: "put,P", Usage: "Put files away in destination, default behavior"},
		cli.BoolFlag{Name: "rm,R", Usage: "Remove links from destination"},
		cli.StringFlag{Name: "dest,d", Usage: "Destination path", Value: ".."},
		cli.StringFlag{Name: "dir-mode,m", Usage: "Directory creation mode", Value: "0770"},
	}
	app.Action = func(c *cli.Context) {
		// Setup plan options to configure run
		options := PlanOptions{}
		options.DryRun = c.Bool("dry-run")
		options.LinkFilesOnly = c.Bool("deep-link")
		options.DirMode = parseFileMode(c.String("dir-mode"))

		// Establish plan executor
		var executor Executor
		executor = Putaway{}
		if c.Bool("put") {
			// Put files
			executor = Putaway{}
		} else if c.Bool("rm") {
			// Unput files
			executor = Unputaway{}
		}

		// Run for all sources given
		for _, source := range c.Args() {
			err := Process(source, c.String("dest"), executor, options)
			if err != nil {
				log.Error(err)
				return // bail on errors
			}
		}
	}
	app.Run(os.Args)
}

// Run the operation for a source that will execute using dest as the
// target destination
func Process(source, dest string, exec Executor, options PlanOptions) (err error) {
	plan, err := NewPlan(source, dest, options)
	if err != nil {
		return err
	}
	err = plan.FindNodes()
	if err != nil {
		return err
	}

	if options.DryRun {
		fmt.Println(exec.Describe(plan))
		return
	}
	return exec.Execute(plan)
}

// Parse user string for filemode
func parseFileMode(modeStr string) os.FileMode {
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		return 0770
	}
	return os.FileMode(mode)
}
