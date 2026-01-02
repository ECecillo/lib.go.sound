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

func TestNyquistFiltering(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		frequency      float64
		samplingRate   float64
		expectFiltered bool
	}{
		{
			name:           "below_nyquist_limit",
			frequency:      1000.0,
			samplingRate:   44100.0,
			expectFiltered: false,
			description:    "Frequency well below Nyquist limit should pass through",
		},
		{
			name:           "at_nyquist_limit",
			frequency:      22050.0,
			samplingRate:   44100.0,
			expectFiltered: true,
			description:    "Frequency at exactly Nyquist limit should be filtered",
		},
		{
			name:           "above_nyquist_limit",
			frequency:      25000.0,
			samplingRate:   44100.0,
			expectFiltered: true,
			description:    "Frequency above Nyquist limit should be filtered",
		},
		{
			name:           "far_above_nyquist_limit",
			frequency:      50000.0,
			samplingRate:   44100.0,
			expectFiltered: true,
			description:    "Frequency far above Nyquist limit should be filtered",
		},
		{
			name:           "just_below_nyquist_limit",
			frequency:      22049.0,
			samplingRate:   44100.0,
			expectFiltered: false,
			description:    "Frequency just below Nyquist limit should pass through",
		},
		{
			name:           "low_sampling_rate_above_nyquist",
			frequency:      1000.0,
			samplingRate:   1500.0,
			expectFiltered: true,
			description:    "1000 Hz with 1500 Hz sampling (Nyquist = 750 Hz) should be filtered",
		},
		{
			name:           "low_sampling_rate_below_nyquist",
			frequency:      600.0,
			samplingRate:   1500.0,
			expectFiltered: false,
			description:    "600 Hz with 1500 Hz sampling (Nyquist = 750 Hz) should pass through",
		},
		{
			name:           "low_sampling_rate_at_nyquist",
			frequency:      500.0,
			samplingRate:   1000.0,
			expectFiltered: true,
			description:    "500 Hz at Nyquist limit with 1000 Hz sampling should be filtered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sine := NewSine(
				tt.frequency,
				100*time.Millisecond,
				WithSamplingRate(tt.samplingRate),
				WithAmplitude(1.0),
			)

			samples, err := sine.Generate()
			require.NoError(t, err)
			require.NotEmpty(t, samples, "Should generate samples")

			if tt.expectFiltered {
				// All samples should be zero when filtered
				for i, sample := range samples {
					require.Equal(t, 0.0, sample,
						"Sample %d should be 0.0 (filtered) but got %f. %s",
						i, sample, tt.description)
				}
			} else {
				// At least some samples should be non-zero when not filtered
				hasNonZero := false
				for _, sample := range samples {
					if sample != 0.0 {
						hasNonZero = true
						break
					}
				}
				require.True(t, hasNonZero,
					"Signal should have non-zero samples (not filtered). %s",
					tt.description)
			}
		})
	}
}

func TestAntiAliasingFilter(t *testing.T) {
	tests := []struct {
		name         string
		frequency    float64
		samplingRate float64
		signal       float64
		expected     float64
	}{
		{
			name:         "pass_through_below_nyquist",
			frequency:    1000.0,
			samplingRate: 44100.0,
			signal:       0.5,
			expected:     0.5,
		},
		{
			name:         "filter_at_nyquist",
			frequency:    22050.0,
			samplingRate: 44100.0,
			signal:       0.8,
			expected:     0.0,
		},
		{
			name:         "filter_above_nyquist",
			frequency:    30000.0,
			samplingRate: 44100.0,
			signal:       1.0,
			expected:     0.0,
		},
		{
			name:         "pass_through_zero_signal",
			frequency:    440.0,
			samplingRate: 44100.0,
			signal:       0.0,
			expected:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sine := NewSine(
				tt.frequency,
				time.Second,
				WithSamplingRate(tt.samplingRate),
			)

			result := sine.applyAntiAliasingFilter(tt.signal)
			require.Equal(t, tt.expected, result,
				"Filter output mismatch for frequency=%f, samplingRate=%f",
				tt.frequency, tt.samplingRate)
		})
	}
}

