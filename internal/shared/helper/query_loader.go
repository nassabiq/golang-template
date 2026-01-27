package helper

import "strings"

func LoadQuery(sql string) map[string]string {
	queries := make(map[string]string)

	blocks := strings.Split(sql, "-- name:")

	for _, block := range blocks[1:] {
		parts := strings.SplitN(block, "\n", 2)
		if len(parts) < 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		query := strings.TrimSpace(parts[1])

		queries[name] = query
	}

	return queries
}
