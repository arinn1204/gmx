package mbean

import (
	"testing"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/arinn1204/gmx/pkg/extensions"
	"github.com/stretchr/testify/assert"
)

func TestRegisterClassHandler_SingleTypedTest(t *testing.T) {
	client := &Client{
		classHandlers: make(map[string]extensions.IHandler),
	}
	doubleHandler := handlers.DoubleHandler{}
	client.RegisterClassHandler(handlers.JNI_DOUBLE, &doubleHandler)

	assert.Equal(t, &doubleHandler, client.classHandlers[handlers.JNI_DOUBLE])
}

func TestRegisterClassHandler_MultipleTypedTest(t *testing.T) {
	client := &Client{
		classHandlers: make(map[string]extensions.IHandler),
	}
	doubleHandler := handlers.DoubleHandler{}
	floatHandler := handlers.FloatHandler{}
	client.RegisterClassHandler(handlers.JNI_DOUBLE, &doubleHandler)
	client.RegisterClassHandler(handlers.JNI_FLOAT, &floatHandler)

	assert.Equal(t, &doubleHandler, client.classHandlers[handlers.JNI_DOUBLE])
	assert.Equal(t, &floatHandler, client.classHandlers[handlers.JNI_FLOAT])
}
