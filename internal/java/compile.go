package java

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

// Compile will compile the given file using the JVM compiler (javac)
func Compile(env *jnigi.Env, fileName string) error {
	dest := jnigi.NewObjectRef("javax/tools/JavaCompiler")
	if err := getCompiler(env, dest); err != nil {
		return errors.New("failed to create the java compiler::" + err.Error())
	}

	responseType := jnigi.NewObjectRef("I")
	in := jnigi.NewObjectRef("java/io/InputStream")
	out := jnigi.NewObjectRef("java/io/OutputStream")

	arrayRef, err := toStringArray(env, fileName)

	if err != nil {
		return fmt.Errorf("failed to turn %s into a parameter array::%s", fileName, err.Error())
	}

	env.PrecalculateSignature("(Ljava/io/InputStream;Ljava/io/OutputStream;Ljava/io/OutputStream;[Ljava/lang/String;)I")

	if err := dest.CallMethod(env, "run", responseType, in, out, out, arrayRef); err != nil {
		return fmt.Errorf("failed to compile %s::%s", fileName, err.Error())
	}

	return nil
}

func getCompiler(env *jnigi.Env, dest *jnigi.ObjectRef) error {
	staticClassReference := "javax/tools/ToolProvider"

	if err := env.CallStaticMethod(staticClassReference, "getSystemJavaCompiler", dest); err != nil {
		return errors.New("Failed to create the compiler::" + err.Error())
	}

	return nil
}

func toStringArray(env *jnigi.Env, str string) (*jnigi.ObjectRef, error) {
	fileNameRef, err := env.NewObject("java/lang/String", []byte(str))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", str, err.Error())
	}

	inputParams := []*jnigi.ObjectRef{fileNameRef}

	arrayRef := env.ToObjectArray(inputParams, "java/lang/String")

	return arrayRef, nil
}
