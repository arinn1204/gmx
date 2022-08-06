package gmx

import (
	"gmx/internal/java"
	"log"
)

func DoWork() {
	log.Println("doin some work")
	java.CreateJvm()
}
