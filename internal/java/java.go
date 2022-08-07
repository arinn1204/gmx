package java

import (
	"errors"
	"runtime"
	"sync"

	"tekao.net/jnigi"
)

type Java struct {
	jvm     *jnigi.JVM
	env     *jnigi.Env
	lock    *sync.Mutex
	started bool
}

type IJava interface {
	CreateJvm() (*jnigi.Env, error)
	ShutdownJvm() error
	IsStarted() bool
}

var JavaRef *Java

func init() {
	JavaRef = &Java{
		jvm:     nil,
		env:     nil,
		lock:    &sync.Mutex{},
		started: false,
	}
}

// CreateJVM will create a JVM for the consumer to execute against
func (java *Java) CreateJvm() (*jnigi.Env, error) {
	//get a lock to ensure you are the only one trying to get the JVM started
	java.lock.Lock()
	defer java.lock.Unlock()

	if java.jvm != nil {
		return java.env, nil
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

	return env, nil
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
