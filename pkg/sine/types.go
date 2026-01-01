package sine

import (
	"time"

	"github.com/ECecillo/lib.go.sound/pkg/format"
)

type Sine struct {
	Format       format.AudioFormat
	Duration     time.Duration // Duration of the signal
	Frequency    float64       // Frequency in Hz
	Amplitude    float64       // Amplitude (optional, default 1.0)
	SamplingRate float64       // Sampling frequency in Hz
}

type Option func(*Sine)

func NewSine(frequency float64, duration time.Duration, options ...Option) *Sine {
	sine := &Sine{
		Frequency:    frequency,
		Duration:     duration,
		Amplitude:    1.0,
		SamplingRate: 44100.0,
		Format:       format.PCM16{},
	}

	for _, opt := range options {
		opt(sine)
	}

	return sine
}

func WithAmplitude(amplitude float64) Option {
	return func(s *Sine) {
		s.Amplitude = amplitude
	}
}

func WithSamplingRate(rate float64) Option {
	return func(s *Sine) {
		s.SamplingRate = rate
	}
}

func WithFormat(fmt format.AudioFormat) Option {
	return func(s *Sine) {
		s.Format = fmt
	}
}
