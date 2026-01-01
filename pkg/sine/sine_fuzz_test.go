package sine

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/ECecillo/lib.go.sound/pkg/format"
)

// FuzzSineGeneration tests sine generation with random parameters
func FuzzSineGeneration(f *testing.F) {
	// Seed corpus with typical and edge case values
	// Format: frequency, durationMs, amplitude, samplingRate
	f.Add(440.0, int64(1000), 1.0, 44100.0)
	f.Add(1.0, int64(1000), 0.5, 44100.0)
	f.Add(1000.0, int64(100), 0.8, 48000.0)
	f.Add(20.0, int64(50), 0.3, 8000.0)
	f.Add(20000.0, int64(10), 1.0, 96000.0)

	f.Fuzz(func(t *testing.T, frequency float64, durationMs int64, amplitude, samplingRate float64) {
		// Skip invalid inputs
		if math.IsNaN(frequency) || math.IsInf(frequency, 0) || frequency <= 0 {
			t.Skip("invalid frequency")
		}
		if math.IsNaN(amplitude) || math.IsInf(amplitude, 0) || amplitude < 0 {
			t.Skip("invalid amplitude")
		}
		if math.IsNaN(samplingRate) || math.IsInf(samplingRate, 0) || samplingRate <= 0 {
			t.Skip("invalid sampling rate")
		}
		if durationMs <= 0 || durationMs > 10000 { // Cap at 10 seconds for performance
			t.Skip("invalid or too long duration")
		}

		// Limit extreme values for performance
		if frequency > 1e6 {
			t.Skip("frequency too high")
		}
		if samplingRate > 1e6 {
			t.Skip("sampling rate too high")
		}
		if amplitude > 1000 {
			amplitude = 1000 // Cap amplitude
		}

		// Skip cases that violate Nyquist-Shannon sampling theorem
		// Frequency must be less than half the sampling rate to avoid aliasing
		if frequency >= samplingRate/2 {
			t.Skip("frequency violates Nyquist criterion")
		}

		duration := time.Duration(durationMs) * time.Millisecond

		// Create sine wave - should never panic
		sine := NewSine(frequency, duration, WithAmplitude(amplitude), WithSamplingRate(samplingRate))

		// Generate samples - should never panic
		samples, err := sine.Generate()
		if err != nil {
			t.Fatalf("Generate() failed: %v", err)
		}

		// Verify sample count
		expectedSamples := int(samplingRate * duration.Seconds())
		if len(samples) != expectedSamples {
			t.Errorf("Expected %d samples, got %d", expectedSamples, len(samples))
		}

		// Verify all samples are within amplitude bounds
		for i, sample := range samples {
			if math.IsNaN(sample) {
				t.Errorf("Sample %d is NaN", i)
			}
			if math.IsInf(sample, 0) {
				t.Errorf("Sample %d is Inf", i)
			}
			if math.Abs(sample) > amplitude {
				t.Errorf("Sample %d has amplitude %f, exceeds limit %f", i, math.Abs(sample), amplitude)
			}
		}

		// Verify samples are not all zero (unless amplitude is zero)
		// Note: With very short durations or specific frequency/sampling combinations,
		// we might sample at zero-crossing points, resulting in near-zero values.
		// Only check this if we have good sampling density.
		if amplitude > 0 && len(samples) >= 10 {
			period := 1.0 / frequency                          // Period in seconds
			periodsInDuration := duration.Seconds() / period   // Number of periods
			samplesPerPeriod := samplingRate / frequency       // Samples per period

			// Only check for non-zero values if:
			// 1. We have at least 0.5 periods (enough to reach a peak or trough)
			// 2. We have at least 10 samples per period (good sampling density)
			// This avoids false positives from undersampling near Nyquist limit
			if periodsInDuration >= 0.5 && samplesPerPeriod >= 10 {
				hasSignificantValue := false
				threshold := amplitude * 0.01 // 1% of amplitude
				for _, sample := range samples {
					if math.Abs(sample) > threshold {
						hasSignificantValue = true
						break
					}
				}
				if !hasSignificantValue {
					t.Errorf("All samples are near-zero despite non-zero amplitude %f, %d samples, %.2f periods, %.2f samples/period",
						amplitude, len(samples), periodsInDuration, samplesPerPeriod)
				}
			}
		}
	})
}

