JAVA			:= javac -h .

LD_LIBRARY_PATH := $(JAVA_HOME)/libexec/openjdk.jdk/Contents/Home/lib/server
CLASSPATH 		:= internal/java

JNIConnector.class:
	$(JAVA) java/JNIConnector.java
	javago --classfile=JNIConnector.class

.PHONY:

build: JNIConnector.class
	go build -o gmx ./cmd/main 

clean:
	rm JNIConnector.class java/JNIConnector.class.go gmx  2>/dev/null || exit 0