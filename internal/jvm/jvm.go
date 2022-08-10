package jvm

import (
	"errors"
	"gmx/internal/mbean"
	"runtime"

	"tekao.net/jnigi"
)

// Java is the structure that will contain JVM pertinent information.
type Java struct {
	Env     *jnigi.Env
	jvm     *jnigi.JVM
	started bool
	beans   []*mbean.Client
}

// IJava is the interface that wraps around the JVM.
// It allows for creation and cleanup. Only one JVM needs to be started.
// It can then be shared out between goroutines to do with as needed
type IJava interface {
	CreateJvm() (*Java, error)
	ShutdownJvm() error
}

// CreateJVM will create a JVM for the consumer to execute against
func CreateJVM() (*Java, error) {
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
