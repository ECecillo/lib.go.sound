package format

import (
	"math"
	"testing"
)

// FuzzPCM16_ConvertSample tests PCM16 conversion with random float64 values
func FuzzPCM16_ConvertSample(f *testing.F) {
	// Seed corpus with interesting values
	f.Add(0.0)
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(0.5)
	f.Add(-0.5)
	f.Add(2.0)   // Out of range
	f.Add(-2.0)  // Out of range
	f.Add(math.MaxFloat64)
	f.Add(-math.MaxFloat64)

	format := PCM16{}

	f.Fuzz(func(t *testing.T, sample float64) {
		// Skip NaN and Inf as they're not valid audio samples
		if math.IsNaN(sample) || math.IsInf(sample, 0) {
			t.Skip()
		}

		// Should never panic
		result := format.ConvertSample(sample)

		// Should always produce exactly 2 bytes
		if len(result) != 2 {
			t.Errorf("Expected 2 bytes, got %d", len(result))
		}

		// Reconstruct the value to verify it's within int16 range
		value := int16(result[0]) | int16(result[1])<<8

		// Value should be within int16 range
		if value < math.MinInt16 || value > math.MaxInt16 {
			t.Errorf("Reconstructed value %d out of int16 range", value)
		}

		// If input was in [-1, 1], verify output scaling is reasonable
		if sample >= -1.0 && sample <= 1.0 {
			expectedMax := int16(32767)
			if value < -expectedMax || value > expectedMax {
				t.Errorf("For input %f in range [-1,1], got value %d outside expected range", sample, value)
			}
		}
	})
}

// FuzzPCM32_ConvertSample tests PCM32 conversion with random float64 values
func FuzzPCM32_ConvertSample(f *testing.F) {
	// Seed corpus
	f.Add(0.0)
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(0.5)
	f.Add(-0.5)
	f.Add(10.0)
	f.Add(-10.0)

	format := PCM32{}

	f.Fuzz(func(t *testing.T, sample float64) {
		if math.IsNaN(sample) || math.IsInf(sample, 0) {
			t.Skip()
		}

		// Should never panic
		result := format.ConvertSample(sample)

		// Should always produce exactly 4 bytes
		if len(result) != 4 {
			t.Errorf("Expected 4 bytes, got %d", len(result))
		}

		// Reconstruct to verify it's within int32 range
		value := int32(result[0]) | int32(result[1])<<8 | int32(result[2])<<16 | int32(result[3])<<24

		// Value should be within int32 range (this is implicit, but we check for sanity)
		if value < math.MinInt32 || value > math.MaxInt32 {
			t.Errorf("Reconstructed value %d out of int32 range", value)
		}
	})
}

// FuzzFloat64_ConvertSample tests Float64 conversion with random values
func FuzzFloat64_ConvertSample(f *testing.F) {
	// Seed corpus with diverse values
	f.Add(0.0)
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(math.Pi)
	f.Add(math.E)
	f.Add(0.123456789)
	f.Add(-0.987654321)
	f.Add(1e10)
	f.Add(1e-10)

	format := Float64{}

	f.Fuzz(func(t *testing.T, sample float64) {
		// Float64 should handle all values, including NaN and Inf
		result := format.ConvertSample(sample)

		// Should always produce exactly 8 bytes
		if len(result) != 8 {
			t.Errorf("Expected 8 bytes, got %d", len(result))
		}

		// Reconstruct and verify round-trip
		bits := uint64(result[0]) | uint64(result[1])<<8 | uint64(result[2])<<16 | uint64(result[3])<<24 |
			uint64(result[4])<<32 | uint64(result[5])<<40 | uint64(result[6])<<48 | uint64(result[7])<<56
		reconstructed := math.Float64frombits(bits)

		// For normal values, should round-trip exactly
		if math.IsNaN(sample) {
			if !math.IsNaN(reconstructed) {
				t.Errorf("NaN did not round-trip correctly")
			}
		} else {
			if sample != reconstructed {
				t.Errorf("Round-trip failed: input=%v, output=%v", sample, reconstructed)
			}
		}
	})
}

// FuzzClamp tests the Clamp utility function with random values
func FuzzClamp(f *testing.F) {
	// Seed corpus
	f.Add(0.0, -1.0, 1.0)
	f.Add(2.0, -1.0, 1.0)
	f.Add(-2.0, -1.0, 1.0)
	f.Add(0.5, 0.0, 1.0)
	f.Add(100.0, -10.0, 10.0)

	f.Fuzz(func(t *testing.T, value, min, max float64) {
		// Skip invalid ranges or special values
		if math.IsNaN(value) || math.IsNaN(min) || math.IsNaN(max) ||
			math.IsInf(value, 0) || math.IsInf(min, 0) || math.IsInf(max, 0) {
			t.Skip()
		}

		if min > max {
			t.Skip() // Invalid range
		}

		result := Clamp(value, min, max)

		// Result should always be within [min, max]
		if result < min {
			t.Errorf("Clamp(%f, %f, %f) = %f, which is less than min", value, min, max, result)
		}
		if result > max {
			t.Errorf("Clamp(%f, %f, %f) = %f, which is greater than max", value, min, max, result)
		}

		// Verify expected behavior
		if value < min && result != min {
			t.Errorf("Expected %f to be clamped to min %f, got %f", value, min, result)
		}
		if value > max && result != max {
			t.Errorf("Expected %f to be clamped to max %f, got %f", value, max, result)
		}
		if value >= min && value <= max && result != value {
			t.Errorf("Expected %f to remain unchanged, got %f", value, result)
		}
	})
}

// FuzzAllFormats_ConsistentBehavior ensures all formats handle the same input consistently
func FuzzAllFormats_ConsistentBehavior(f *testing.F) {
	f.Add(0.0)
	f.Add(0.5)
	f.Add(-0.5)
	f.Add(1.0)
	f.Add(-1.0)

	formats := []struct {
		name         string
		format       AudioFormat
		expectedSize int
	}{
		{"PCM16", PCM16{}, 2},
		{"PCM32", PCM32{}, 4},
		{"Float64", Float64{}, 8},
	}

	f.Fuzz(func(t *testing.T, sample float64) {
		if math.IsNaN(sample) || math.IsInf(sample, 0) {
			t.Skip()
		}

		for _, ft := range formats {
			// Should never panic
			result := ft.format.ConvertSample(sample)

			// Should produce correct byte count
			if len(result) != ft.expectedSize {
				t.Errorf("%s: expected %d bytes, got %d", ft.name, ft.expectedSize, len(result))
			}

			// Should match bit depth
			expectedBytes := ft.format.BitDepth() / 8
			if len(result) != expectedBytes {
				t.Errorf("%s: bit depth %d implies %d bytes, but got %d",
					ft.name, ft.format.BitDepth(), expectedBytes, len(result))
			}
		}
	})
}
