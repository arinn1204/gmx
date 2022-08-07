package java

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	"tekao.net/jnigi"
)

// CreateJVM will create a JVM for the consumer to execute against
func CreateJvm() (*Java, error) {
	//get a lock to ensure you are the only one trying to get the JVM started
	java := &Java{
		lock: &sync.Mutex{},
	}

	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		return nil, errors.New("Failed to create a JVM::" + err.Error())
	}

	runtime.LockOSThread()

	args := []string{"-Xcheck:jni"}

	jvm, env, err := jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, jnigi.DEFAULT_VERSION, args))

	if err != nil {
		return nil, errors.New("Failed to create the JVM::" + err.Error())
	}

	java.jvm = jvm
	java.env = env
	java.started = true

	return java, nil
}

func (java *Java) IsStarted() bool {
	return java.started
}

func (java *Java) ShutdownJvm() error {
	java.lock.Lock()
	defer java.lock.Unlock()

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
