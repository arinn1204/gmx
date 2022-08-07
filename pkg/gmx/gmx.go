package gmx

type MBeanArg struct {
	Value     string
	ClassName string
}

type Client struct {
	Hostname string
	Port     int
}

type MBeanOperator interface {
	GetValue(domain string, name string, operation string, args ...MBeanArg) (string, error)
	PutValue(domain string, name string, operation string, args ...MBeanArg) (string, error)
}
