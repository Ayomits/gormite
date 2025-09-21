package next_id

import (
	"fmt"
	gqb "github.com/KoNekoD/gormite/pkg/gormite_query_builders"
	"github.com/pkg/errors"
	"iter"
)

type option struct {
	alias     string
	startFrom any
	limit     int
	onError   func(error)
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

func WithOnError(onError func(error)) OptionFn {
	return func(o *option) { o.onError = onError }
}

func defaultOnError(err error) {
	fmt.Println(errors.Wrap(err, "failed to get next id"))
}

func getOptions(opts []OptionFn) (alias string, startFrom any, limit int, onError func(error)) {
	options := &option{alias: "e", limit: 1000, onError: defaultOnError}

	for _, opt := range opts {
		opt(options)
	}

	return options.alias, options.startFrom, options.limit, options.onError
}

// GetNextId - Helper function to get next id from the query builder through iteration.
func GetNextId[T int | string](qb *gqb.QueryBuilder[T], opts ...OptionFn) iter.Seq[T] {
	alias, startFrom, limit, onError := getOptions(opts)

	column := fmt.Sprintf("%s.id", alias)

	return func(yield func(T) bool) {
		for {
			doQb := qb.Clone().Select(column).SetMaxResults(limit).OrderBy(column)

			if startFrom != nil {
				doQb.AndWhere(fmt.Sprintf("%s > @startFrom", column)).SetParameter("startFrom", startFrom)
			}

			ids, err := doQb.GetLiteralResult()
			if err != nil {
				onError(err)
			}

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
