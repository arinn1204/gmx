package main

import (
	"fmt"
	"gmx/internal/jvm"
	"gmx/internal/mbean"
	"log"
)

func main() {

	javaVM, err := jvm.CreateJVM()

	if err != nil {
		log.Panicf("failed to start jvm::%s", err.Error())
	}

	defer javaVM.ShutdownJvm()

	beanExecutor, err := jvm.CreateMBeanConnection(javaVM, "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")

	if err != nil {
		log.Panicf("failed to initialize the connection::%s", err.Error())
	}

	operation := mbean.Operation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "putString",
		Args: []mbean.OperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
			{
				Value: "fan369",
				Type:  "java.lang.String",
			},
		},
	}
	beanExecutor.Execute(operation)

	operation = mbean.Operation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "getString",
		Args: []mbean.OperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
		},
	}
	res, err := beanExecutor.Execute(operation)

	fmt.Println(res)
	fmt.Println(err)
}
