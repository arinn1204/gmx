package mbean

import (
	"errors"
	"fmt"

	"github.com/arinn1204/gmx/internal/handlers"
	"tekao.net/jnigi"
)

// Put is the method that will set the declared attribute to the given value
// The value must be able to be converted to the java representation that the server is expecting
// If this is a custom type, then a handler will need to be defined
func (mbean *Client) Put(domainName string, beanName string, attributeName string, args OperationArgs) (string, error) {
	mBeanServerConnector, err := createMBeanServerConnection(mbean.Env, mbean)

	if err != nil {
		return "", errors.New("failed to create the mbean server connection::" + err.Error())
	}

	defer mbean.Env.DeleteLocalRef(mBeanServerConnector)

	var attributeType string
	var innerType string

	operation := Operation{
		Domain:    domainName,
		Name:      beanName,
		Operation: attributeName,
		Args:      []OperationArgs{args},
	}

	if args.JavaContainerType != "" {
		attributeType = args.JavaContainerType
		innerType = args.JavaType
	} else {
		innerType = ""
		attributeType, err = mbean.getAttributeType(mBeanServerConnector, mbean.Env, operation)
	}

	if err != nil {
		return "", err
	}

	objectName, err := getObjectName(mbean.Env, operation)
	if err != nil {
		return "", errors.New("failed to create ObjectName::" + err.Error())
	}

	defer mbean.Env.DeleteLocalRef(objectName)

	attribute, err := mbean.createAttributeJni(mbean.Env, operation, innerType, attributeType, args.Value)
	if err != nil {
		return "", err
	}

	defer mbean.Env.DeleteLocalRef(attribute)

	if err := mBeanServerConnector.CallMethod(mbean.Env, "setAttribute", nil, objectName, attribute); err != nil {
		return "", err
	}

	return "", nil
}

func (mbean *Client) createAttributeJni(env *jnigi.Env, operation Operation, innerType string, attributeType string, value string) (*jnigi.ObjectRef, error) {

	nameRef, err := stringHandler.ToJniRepresentation(env, operation.Operation)

	if err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(nameRef)

	attributeRef, err := toJni(mbean, attributeType, innerType, attributeType, value)

	if err != nil {
		return nil, err
	}

	defer mbean.Env.DeleteLocalRef(attributeRef)

	env.PrecalculateSignature("(Ljava/lang/String;Ljava/lang/Object;)V")
	return env.NewObject("javax/management/Attribute", nameRef, attributeRef)
}

func (mbean *Client) getAttributeType(connection *jnigi.ObjectRef, env *jnigi.Env, operation Operation) (string, error) {
	objectName, err := getObjectName(env, operation)
	if err != nil {
		return "", errors.New("failed to create ObjectName::" + err.Error())
	}

	defer env.DeleteLocalRef(objectName)

	beanInfo := jnigi.NewObjectRef("javax/management/MBeanInfo")

	if err := connection.CallMethod(env, "getMBeanInfo", beanInfo, objectName); err != nil {
		return "", errors.New("failed to get the mbean info::" + err.Error())
	}

	defer env.DeleteLocalRef(beanInfo)
	attributeInfoArr := jnigi.NewObjectArrayRef("javax/management/MBeanAttributeInfo")

	if err := beanInfo.CallMethod(env, "getAttributes", attributeInfoArr); err != nil {
		return "", err
	}

	defer env.DeleteLocalRef(attributeInfoArr)
	attributeInfos := env.FromObjectArray(attributeInfoArr)

	nameRef := jnigi.NewObjectRef(handlers.StringJniRepresentation)
	typeRef := jnigi.NewObjectRef(handlers.StringJniRepresentation)

	defer env.DeleteLocalRef(nameRef)
	defer env.DeleteLocalRef(typeRef)

	for _, attributeInfo := range attributeInfos {

		if err := attributeInfo.CallMethod(env, "getName", nameRef); err != nil {
			return "", err
		}

		var dest string
		if err = stringHandler.ToGoRepresentation(env, nameRef, &dest); err != nil {
			return "", err
		}

		if dest == operation.Operation {
			if err := attributeInfo.CallMethod(env, "getType", typeRef); err != nil {
				return "", err
			}

			if err = stringHandler.ToGoRepresentation(env, typeRef, &dest); err != nil {
				return "", err
			}

			return dest, nil
		}

		env.DeleteLocalRef(attributeInfo)
	}

	return "", fmt.Errorf("failed to find an attribute with a name of %s", operation.Operation)
}
