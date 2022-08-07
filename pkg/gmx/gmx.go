package gmx

type Client struct {
	Hostname string
	Port     int
}

type MBeanOperator interface {
	GetString(domain string, beanName string, operation string, argName string) (string, error)
	PutString(domain string, name string, operation string, argName string, arvValue string) (string, error)
}
