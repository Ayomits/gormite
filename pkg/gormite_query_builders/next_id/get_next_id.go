package next_id

import (
	"fmt"
	"iter"

	gqb "github.com/KoNekoD/gormite/pkg/gormite_query_builders"
)

type option struct {
	alias     string
	startFrom any
	limit     int
}

type OptionFn func(*option)

func WithAlias(alias string) OptionFn {
	return func(o *option) { o.alias = alias }
}

func WithStartFrom(startFrom any) OptionFn {
	return func(o *option) { o.startFrom = startFrom }
}

func WithLimit(limit int) OptionFn {
	return func(o *option) { o.limit = limit }
}

func getOptions(opts []OptionFn) (alias string, startFrom any, limit int) {
	options := &option{alias: "e", limit: 1000}

	for _, opt := range opts {
		opt(options)
	}

	return options.alias, options.startFrom, options.limit
}

// GetNextId - Helper function to get next id from the query builder through iteration.
func GetNextId[T int | string](qb *gqb.QueryBuilder[T], opts ...OptionFn) iter.Seq[T] {
	alias, startFrom, limit := getOptions(opts)

	column := fmt.Sprintf("%s.id", alias)

	return func(yield func(T) bool) {
		for {
			doQb := qb.Clone().Select(column).SetMaxResults(limit).OrderBy(column)

			if startFrom != nil {
				doQb.AndWhere(fmt.Sprintf("%s > @startFrom", column)).SetParameter("startFrom", startFrom)
			}

			ids, err := doQb.GetLiteralResult()
			if err != nil {
				panic(err)
			}

			fmt.Printf("got %d ids\n", len(ids))

			if len(ids) == 0 {
				break
			}

			startFrom = &ids[len(ids)-1]

			for _, id := range ids {
				if !yield(id) {
					return
				}
			}
		}
	}
}
