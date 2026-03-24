package app

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func envFloat(key string, defaultVal float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultVal
}

func joinStrings(ss []string) string {
	return strings.Join(ss, ",")
}

func f64str(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }
func i32str(i int32) string   { return strconv.Itoa(int(i)) }

func boolstr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func timestr(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}
