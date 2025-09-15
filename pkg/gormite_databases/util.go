package gormite_databases

import "strings"

func trimSQL(sql string) string {
	lines := strings.Split(sql, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	sql = strings.Join(lines, " ")

	return sql
}
