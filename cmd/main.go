package main

import (
	"os"
	"time"

	"github.com/ECecillo/lib.go.sound/pkg/sine"
)

func main() {

	frequency := 440.0
	duration := 2 * time.Second

	s := sine.NewSine(frequency, duration)

	// Créer un fichier pour écrire les données audio
	file, err := os.Create("data/output.bin")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Écrire les données dans le fichier
	bytesWritten, err := s.WriteTo(file)
	if err != nil {
		panic(err)
	}

	println("Bytes written:", bytesWritten)
}
