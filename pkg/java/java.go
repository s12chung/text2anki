// Package java contains Java integration helpers
package java

import (
	"fmt"

	"tekao.net/jnigi"
)

// Instance is an instance of java running
type Instance struct {
	JVM *jnigi.JVM
	Env *jnigi.Env
}

// Setup sets up the jvm instance
func (i *Instance) Setup(classPath string) error {
	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		return err
	}

	args := []string{
		"-Xcheck:jni",
		"-Djava.class.path=" + classPath,
	}

	var err error
	i.JVM, i.Env, err = jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, int(jnigi.DEFAULT_VERSION), args))
	if err != nil {
		return err
	}
	return nil
}

// Cleanup cleans up the java instance
func (i *Instance) Cleanup() error {
	err := i.JVM.Destroy()
	fmt.Println("") // space out the JVM stuff
	if err != nil {
		return err
	}
	i.JVM, i.Env = nil, nil
	return nil
}

// IsSetup returns true if the java instance is setup
func (i *Instance) IsSetup() bool {
	return i.JVM != nil && i.Env != nil
}

// JStringToString converts a java string to a string
func (i *Instance) JStringToString(jString interface{}) (string, error) {
	jStringBytes, err := jString.(*jnigi.ObjectRef).CallMethod(i.Env, "getBytes", jnigi.Byte|jnigi.Array)
	if err != nil {
		return "", err
	}
	return string(jStringBytes.([]byte)), nil
}
