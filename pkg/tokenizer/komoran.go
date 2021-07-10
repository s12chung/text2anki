package tokenizer

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"tekao.net/jnigi"
)

// NewKomoran returns a Komoran Korean tokenizer
func NewKomoran() Tokenizer {
	return &Komoran{}
}

// Komoran is a Korean Tokenizer in java
type Komoran struct {
	jvm *jnigi.JVM
	env *jnigi.Env
}

const jarPath = "tokenizers/build/komoran"

// Setup setups up the JVM for Komoran to run
func (k *Komoran) Setup() error {
	classPathArray, err := jarPaths()
	if err != nil {
		return err
	}

	k.jvm, k.env, err = createJVM(strings.Join(append(classPathArray, ""), ":"))
	if err != nil {
		return err
	}
	return nil
}

func createJVM(classPath string) (*jnigi.JVM, *jnigi.Env, error) {
	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		return nil, nil, err
	}

	args := []string{
		"-Xcheck:jni",
		"-Djava.class.path=" + classPath,
	}
	return jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, int(jnigi.DEFAULT_VERSION), args))
}

func jarPaths() ([]string, error) {
	files, err := ioutil.ReadDir(jarPath)
	if err != nil {
		return nil, err
	}

	a := make([]string, len(files))
	for i, f := range files {
		a[i] = filepath.Join(jarPath, f.Name())
	}
	return a, nil
}

// Cleanup cleans up the jvm
func (k *Komoran) Cleanup() error {
	if err := k.jvm.Destroy(); err != nil {
		return err
	}
	k.jvm, k.env = nil, nil
	return nil
}

// IsSetup returns true if the Komoran tokenizer is set up
func (k *Komoran) IsSetup() bool {
	return k.jvm != nil && k.env != nil
}

// GetTokens returns the grammar toekns of the given string
func (k *Komoran) GetTokens() (string, error) {
	if !k.IsSetup() {
		return "", &NotSetupError{}
	}
	return k.callStringMethod("text2anki/tokenizer/komoran/Tokenizer", "testy")
}

func (k *Komoran) callStringMethod(className string, methodName string, args ...interface{}) (string, error) {
	result, err := k.env.CallStaticMethod(className, methodName, jnigi.ObjectType("java/lang/String"), args...)
	if err != nil {
		return "", err
	}

	result, err = result.(*jnigi.ObjectRef).CallMethod(k.env, "getBytes", jnigi.Byte|jnigi.Array)
	if err != nil {
		return "", err
	}
	return string(result.([]byte)), nil
}