// FuzzSineWriteTo tests the full pipeline with random parameters
func FuzzSineWriteTo(f *testing.F) {
	// Seed corpus
	f.Add(440.0, int64(100), 1.0, 44100.0, uint8(0)) // PCM16
	f.Add(440.0, int64(100), 1.0, 44100.0, uint8(1)) // PCM32
	f.Add(440.0, int64(100), 1.0, 44100.0, uint8(2)) // Float64

	f.Fuzz(func(t *testing.T, frequency float64, durationMs int64, amplitude, samplingRate float64, formatType uint8) {
		// Skip invalid inputs
		if math.IsNaN(frequency) || math.IsInf(frequency, 0) || frequency <= 0 || frequency > 1e6 {
			t.Skip()
		}
		if math.IsNaN(amplitude) || math.IsInf(amplitude, 0) || amplitude < 0 || amplitude > 1000 {
			t.Skip()
		}
		if math.IsNaN(samplingRate) || math.IsInf(samplingRate, 0) || samplingRate <= 0 || samplingRate > 1e6 {
			t.Skip()
		}
		if durationMs <= 0 || durationMs > 1000 { // Cap at 1 second for WriteTo performance
			t.Skip()
		}
		// Skip cases that violate Nyquist criterion
		if frequency >= samplingRate/2 {
			t.Skip()
		}

		duration := time.Duration(durationMs) * time.Millisecond

		// Skip if duration is too short to produce at least one sample
		if int(samplingRate*duration.Seconds()) < 1 {
			t.Skip()
		}

		// Select format based on formatType
		var audioFormat format.AudioFormat
		var bytesPerSample int
		switch formatType % 3 {
		case 0:
			audioFormat = format.PCM16{}
			bytesPerSample = 2
		case 1:
			audioFormat = format.PCM32{}
			bytesPerSample = 4
		case 2:
			audioFormat = format.Float64{}
			bytesPerSample = 8
		}

		// Create and write - should never panic
		sine := NewSine(
			frequency,
			duration,
			WithAmplitude(amplitude),
			WithSamplingRate(samplingRate),
			WithFormat(audioFormat),
		)

		var buf bytes.Buffer
		bytesWritten, err := sine.WriteTo(&buf)
		if err != nil {
			t.Fatalf("WriteTo() failed: %v", err)
		}

		// Verify bytes written
		if bytesWritten <= 0 {
			t.Errorf("No bytes written")
		}

		expectedSamples := int(samplingRate * duration.Seconds())
		expectedBytes := int64(expectedSamples * bytesPerSample)

		if bytesWritten != expectedBytes {
			t.Errorf("Expected %d bytes (%d samples * %d bytes/sample), got %d",
				expectedBytes, expectedSamples, bytesPerSample, bytesWritten)
		}

		// Verify buffer length matches
		if int64(buf.Len()) != bytesWritten {
			t.Errorf("Buffer length %d doesn't match bytesWritten %d", buf.Len(), bytesWritten)
		}
	})
}

// FuzzCalculateSampleValue tests the core sample calculation with random inputs
func FuzzCalculateSampleValue(f *testing.F) {
	// Seed corpus
	f.Add(440.0, 1.0, 44100.0, 0)
	f.Add(440.0, 1.0, 44100.0, 100)
	f.Add(440.0, 1.0, 44100.0, 44099)
	f.Add(1.0, 0.5, 100.0, 0)
	f.Add(1000.0, 0.8, 48000.0, 1000)

	f.Fuzz(func(t *testing.T, frequency, amplitude, samplingRate float64, sampleIndex int) {
		// Skip invalid inputs
		if math.IsNaN(frequency) || math.IsInf(frequency, 0) || frequency <= 0 {
			t.Skip()
		}
		if math.IsNaN(amplitude) || math.IsInf(amplitude, 0) || amplitude < 0 {
			t.Skip()
		}
		if math.IsNaN(samplingRate) || math.IsInf(samplingRate, 0) || samplingRate <= 0 {
			t.Skip()
		}
		if sampleIndex < 0 || sampleIndex > 1000000 {
			t.Skip()
		}

		sine := Sine{
			Frequency:    frequency,
			Amplitude:    amplitude,
			SamplingRate: samplingRate,
		}

		// Should never panic
		value := sine.calculateSampleValue(sampleIndex)

		// Verify output properties
		if math.IsNaN(value) {
			t.Errorf("calculateSampleValue returned NaN")
		}
		if math.IsInf(value, 0) {
			t.Errorf("calculateSampleValue returned Inf")
		}

		// Value should be within amplitude bounds
		if math.Abs(value) > amplitude {
			t.Errorf("Sample value %f exceeds amplitude %f", value, amplitude)
		}
	})
}

// FuzzReproducibility ensures that generating with same parameters produces same results
func FuzzReproducibility(f *testing.F) {
	f.Add(440.0, int64(100), 1.0, 44100.0)
	f.Add(1000.0, int64(50), 0.5, 48000.0)

	f.Fuzz(func(t *testing.T, frequency float64, durationMs int64, amplitude, samplingRate float64) {
		// Skip invalid inputs
		if math.IsNaN(frequency) || math.IsInf(frequency, 0) || frequency <= 0 || frequency > 1e6 {
			t.Skip()
		}
		if math.IsNaN(amplitude) || math.IsInf(amplitude, 0) || amplitude < 0 || amplitude > 1000 {
			t.Skip()
		}
		if math.IsNaN(samplingRate) || math.IsInf(samplingRate, 0) || samplingRate <= 0 || samplingRate > 1e6 {
			t.Skip()
		}
		if durationMs <= 0 || durationMs > 1000 {
			t.Skip()
		}
		// Skip cases that violate Nyquist criterion
		if frequency >= samplingRate/2 {
			t.Skip()
		}

		duration := time.Duration(durationMs) * time.Millisecond

		// Generate twice with same parameters
		sine1 := NewSine(frequency, duration, WithAmplitude(amplitude), WithSamplingRate(samplingRate))
		sine2 := NewSine(frequency, duration, WithAmplitude(amplitude), WithSamplingRate(samplingRate))

		samples1, err1 := sine1.Generate()
		if err1 != nil {
			t.Fatalf("First generation failed: %v", err1)
		}

		samples2, err2 := sine2.Generate()
		if err2 != nil {
			t.Fatalf("Second generation failed: %v", err2)
		}

		// Should produce identical results
		if len(samples1) != len(samples2) {
			t.Fatalf("Sample counts differ: %d vs %d", len(samples1), len(samples2))
		}

		for i := range samples1 {
			if samples1[i] != samples2[i] {
				t.Errorf("Sample %d differs: %f vs %f", i, samples1[i], samples2[i])
				break // Report first difference only
			}
		}
	})
}
