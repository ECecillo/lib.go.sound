package sine

import (
	"fmt"
	"io"
	"math"
)

// WriteTo will generate samples and write them to the given Writer.
func (s Sine) WriteTo(w io.Writer) (int64, error) {
	// Générer les données de l'onde sinus
	samples, err := s.Generate()
	if err != nil {
		return 0, fmt.Errorf("unable to generate samples, err: %w", err)
	}

	// Will help us count the number of bytes written.
	var totalBytesWritten int64

	// Write each sampel to a little endian byte encoding (int16 -> 2 octets)
	for i := range len(samples) {
		data := s.Format.ConvertSample(s.Amplitude * samples[i])

		// Convert sample to little-endian
		n, err := w.Write(data)
		if err != nil {
			return totalBytesWritten, fmt.Errorf("unable to write data, err: %w", err)
		}
		totalBytesWritten += int64(n)

	}

	return totalBytesWritten, nil
}

// FIXME: we need to find a way to also create a function that return []int32
func (s Sine) Generate() ([]float64, error) {

	totalSamples := int(s.SamplingRate * s.Duration.Seconds())
	result := make([]float64, totalSamples)

	for n := range totalSamples {

		t := float64(n) / s.SamplingRate

		angle := 2 * math.Pi * s.Frequency * t
		value := math.Sin(angle)

		result[n] = value

	}
	return result, nil
}
