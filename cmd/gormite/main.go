package main

import (
	"context"
	"os"
	"slices"

	"github.com/KoNekoD/gormite/pkg/runners"
	"github.com/gookit/goutil/cflag"
	"github.com/pkg/errors"
)

func toolValidate(val any) (err error) {
	if !slices.Contains([]string{"migrate", "goose"}, val.(string)) {
		err = errors.New("invalid migration tool name")
	}
	return err
}

var possibleScenarios = []string{runners.ScenarioTypeDiff, runners.ScenarioTypeValidate}

func main() {
	scenario := runners.ScenarioTypeDiff
	args := os.Args[1:]
	if len(args) > 0 && slices.Contains(possibleScenarios, args[0]) {
		scenario = args[0]
		args = args[1:]
	}

	c := cflag.New(func(c *cflag.CFlags) { c.Desc = "Gormite CLI" })
	ctx := context.Background()
	opts := runners.DiffRunnerOptions{Scenario: scenario}

	switch scenario {
	case runners.ScenarioTypeDiff:
		c.StringVar(&opts.Tool, "tool", "", "migration tool, allowed: migrate, goose;true;t")
		c.AddValidator("tool", toolValidate)
	}

	c.StringVar(&opts.Dsn, "dsn", "", "database connection string;true")
	c.StringVar(&opts.ConfigPath, "config", "gormite.yaml", "config file path;true;config,c")

	c.Func = func(c *cflag.CFlags) error { return runners.NewDiffRunner(opts).Run(ctx) }

	c.MustParse(args)
}
