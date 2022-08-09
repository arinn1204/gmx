package java

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

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
	java.env = env
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
	java.env = nil

	return nil
}

func (java *Java) createString(str string) (*jnigi.ObjectRef, error) {
	fileNameRef, err := java.env.NewObject(STRING, []byte(str))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", str, err.Error())
	}

	return fileNameRef, nil
}

func (java *Java) createLong(obj int64) (*jnigi.ObjectRef, error) {
	fileNameRef, err := java.env.NewObject(LONG, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to turn %d into an object::%s", obj, err.Error())
	}

	return fileNameRef, nil
}

func (java *Java) createInteger(obj int) (*jnigi.ObjectRef, error) {
	fileNameRef, err := java.env.NewObject(INTEGER, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to turn %d into an object::%s", obj, err.Error())
	}

	return fileNameRef, nil
}
