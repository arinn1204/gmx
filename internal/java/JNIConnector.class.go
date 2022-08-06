
package java

import "github.com/juntaki/jnigo"

var jvm *jnigo.JVM

func init() {
	jvm = jnigo.CreateJVM()
}

type class struct{
	jclass *jnigo.JClass
}

func Newclass(args ...interface{}) (*class, error) {
	convertedArgs, err := jvm.ConvertAll(args)
	if err != nil{
		return nil, err
	}
	jclass, err := jvm.NewJClass("JNIConnector/class", convertedArgs)
	if err != nil{
		return nil, err
	}
	return &class{
		jclass: jclass,
	}, nil
}






