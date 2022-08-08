package gmx

import "gmx/internal/java"

type Client struct {
	Hostname string
	Port     int
	jvm      *java.Java
}

type MBeanOperator interface {
	Initialize()
	Close()
	GetString(domain string, beanName string, operation string, argName string) (string, error)
	PutString(domain string, name string, operation string, argName string, arvValue string) (string, error)
}
