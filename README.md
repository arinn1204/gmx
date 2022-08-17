# Go JMX Client (GMX)

## What it does  

GMX is built to interact with live JMX RMI servers. It has the ability to interact with a variety of different types of operations. This is done mostly in go using CGO. It will spin up a JVM in order to do the different JMX operations. The JVM will be created/destroyed during the life cycle of the GMX client.

The parameter order when using the clients are required to be in the same order as they are defined on the operation. GMX will query the operation in order to find the ordered parameters. The inputted types are only required when executing operations that require a generic type as a parameter.

## Limitations

Due to java's type erasure, GMX does not support setting nested lists or map attributes. It also does not support using operations that have parameters that require a nested list.

In this context, a nested list would be something like: `List<List<String>> nestedList = new ArrayList<>()`

## Examples

See the example's directory

## Using GMX

This will require a valid version of java on the machine. This was built and tested with java 18.

Since GMX relies on CGO, we will need a few environment variables defined to facilitate building and linking against `jni.h` and `libjvm.so`/`libjvm.dylib`/`libjvm.dll`

The Makefile is setup for a Mac with the includes/compiler flags. But to set it up you will need:

* `CGO_CFLAGS=-I$JAVA_HOME/include`

This will set the compiler flags. JNIGI (the JNI go bridge) will need access to `jni.h` in order to do any of the jni operations

* `-L$JAVA_HOME/lib -ljvm`

This will include the JVM shared library from your JDK. This is required for runtime ability to create a JVM and actually perform the JNI operations