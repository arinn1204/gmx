package java

import (
	"errors"
	"runtime"

	"tekao.net/jnigi"
)

var jvm *jnigi.JVM

// CreateJVM will create a JVM for the consumer to execute against
func CreateJvm() error {
	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		return errors.New("Failed to create a JVM::" + err.Error())
	}

	runtime.LockOSThread()
	var err error

	jvm, _, err = jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, jnigi.DEFAULT_VERSION, []string{"-Xcheck:jni"}))

	if err != nil {
		return errors.New("Failed to create the JVM::" + err.Error())
	}

	return nil
}

func ShutdownJvm() {
	jvm.Destroy()
}
