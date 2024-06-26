package toolbox

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadSQLQueries loads SQL queries from a file and populates the QUERY map.
// This tool aims to facilitate the use of the go language's database/sql standard library.
// Writing SQL queries directly in the code can make it messy, so writing SQL queries in .sql files
// and then calling them from the code helps prevent code clutter,
// allowing SQL queries to be centralized in one place for better organization.
func (t *Tools) LoadSQLQueries(fileName string) (map[string]string, error) {
	query := make(map[string]string)

	file, err := os.Open(fileName)
	if err != nil {
		return query, err
	}
	defer func() {
		_ = file.Close()
	}()

	query, err = parseSQLQueries(file, query)
	return query, err
}

// parseSQLQueries reads the SQL queries from the provided file and populates the QUERY map.
func parseSQLQueries(file *os.File, query map[string]string) (map[string]string, error) {
	scanner := bufio.NewScanner(file)
	var key string
	var queryBuilder strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if isSQLQuery(line) || len(key) > 0 {
			if len(key) > 0 {
				if strings.HasSuffix(line, ";") {
					queryBuilder.WriteString(line)
					query[key] = queryBuilder.String()
					key, queryBuilder = "", strings.Builder{}
				} else {
					queryBuilder.WriteString(line + " ")
				}
			} else {
				key = extractKey(line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return query, fmt.Errorf("error reading file: %w", err)
	}
	return query, nil
}

// isSQLQuery checks if the given line is an SQL query or a comment.
func isSQLQuery(line string) bool {
	return hasPrefixInList(line, []string{"-- ", "SELECT", "INSERT", "UPDATE", "DELETE"})
}

// extractKey extracts the key from the comment line.
func extractKey(line string) string {
	if strings.HasPrefix(line, "-- ") {
		return strings.Split(line, "-- ")[1]
	}
	return ""
}

// hasPrefixInList is a prefix checker
func hasPrefixInList(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}
