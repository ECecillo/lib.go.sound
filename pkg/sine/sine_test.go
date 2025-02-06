package sine

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAmplitude(t *testing.T) {
	sine := NewSine(440.0, time.Second, WithAmplitude(1.0))
	samples, err := sine.Generate()
	require.NoError(t, err)

	for _, v := range samples {
		require.LessOrEqual(t, math.Abs(v), sine.Amplitude, "Amplitude exceeds limit")
	}
}

func TestGenerateConstantPeriodSampleCount(t *testing.T) {
	// NOTE: Important value here because if we choose something else we
	// would be using a float value that would be round at some point
	// and since we only want to
	frequency := 1.0
	samplingRate := 44100.0

	periodDuration := time.Duration(1 / frequency)

	// Création d'une instance de Sine avec la durée fixée à une période.
	sineWave := NewSine(frequency, periodDuration, WithSamplingRate(samplingRate))

	samples, err := sineWave.Generate()
	require.NoError(t, err)
	totalSamples := len(samples)

	totalExpectedSamples := int(samplingRate * periodDuration.Seconds())

	require.Equal(t, totalExpectedSamples, totalSamples, "incorrect number of samples for one period, expected %d but got %d", totalExpectedSamples, totalSamples)
}

func TestReproducibility(t *testing.T) {
	sine1 := NewSine(440.0, time.Second)
	sine2 := NewSine(440.0, time.Second)

	samples1, err := sine1.Generate()
	require.NoError(t, err)

	samples2, err := sine2.Generate()
	require.NoError(t, err)

	require.Equal(t, samples1, samples2, "Signals are not reproducible")
}

func TestCorrectness(t *testing.T) {
	frequency := 1.0
	duration := time.Second
	samplingRate := 10.0
	amplitude := 1.0

	sine := NewSine(frequency, duration, WithSamplingRate(samplingRate), WithAmplitude(amplitude))
	samples, err := sine.Generate()
	require.NoError(t, err)

	for i, sample := range samples {
		expected := sine.calculateSampleValue(i)

		require.Equal(t, expected, sample, "Mismatch at sample %d", sample, expected)
	}
}

func TestWriteTo(t *testing.T) {
	sine := NewSine(440.0, time.Second)
	buffer := &bytes.Buffer{}

	bytesWritten, err := sine.WriteTo(buffer)
	require.NoError(t, err)
	require.Greater(t, bytesWritten, int64(0), "No data written to writer")
	require.Equal(t, bytesWritten, int64(buffer.Len()), "Mismatch between bytes written and buffer length")
}

func TestExtremeParameters(t *testing.T) {
	sine := NewSine(440.0, time.Second, WithAmplitude(0.0))
	samples, err := sine.Generate()
	require.NoError(t, err)

	for _, v := range samples {
		require.Equal(t, 0.0, v, "Signal is not zero with zero amplitude")
	}

	sine = NewSine(1e6, time.Second, WithSamplingRate(1e7)) // 1 MHz with 10 MHz sampling rate
	samples, err = sine.Generate()
	require.NoError(t, err)
	require.Len(t, samples, int(sine.SamplingRate*sine.Duration.Seconds()), "Unexpected number of samples")
}
