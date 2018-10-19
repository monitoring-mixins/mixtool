package main

import (
	"fmt"
	"os"

	"github.com/metalmatze/mixtool/pkg/mixer"
	"github.com/urfave/cli"
)

func testCommand() cli.Command {
	return cli.Command{
		Name:        "test",
		Usage:       "Run unit tests for alerts",
		Description: "Run the embedded unit tests inside alert objects",
		Action:      testAction,
	}
}

func testAction(c *cli.Context) error {
	files := c.Args()

	for _, filename := range files {
		if err := mixer.Test(os.Stdout, filename, mixer.TestOptions{JPaths: []string{"vendor"}}); err != nil {
			return fmt.Errorf("failed to test the alerts in %s: %v", filename, err)
		}
	}

	return nil
}
