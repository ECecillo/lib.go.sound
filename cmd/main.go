package main

import (
	"os"
	"time"

	"github.com/ECecillo/lib.go.sound/pkg/format"
	"github.com/ECecillo/lib.go.sound/pkg/sine"
)

func main() {
	conf := sine.Config{
		Frequency:    440.0,
		Duration:     4 * time.Second,
		Amplitude:    32767.0,
		SamplingRate: 44100.0,
		Format:       format.PCM16{},
	}

	s := sine.NewSine(conf)

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
