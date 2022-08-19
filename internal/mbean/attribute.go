package mbean

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

// Get is the method that will fetch the attribute given by the name
// The return value must be able to be converted to the java representation that the server is expecting
// If this is a complex type, then a handler will need to be defined
func (mbean *Client) Get(domainName string, beanName string, attributeName string, args OperationArgs) (string, error) {
	mBeanServerConnector, err := createMBeanServerConnection(mbean.Env, mbean)

	if err != nil {
		return "", errors.New("failed to create the mbean server connection::" + err.Error())
	}

	defer mbean.Env.DeleteLocalRef(mBeanServerConnector)

	attribute, err := getAttribute(mBeanServerConnector, mbean.Env, Operation{
		Domain:    domainName,
		Name:      beanName,
		Operation: attributeName,
		Args:      []OperationArgs{args},
	})

	if err != nil {
		return "", err
	}

	return mbean.toGoString(mbean.Env, attribute)
}

func getAttribute(connection *jnigi.ObjectRef, env *jnigi.Env, operation Operation) (*jnigi.ObjectRef, error) {
	objectName, err := getObjectName(env, operation)
	if err != nil {
		return nil, errors.New("failed to create ObjectName::" + err.Error())
	}

	defer env.DeleteLocalRef(objectName)

	attributeNameRef, err := stringHandler.ToJniRepresentation(env, operation.Operation)

	if err != nil {
		return nil, err
	}
	defer env.DeleteLocalRef(attributeNameRef)

	attributeList := jnigi.NewObjectRef(OBJECT)
	if err = connection.CallMethod(env, "getAttribute", attributeList, objectName, attributeNameRef); err != nil {
		return nil, err
	}

	return attributeList, nil
}

func getObjectName(env *jnigi.Env, operation Operation) (*jnigi.ObjectRef, error) {
	mbeanName := fmt.Sprintf("%s:name=%s", operation.Domain, operation.Name)
	objectParam, err := stringHandler.ToJniRepresentation(env, mbeanName)

	defer env.DeleteLocalRef(objectParam)
	if err != nil {
		return nil, err
	}

	return env.NewObject("javax/management/ObjectName", objectParam)
}
