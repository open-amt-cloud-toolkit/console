package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	ErrUpgrade  = errors.New("upgrade error")
	ErrRedirect = errors.New("redirection error")
)

func TestWebSocketHandler(t *testing.T) { //nolint:paralleltest // logging library is not thread-safe for tests
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockFeature := NewMockFeature(ctrl)
	mockUpgrader := NewMockUpgrader(ctrl)
	mockLogger := &MockLogger{}

	tests := []struct {
		name           string
		upgraderError  error
		redirectError  error
		expectedStatus int
	}{
		{
			name:           "Success case",
			upgraderError:  nil,
			redirectError:  nil,
			expectedStatus: http.StatusSwitchingProtocols,
		},
		{
			name:           "Upgrade error",
			upgraderError:  ErrUpgrade,
			redirectError:  nil,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Redirect error",
			upgraderError:  nil,
			redirectError:  ErrRedirect,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // logging library is not thread-safe for tests
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgraderError != nil {
				mockUpgrader.EXPECT().
					Upgrade(gomock.Any(), gomock.Any(), nil).
					Return(nil, tc.upgraderError)
			} else {
				mockUpgrader.EXPECT().
					Upgrade(gomock.Any(), gomock.Any(), nil).
					Return(&websocket.Conn{}, nil)

				mockFeature.EXPECT().
					Redirect(gomock.Any(), gomock.Any(), "someHost", "someMode").
					Return(tc.redirectError)
			}

			r := gin.Default()
			RegisterRoutes(r, mockLogger, mockFeature, mockUpgrader)

			req := httptest.NewRequest(http.MethodGet, "/relay/webrelay.ashx?host=someHost&mode=someMode", http.NoBody)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