func TestContinuousSignalAt(t *testing.T) {
	tests := []struct {
		name      string
		frequency float64
		amplitude float64
		timePoint float64
	}{
		{
			name:      "440hz_at_zero",
			frequency: 440.0,
			amplitude: 1.0,
			timePoint: 0.0,
		},
		{
			name:      "440hz_at_quarter_period",
			frequency: 440.0,
			amplitude: 1.0,
			timePoint: 1.0 / (4.0 * 440.0), // Quarter period should give max amplitude
		},
		{
			name:      "1000hz_various_amplitude",
			frequency: 1000.0,
			amplitude: 0.5,
			timePoint: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sine := NewSine(
				tt.frequency,
				time.Second,
				WithAmplitude(tt.amplitude),
			)

			result := sine.continuousSignalAt(tt.timePoint)

			// Verify the result is within amplitude bounds
			require.LessOrEqual(t, math.Abs(result), tt.amplitude,
				"Signal amplitude exceeds maximum at t=%f", tt.timePoint)

			// Verify it matches the expected sine calculation
			expectedAngle := 2 * math.Pi * tt.frequency * tt.timePoint
			expected := tt.amplitude * math.Sin(expectedAngle)
			require.Equal(t, expected, result,
				"Continuous signal value mismatch at t=%f", tt.timePoint)
		})
	}
}

// Golden file test helpers

// compareWithGoldenFile compares generated audio data with a golden file
func compareWithGoldenFile(t *testing.T, goldenFilePath string, data []byte) {
	t.Helper()

	if *updateGolden {
		// Update mode: write the new golden file
		err := os.MkdirAll(filepath.Dir(goldenFilePath), 0o755)
		require.NoError(t, err, "failed to create testdata directory")

		err = os.WriteFile(goldenFilePath, data, 0o644)
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
		format    format.AudioFormat
		name      string
		filename  string
		duration  time.Duration
		frequency float64
		amplitude float64
		sampling  float64
	}{
		{
			format:    format.PCM16{},
			duration:  time.Second,
			name:      "440hz_1sec_pcm16",
			filename:  "440hz_1sec_pcm16.bin",
			frequency: 440.0,
			amplitude: 1.0,
			sampling:  44100.0,
		},
		{
			format:    format.PCM32{},
			duration:  time.Second,
			name:      "440hz_1sec_pcm32",
			filename:  "440hz_1sec_pcm32.bin",
			frequency: 440.0,
			amplitude: 1.0,
			sampling:  44100.0,
		},
		{
			format:    format.PCM16{},
			duration:  500 * time.Millisecond,
			name:      "1000hz_500ms_pcm16",
			filename:  "1000hz_500ms_pcm16.bin",
			frequency: 1000.0,
			amplitude: 0.8,
			sampling:  44100.0,
		},
		{
			format:    format.PCM16{},
			duration:  100 * time.Millisecond,
			name:      "220hz_100ms_low_amp",
			filename:  "220hz_100ms_low_amp.bin",
			frequency: 220.0,
			amplitude: 0.3,
			sampling:  44100.0,
		},
		{
			format:    format.PCM16{},
			duration:  time.Second,
			name:      "1hz_1sec_pcm16_low_sampling",
			filename:  "1hz_1sec_pcm16_low_sampling.bin",
			frequency: 1.0,
			amplitude: 1.0,
			sampling:  100.0,
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
		format    format.AudioFormat
		name      string
		filename  string
		duration  time.Duration
		frequency float64
		amplitude float64
		sampling  float64
	}{
		{
			format:    format.PCM16{},
			duration:  100 * time.Millisecond,
			name:      "zero_amplitude",
			filename:  "zero_amplitude.bin",
			frequency: 440.0,
			amplitude: 0.0,
			sampling:  44100.0,
		},
		{
			format:    format.PCM16{},
			duration:  50 * time.Millisecond,
			name:      "very_high_frequency",
			filename:  "very_high_frequency.bin",
			frequency: 10000.0,
			amplitude: 0.5,
			sampling:  44100.0,
		},
		{
			format:    format.PCM16{},
			duration:  time.Second,
			name:      "very_low_frequency",
			filename:  "very_low_frequency.bin",
			frequency: 10.0,
			amplitude: 1.0,
			sampling:  44100.0,
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
