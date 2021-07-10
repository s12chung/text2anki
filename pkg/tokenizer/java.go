package tokenizer

import "tekao.net/jnigi"

type javaInstance struct {
	jvm *jnigi.JVM
	env *jnigi.Env
}

func (j *javaInstance) setup(classPath string) error {
	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		return err
	}

	args := []string{
		"-Xcheck:jni",
		"-Djava.class.path=" + classPath,
	}

	var err error
	j.jvm, j.env, err = jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, int(jnigi.DEFAULT_VERSION), args))
	if err != nil {
		return err
	}
	return nil
}

// Cleanup cleans up the java instance
func (j *javaInstance) Cleanup() error {
	if err := j.jvm.Destroy(); err != nil {
		return err
	}
	j.jvm, j.env = nil, nil
	return nil
}

// IsSetup returns true if the java instance is setup
func (j *javaInstance) IsSetup() bool {
	return j.jvm != nil && j.env != nil
}

func (j *javaInstance) jStringToString(jString interface{}) (string, error) {
	jStringBytes, err := jString.(*jnigi.ObjectRef).CallMethod(j.env, "getBytes", jnigi.Byte|jnigi.Array)
	if err != nil {
		return "", err
	}
	return string(jStringBytes.([]byte)), nil
}
