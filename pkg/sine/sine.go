package sine

import (
	"fmt"
	"io"
	"math"
)

// WriteTo will generate samples and write them to the given Writer.
func (s Sine) WriteTo(w io.Writer) (int64, error) {
	samples, err := s.Generate()
	if err != nil {
		return 0, fmt.Errorf("unable to generate samples, err: %w", err)
	}

	// Will help us count the number of bytes written.
	var totalBytesWritten int64

	for i := range len(samples) {
		data := s.Format.ConvertSample(samples[i])

		n, err := w.Write(data)
		if err != nil {
			return totalBytesWritten, fmt.Errorf("unable to write data, err: %w", err)
		}
		totalBytesWritten += int64(n)

	}

	return totalBytesWritten, nil
}

func (s Sine) Generate() ([]float64, error) {
	totalSamples := int(s.SamplingRate * s.Duration.Seconds())
	result := make([]float64, 0, totalSamples)

	for n := range totalSamples {
		value := s.calculateSampleValue(n)
		result = append(result, value)
	}
	return result, nil
}

// continuousSignalAt simulates the continuous sine wave signal at time t.
// This represents the physical sound wave before any electronic processing.
func (s Sine) continuousSignalAt(t float64) float64 {
	angle := 2 * math.Pi * s.Frequency * t
	return s.Amplitude * math.Sin(angle)
}

// applyAntiAliasingFilter simulates an analog anti-aliasing filter.
// If the frequency exceeds the Nyquist limit (SamplingRate/2), the filter
// cuts off the signal completely to prevent aliasing artifacts.
func (s Sine) applyAntiAliasingFilter(signal float64) float64 {
	nyquistLimit := s.SamplingRate / 2.0

	if s.Frequency >= nyquistLimit {
		return 0.0
	}

	return signal
}

// calculateSampleValue orchestrates the signal processing pipeline:
// 1. Generate the continuous signal at time t
// 2. Apply anti-aliasing filter
// 3. Return the filtered sample value
func (s Sine) calculateSampleValue(sampleIndex int) float64 {
	// Calculate time for this sample
	t := float64(sampleIndex) / s.SamplingRate

	// Step 1: Get the continuous signal value
	signal := s.continuousSignalAt(t)

	// Step 2: Apply anti-aliasing filter
	filteredSignal := s.applyAntiAliasingFilter(signal)

	return filteredSignal
}
