package handler

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/mocks"
)

func TestImportPlacesUploadLimit(t *testing.T) {
	for _, tc := range []struct {
		name          string
		requestSize   int
		unknownLength bool
		wantStatus    int
	}{
		{name: "below limit", requestSize: 1024, wantStatus: http.StatusOK},
		{name: "at limit", requestSize: maxMultipartSize, wantStatus: http.StatusOK},
		{name: "above limit", requestSize: maxMultipartSize + 1, wantStatus: http.StatusRequestEntityTooLarge},
		{name: "above limit without content length", requestSize: maxMultipartSize + 1, unknownLength: true, wantStatus: http.StatusRequestEntityTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)
			part, err := writer.CreateFormFile("file", "places.csv")
			if err != nil {
				t.Fatal(err)
			}
			// Account for the closing boundary so requestSize is the entire body size.
			closingSize := len("\r\n--" + writer.Boundary() + "--\r\n")
			contents := strings.Repeat("x", tc.requestSize-body.Len()-closingSize)
			if _, err := io.WriteString(part, contents); err != nil {
				t.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			if body.Len() != tc.requestSize {
				t.Fatalf("body size = %d, want %d", body.Len(), tc.requestSize)
			}

			service := mocks.NewMockService(t)
			if tc.wantStatus == http.StatusOK {
				service.EXPECT().ImportCSV(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, r io.Reader) (int, error) {
						got, err := io.ReadAll(r)
						if err != nil || string(got) != contents {
							t.Fatal("import received incorrect file contents", err)
						}
						return 1, nil
					}).Once()
			}
			req := httptest.NewRequest(http.MethodPost, "/places/import", &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			if tc.unknownLength {
				req.ContentLength = -1
			}
			response := httptest.NewRecorder()
			NewHandler(service, zap.NewNop()).ImportPlaces(response, req)
			if response.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.wantStatus, response.Body.String())
			}
		})
	}
}
