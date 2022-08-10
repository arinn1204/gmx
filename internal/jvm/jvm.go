package jvm

import (
	"errors"
	"gmx/internal/mbean"
	"log"
	"runtime"

	"tekao.net/jnigi"
)

// Java is the structure that will contain JVM pertinent information.
type Java struct {
	Env     *jnigi.Env
	jvm     *jnigi.JVM
	started bool
}

// IJava is the interface that wraps around the JVM.
// It allows for creation and cleanup. Only one JVM needs to be started.
// It can then be shared out between goroutines to do with as needed
type IJava interface {
	CreateJVM() (*Java, error)                               // Will create and start the JVM for any JNI threads to communicate with
	ShutdownJvm() error                                      // Will cleanup any threads remaining and close the JVM
	Attach() *jnigi.Env                                      // Will attach the current running thread to the JVM
	IsStarted() bool                                         // A simple flag indicating whether or not the JVM has started running
	CreateMBeanConnection(uri string) (*mbean.Client, error) // The factory method for all bean connections
}

// Attach is the method to attach the current thread to the JNI environment.
// This is required in order to execute any JNI commands in the actively running threads
// and should be used whenever a new go routine is used
func (java *Java) Attach() *jnigi.Env {
	runtime.LockOSThread()
	return java.jvm.AttachCurrentThread()
}

// IsStarted will indicate whether or not the JVM has already been started
func (java *Java) IsStarted() bool {
	return java.started
}

// CreateJVM will create a JVM for the consumer to execute against
func (java *Java) CreateJVM() (*Java, error) {

	if java.IsStarted() {
		log.Fatal("The JVM has already been started.")
		return nil, nil
	}

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
