package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type errorResponseTest struct {
	name           string
	code           int
	msg            string
	expectedStatus int
	expectedBody   string
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	tests := []errorResponseTest{
		{
			name:           "ErrorResponse with 400 status",
			code:           http.StatusBadRequest,
			msg:            "bad request",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"bad request"}`,
		},
		{
			name:           "ErrorResponse with 500 status",
			code:           http.StatusInternalServerError,
			msg:            "internal server error",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)

			errorResponse(c, tc.code, tc.msg)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
