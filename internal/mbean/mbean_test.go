package mbean

import (
	"testing"

	"github.com/arinn1204/gmx/internal/handlers/class"
	"github.com/arinn1204/gmx/pkg/extensions"
	"github.com/stretchr/testify/assert"
)

func TestRegisterClassHandler_SingleTypedTest(t *testing.T) {
	client := &Client{
		classHandlers: make(map[string]extensions.IClassHandler),
	}
	doubleHandler := class.DoubleHandler{}
	client.RegisterClassHandler(class.DOUBLE, &doubleHandler)

	assert.Equal(t, &doubleHandler, client.classHandlers[class.DOUBLE])
}

func TestRegisterClassHandler_MultipleTypedTest(t *testing.T) {
	client := &Client{
		classHandlers: make(map[string]extensions.IClassHandler),
	}
	doubleHandler := class.DoubleHandler{}
	floatHandler := class.FloatHandler{}
	client.RegisterClassHandler(class.DOUBLE, &doubleHandler)
	client.RegisterClassHandler(class.FLOAT, &floatHandler)

	assert.Equal(t, &doubleHandler, client.classHandlers[class.DOUBLE])
	assert.Equal(t, &floatHandler, client.classHandlers[class.FLOAT])
}
