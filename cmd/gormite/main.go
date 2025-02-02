package main

import (
	"github.com/KoNekoD/gormite/pkg/gormite/runners"
	"github.com/gookit/goutil/cflag"
	"github.com/pkg/errors"
	"slices"
)

func toolValidate(val any) (err error) {
	if !slices.Contains([]string{"migrate", "goose"}, val.(string)) {
		err = errors.New("invalid migration tool name")
	}
	return err
}

func main() {
	opts := runners.DiffRunnerOptions{}

	c := cflag.New(func(c *cflag.CFlags) { c.Desc = "Create migrations diff" })

	c.StringVar(
		&opts.Tool,
		"tool",
		"",
		"migration tool, allowed: migrate, goose;true;t,mt,m",
	)
	c.StringVar(
		&opts.Dsn,
		"dsn",
		"",
		"database connection string;true;dsn,db,d",
	)

	c.AddValidator("tool", toolValidate)

	c.Func = func(c *cflag.CFlags) error { return runners.NewDiffRunner(opts).Run() }

	c.MustParse(nil)
}
