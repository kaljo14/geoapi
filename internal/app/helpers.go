package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func pstr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

func pf64(v pgtype.Float8) *float64 {
	if !v.Valid {
		return nil
	}
	return &v.Float64
}

func pbool(v pgtype.Bool) *bool {
	if !v.Valid {
		return nil
	}
	return &v.Bool
}

func pint32(v pgtype.Int4) *int32 {
	if !v.Valid {
		return nil
	}
	i := int32(v.Int32)
	return &i
}

func ptime(v pgtype.Timestamptz) *string {
	if !v.Valid {
		return nil
	}
	s := v.Time.Format(time.RFC3339)
	return &s
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func f64Val(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
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
