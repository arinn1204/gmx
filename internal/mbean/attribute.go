package mbean

// Get is the method that will fetch the attribute given by the name
// The return value must be able to be converted to the java representation that the server is expecting
// If this is a complex type, then a handler will need to be defined
func (mbean *Client) Get(domainName string, beanName string, attributeName string, args ...OperationArgs) (string, error) {
	return "", nil
}

// Put is the method that will set the declared attribute to the given value
// The value must be able to be converted to the java representation that the server is expecting
// If this is a custom type, then a handler will need to be defined
func (mbean *Client) Put(domainName string, beanName string, attributeName string, args ...OperationArgs) (string, error) {
	return "", nil
}
