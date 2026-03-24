package app

import (
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
