package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggo/swag"
)

func TestSwaggerInfoInitialization(t *testing.T) {
	// Assert that the SwaggerInfo is registered correctly
	spec := swag.GetSwagger(SwaggerInfo.InstanceName())

	//assert.True(t, ok, "SwaggerInfo should be registered")
	assert.Equal(t, SwaggerInfo, spec, "SwaggerInfo should match the registered spec")

	// Assert specific fields
	assert.Equal(t, "1.0", SwaggerInfo.Version)
	assert.Equal(t, "localhost:8080", SwaggerInfo.Host)
	assert.Equal(t, "/v1", SwaggerInfo.BasePath)
	assert.Equal(t, "Console API", SwaggerInfo.Title)
	assert.Equal(t, "Console is an application that provides a 1:1, direct connection for AMT devices for use in an enterprise environment. Users can add activated AMT devices to access device information and device management functionality such as power control, remote keyboard-video-mouse (KVM) control, and more.", SwaggerInfo.Description)
	assert.Equal(t, "{{", SwaggerInfo.LeftDelim)
	assert.Equal(t, "}}", SwaggerInfo.RightDelim)
	assert.Empty(t, SwaggerInfo.Schemes, "Schemes should be empty by default")
	assert.Equal(t, docTemplate, SwaggerInfo.SwaggerTemplate)
}
