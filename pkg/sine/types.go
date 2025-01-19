package sine

import (
	"time"

	"github.com/ECecillo/lib.go.sound/pkg/format"
)

// TODO: remove once we add options in constructor.
type Config struct {
	Frequency    float64       // Fréquence en Hz
	Duration     time.Duration // Durée en secondes
	Amplitude    float64       // Amplitude (optionnelle, par défaut 1.0)
	SamplingRate float64       // Fréquence d'échantillonnage en Hz
	Format       format.AudioFormat
}

type Sine struct {
	Frequency    float64       // Fréquence en Hz
	Duration     time.Duration // Durée en secondes
	Amplitude    float64       // Amplitude (optionnelle, par défaut 1.0)
	SamplingRate float64       // Fréquence d'échantillonnage en Hz
	Format       format.AudioFormat
}

// TODO: Option to configure config value or set default one
// TODO: remove conf once we added options
func NewSine(conf Config) *Sine {
	return &Sine{
		Frequency:    conf.Frequency,
		Duration:     conf.Duration,
		Amplitude:    conf.Amplitude,
		SamplingRate: conf.SamplingRate,
		Format:       conf.Format,
	}
}
