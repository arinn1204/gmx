package jvm

import (
	"errors"
	"gmx/internal/mbean"
	"runtime"

	"tekao.net/jnigi"
)

type Java struct {
	Env     *jnigi.Env
	jvm     *jnigi.JVM
	started bool
	beans   []*mbean.MBean
}

type IJava interface {
	CreateJvm() (*Java, error)
	ShutdownJvm() error
}

// CreateJVM will create a JVM for the consumer to execute against
func CreateJvm() (*Java, error) {
	java := &Java{}

	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		return nil, errors.New("Failed to load the JVM::" + err.Error())
	}

	args := []string{"-Xcheck:jni"}

	runtime.LockOSThread()
	jvm, env, err := jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, jnigi.DEFAULT_VERSION, args))

	if err != nil {
		return nil, errors.New("Failed to create the JVM::" + err.Error())
	}

	configureEnvironment(env)

	java.jvm = jvm
	java.Env = env
	java.started = true

	return java, nil
}

// ShutdownJvm will shut down the JVM, this should be done at the end
func (java *Java) ShutdownJvm() error {
	if java == nil || java.jvm == nil {
		return nil
	}

	for _, bean := range java.beans {
		bean.JmxConnection.CallMethod(java.Env, "close", nil)
	}

	if err := java.jvm.Destroy(); err != nil {
		return err
	}

	java.jvm = nil
	java.Env = nil

	return nil
}

func configureEnvironment(env *jnigi.Env) {
	env.ExceptionHandler = jnigi.ThrowableToStringExceptionHandler
}
