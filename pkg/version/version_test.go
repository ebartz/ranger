package version

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionServeHTTP(t *testing.T) {
	tests := []struct {
		name     string
		setPrime func()
		cleanup  func()
		want     string
	}{
		{
			name:     "unmodified",
			setPrime: func() {},
			cleanup:  func() {},
			want:     `{"Version":"dev","GitCommit":"HEAD","RangerPrime":"false"}`,
		},
		{
			name:     "prime=true",
			setPrime: func() { os.Setenv("RANCHER_PRIME", "true") },
			cleanup:  func() { os.Unsetenv("RANCHER_PRIME") },
			want:     `{"Version":"dev","GitCommit":"HEAD","RangerPrime":"true"}`,
		},
		{
			name:     "prime=false",
			setPrime: func() { os.Setenv("RANCHER_PRIME", "false") },
			cleanup:  func() { os.Unsetenv("RANCHER_PRIME") },
			want:     `{"Version":"dev","GitCommit":"HEAD","RangerPrime":"false"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setPrime()
			defer tt.cleanup()
			req := httptest.NewRequest(http.MethodGet, "/rangerversion", nil)
			rr := httptest.NewRecorder()
			handler := NewVersionHandler()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
			resp := rr.Result()
			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, string(body))
		})
	}
}
