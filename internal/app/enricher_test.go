package app

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPlaceDetailsIntegerBounds(t *testing.T) {
	for _, field := range []string{"price_level", "user_ratings_total", "utc_offset"} {
		for _, tc := range []struct {
			value   string
			wantErr bool
		}{
			{value: "0"},
			{value: "120"},
			{value: "-300"},
			{value: "2147483647"},
			{value: "-2147483648"},
			{value: "2147483648", wantErr: true},
			{value: "-2147483649", wantErr: true},
		} {
			t.Run(field+"/"+tc.value, func(t *testing.T) {
				var response placeDetailsResponse
				err := json.Unmarshal([]byte(fmt.Sprintf(`{"status":"OK","result":{%q:%s}}`, field, tc.value)), &response)
				if (err != nil) != tc.wantErr {
					t.Fatalf("decode error = %v, want error = %v", err, tc.wantErr)
				}
			})
		}
	}
}
