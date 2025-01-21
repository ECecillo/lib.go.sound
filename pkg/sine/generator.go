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
	result := make([]float64, totalSamples)

	for n := range totalSamples {

		value := calculateSampleValue(n, s)
		result[n] = value

	}
	return result, nil
}

func calculateSampleValue(n int, s Sine) float64 {
	t := float64(n) / s.SamplingRate

	angle := 2 * math.Pi * s.Frequency * t
	value := s.Amplitude * math.Sin(angle)
	return value
}
