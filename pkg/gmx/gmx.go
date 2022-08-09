package gmx

import (
	"gmx/internal/mbean"
)

type Client struct {
	Hostname string
	Port     int
	client   *mbean.Client
}

type MBeanOperator interface {
	Initialize()
	Close()
	GetString(domain string, beanName string, operation string, argName string) (string, error)
	PutString(domain string, name string, operation string, argName string, arvValue string) (string, error)
}
