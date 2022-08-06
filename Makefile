JAVA			:= javac -h .

LD_LIBRARY_PATH := $(JAVA_HOME)/libexec/openjdk.jdk/Contents/Home/lib/server
CLASSPATH 		:= internal/java

JNIConnector.class:
	$(JAVA) internal/java/JNIConnector.java
	javago --classpath=internal.java --classfile=JNIConnector.class
	@mv java/* internal/java/
	@rmdir java

.PHONY:

build: JNIConnector.class
	go build ./cmd/main -o gmx

clean:
	rm internal/java/JNIConnector.class gmx