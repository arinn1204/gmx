package jvm

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

// the commonly used types
const (
	STRING  = "java/lang/String"
	OBJECT  = "java/lang/Object"
	LONG    = "java/lang/Long"
	INTEGER = "java/lang/Integer"
	BOOLEAN = "java/lang/Boolean"
	FLOAT   = "java/lang/Float"
	DOUBLE  = "java/lang/Double"
)

type Java struct {
	Env     *jnigi.Env
	jvm     *jnigi.JVM
	started bool
}

type IJava interface {
	CreateJvm() (*jnigi.Env, error)
	ShutdownJvm() error
}

// CreateJVM will create a JVM for the consumer to execute against
func CreateJvm() (*Java, error) {
	java := &Java{}

	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		return nil, errors.New("Failed to create a JVM::" + err.Error())
	}

	args := []string{"-Xcheck:jni"}

	jvm, env, err := jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, jnigi.DEFAULT_VERSION, args))

	if err != nil {
		return nil, errors.New("Failed to create the JVM::" + err.Error())
	}

	env.ExceptionHandler = jnigi.ThrowableToStringExceptionHandler

	java.jvm = jvm
	java.Env = env
	java.started = true

	return java, nil
}

// ShutdownJvm will shut down the JVM, this should be done at the end
func (java *Java) ShutdownJvm() error {
	if java.jvm == nil {
		return nil
	}

	if err := java.jvm.Destroy(); err != nil {
		return err
	}

	java.jvm = nil
	java.Env = nil

	return nil
}

func (java *Java) CreateString(str string) (*jnigi.ObjectRef, error) {
	fileNameRef, err := java.Env.NewObject(STRING, []byte(str))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", str, err.Error())
	}

	return fileNameRef, nil
}

func (java *Java) CreateJavaNative(obj any, typeName string) (*jnigi.ObjectRef, error) {
	ref, err := java.Env.NewObject(typeName, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", obj, err.Error())
	}

	return ref, nil
}
