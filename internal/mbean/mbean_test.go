package mbean

import (
	"testing"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/arinn1204/gmx/pkg/extensions"
	"github.com/stretchr/testify/assert"
)

func TestRegisterClassHandler_SingleTypedTest(t *testing.T) {
	client := &Client{
		classHandlers: make(map[string]extensions.IClassHandler),
	}
	doubleHandler := handlers.DoubleHandler{}
	client.RegisterClassHandler(handlers.DOUBLE, &doubleHandler)

	assert.Equal(t, &doubleHandler, client.classHandlers[handlers.DOUBLE])
}

func TestRegisterClassHandler_MultipleTypedTest(t *testing.T) {
	client := &Client{
		classHandlers: make(map[string]extensions.IClassHandler),
	}
	doubleHandler := handlers.DoubleHandler{}
	floatHandler := handlers.FloatHandler{}
	client.RegisterClassHandler(handlers.DOUBLE, &doubleHandler)
	client.RegisterClassHandler(handlers.FLOAT, &floatHandler)

	assert.Equal(t, &doubleHandler, client.classHandlers[handlers.DOUBLE])
	assert.Equal(t, &floatHandler, client.classHandlers[handlers.FLOAT])
}
