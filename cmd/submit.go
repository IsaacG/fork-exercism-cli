package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Submit posts an iteration to the API.
func Submit(ctx *cli.Context) {
	if len(ctx.Args()) == 0 {
		log.Fatal("Please enter a file name")
	}

	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("verbose") {
		log.Printf("Exercises dir: %s", c.Dir)
		dir, err := os.Getwd()
		if err != nil {
			log.Printf("Unable to get current working directory - %s", err)
		} else {
			log.Printf("Current dir: %s", dir)
		}
	}

	if !c.IsAuthenticated() {
		log.Fatal(msgPleaseAuthenticate)
	}

	dir, err := filepath.EvalSymlinks(c.Dir)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("verbose") {
		log.Printf("eval symlinks (dir): %s", dir)
	}

	files := []string{}
	for _, filename := range ctx.Args() {
		if ctx.GlobalBool("verbose") {
			log.Printf("file name: %s", filename)
		}

		if isTest(filename) && !ctx.Bool("test") {
			log.Fatal("You're trying to submit a test file. If this is really what " +
				"you want, please pass the --test flag to exercism submit.")
		}

		file, err := filepath.Abs(filename)
		if err != nil {
			log.Fatal(err)
		}

		if ctx.GlobalBool("verbose") {
			log.Printf("absolute path: %s", file)
		}

		file, err = filepath.EvalSymlinks(file)
		if err != nil {
			log.Fatal(err)
		}

		if ctx.GlobalBool("verbose") {
			log.Printf("eval symlinks (file): %s", file)
		}

		files = append(files, file)
	}

	iteration, err := api.NewIteration(dir, files)
	if err != nil {
		log.Fatalf("Unable to submit - %s", err)
	}
	iteration.Key = c.APIKey

	client := api.NewClient(c)
	submission, err := client.Submit(iteration)
	if err != nil {
		log.Fatal(err)
	}

	msg := `
Submitted %s in %s.
Your submission can be found online at %s
`

	if submission.Iteration == 1 {
		msg += `
To get the next exercise, run "exercism fetch" again.
`
	}

	fmt.Printf(msg, submission.Name, submission.Language, submission.URL)
}
