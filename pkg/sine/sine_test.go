package sine

import (
	"bytes"
	"flag"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ECecillo/lib.go.sound/pkg/format"
	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update-golden", false, "update golden test files")

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

// Golden file test helpers

// compareWithGoldenFile compares generated audio data with a golden file
func compareWithGoldenFile(t *testing.T, goldenFilePath string, data []byte) {
	t.Helper()

	if *updateGolden {
		// Update mode: write the new golden file
		err := os.MkdirAll(filepath.Dir(goldenFilePath), 0755)
		require.NoError(t, err, "failed to create testdata directory")

		err = os.WriteFile(goldenFilePath, data, 0644)
		require.NoError(t, err, "failed to write golden file")

		t.Logf("Updated golden file: %s", goldenFilePath)
		return
	}

	// Compare mode: read and compare with existing golden file
	goldenData, err := os.ReadFile(goldenFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("Golden file does not exist: %s\nRun with -update-golden to create it", goldenFilePath)
	}
	require.NoError(t, err, "failed to read golden file")

	if !bytes.Equal(data, goldenData) {
		t.Errorf("Generated audio does not match golden file: %s", goldenFilePath)
		t.Errorf("Generated size: %d bytes, Golden size: %d bytes", len(data), len(goldenData))

		// Find first difference
		minLen := min(len(goldenData), len(data))

		for i := range minLen {
			if data[i] != goldenData[i] {
				t.Errorf("First difference at byte %d: got 0x%02x, want 0x%02x", i, data[i], goldenData[i])
				break
			}
		}

		t.Fatal("Golden file mismatch detected")
	}
}

// TestGoldenFiles verifies audio generation produces consistent output
func TestGoldenFiles(t *testing.T) {
	tests := []struct {
		name      string
		frequency float64
		duration  time.Duration
		amplitude float64
		sampling  float64
		format    format.AudioFormat
		filename  string
	}{
		{
			name:      "440hz_1sec_pcm16",
			frequency: 440.0,
			duration:  time.Second,
			amplitude: 1.0,
			sampling:  44100.0,
			format:    format.PCM16{},
			filename:  "440hz_1sec_pcm16.bin",
		},
		{
			name:      "440hz_1sec_pcm32",
			frequency: 440.0,
			duration:  time.Second,
			amplitude: 1.0,
			sampling:  44100.0,
			format:    format.PCM32{},
			filename:  "440hz_1sec_pcm32.bin",
		},
		{
			name:      "440hz_1sec_float64",
			frequency: 440.0,
			duration:  time.Second,
			amplitude: 1.0,
			sampling:  44100.0,
			format:    format.Float64{},
			filename:  "440hz_1sec_float64.bin",
		},
		{
			name:      "1000hz_500ms_pcm16",
			frequency: 1000.0,
			duration:  500 * time.Millisecond,
			amplitude: 0.8,
			sampling:  44100.0,
			format:    format.PCM16{},
			filename:  "1000hz_500ms_pcm16.bin",
		},
		{
			name:      "220hz_100ms_low_amp",
			frequency: 220.0,
			duration:  100 * time.Millisecond,
			amplitude: 0.3,
			sampling:  44100.0,
			format:    format.PCM16{},
			filename:  "220hz_100ms_low_amp.bin",
		},
		{
			name:      "1hz_1sec_pcm16_low_sampling",
			frequency: 1.0,
			duration:  time.Second,
			amplitude: 1.0,
			sampling:  100.0,
			format:    format.PCM16{},
			filename:  "1hz_1sec_pcm16_low_sampling.bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate audio
			sine := NewSine(
				tt.frequency,
				tt.duration,
				WithAmplitude(tt.amplitude),
				WithSamplingRate(tt.sampling),
				WithFormat(tt.format),
			)

			// Write to buffer
			var buf bytes.Buffer
			_, err := sine.WriteTo(&buf)
			require.NoError(t, err, "failed to generate audio")

			// Compare with golden file
			goldenPath := filepath.Join("testdata", tt.filename)
			compareWithGoldenFile(t, goldenPath, buf.Bytes())
		})
	}
}

// TestGoldenFiles_EdgeCases tests edge cases with golden files
func TestGoldenFiles_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		frequency float64
		duration  time.Duration
		amplitude float64
		sampling  float64
		format    format.AudioFormat
		filename  string
	}{
		{
			name:      "zero_amplitude",
			frequency: 440.0,
			duration:  100 * time.Millisecond,
			amplitude: 0.0,
			sampling:  44100.0,
			format:    format.PCM16{},
			filename:  "zero_amplitude.bin",
		},
		{
			name:      "very_high_frequency",
			frequency: 10000.0,
			duration:  50 * time.Millisecond,
			amplitude: 0.5,
			sampling:  44100.0,
			format:    format.PCM16{},
			filename:  "very_high_frequency.bin",
		},
		{
			name:      "very_low_frequency",
			frequency: 10.0,
			duration:  time.Second,
			amplitude: 1.0,
			sampling:  44100.0,
			format:    format.PCM16{},
			filename:  "very_low_frequency.bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sine := NewSine(
				tt.frequency,
				tt.duration,
				WithAmplitude(tt.amplitude),
				WithSamplingRate(tt.sampling),
				WithFormat(tt.format),
			)

			var buf bytes.Buffer
			_, err := sine.WriteTo(&buf)
			require.NoError(t, err)

			goldenPath := filepath.Join("testdata", tt.filename)
			compareWithGoldenFile(t, goldenPath, buf.Bytes())
		})
	}
}

// TestGoldenFiles_Consistency ensures multiple generations produce identical output
func TestGoldenFiles_Consistency(t *testing.T) {
	sine := NewSine(440.0, 100*time.Millisecond, WithAmplitude(0.8))

	// Generate 3 times
	outputs := make([][]byte, 3)
	for i := range 3 {
		var buf bytes.Buffer
		_, err := sine.WriteTo(&buf)
		require.NoError(t, err)
		outputs[i] = buf.Bytes()
	}

	// All outputs should be identical
	require.Equal(t, outputs[0], outputs[1], "First and second generation differ")
	require.Equal(t, outputs[0], outputs[2], "First and third generation differ")
	require.Equal(t, outputs[1], outputs[2], "Second and third generation differ")
}
