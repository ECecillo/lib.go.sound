package format

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClamp tests the Clamp utility function
func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{
			name:     "value within range",
			value:    0.5,
			min:      -1.0,
			max:      1.0,
			expected: 0.5,
		},
		{
			name:     "value below minimum",
			value:    -2.0,
			min:      -1.0,
			max:      1.0,
			expected: -1.0,
		},
		{
			name:     "value above maximum",
			value:    2.0,
			min:      -1.0,
			max:      1.0,
			expected: 1.0,
		},
		{
			name:     "value equals minimum",
			value:    -1.0,
			min:      -1.0,
			max:      1.0,
			expected: -1.0,
		},
		{
			name:     "value equals maximum",
			value:    1.0,
			min:      -1.0,
			max:      1.0,
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clamp(tt.value, tt.min, tt.max)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestPCM16_BitDepth verifies PCM16 reports correct bit depth
func TestPCM16_BitDepth(t *testing.T) {
	format := PCM16{}
	require.Equal(t, 16, format.BitDepth())
}

// TestPCM16_ConvertSample tests PCM16 sample conversion
func TestPCM16_ConvertSample(t *testing.T) {
	format := PCM16{}

	tests := []struct {
		name     string
		expected []byte
		input    float64
	}{
		{
			name:     "zero value",
			expected: []byte{0x00, 0x00},
			input:    0.0,
		},
		{
			name:     "maximum positive value",
			expected: []byte{0xFF, 0x7F}, // 32767 in little-endian
			input:    1.0,
		},
		{
			name:     "maximum negative value",
			input:    -1.0,
			expected: []byte{0x01, 0x80}, // -32767 in little-endian
		},
		{
			name:     "positive mid value",
			input:    0.5,
			expected: []byte{0xFF, 0x3F}, // 16383 in little-endian
		},
		{
			name:     "negative mid value",
			input:    -0.5,
			expected: []byte{0x01, 0xC0}, // -16383 in little-endian
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := format.Encode(format.Quantize(tt.input))
			require.Equal(t, tt.expected, result, "bytes mismatch for input %f", tt.input)
			require.Len(t, result, 2, "PCM16 should produce 2 bytes")
		})
	}
}

// TestPCM16_Clamping verifies that PCM16 clamps out-of-range values
func TestPCM16_Clamping(t *testing.T) {
	format := PCM16{}

	tests := []struct {
		name     string
		expected []byte
		input    float64
	}{
		{
			name:     "value above 1.0 clamped",
			expected: []byte{0xFF, 0x7F}, // should be clamped to 1.0 → 32767
			input:    2.5,
		},
		{
			name:     "value below -1.0 clamped",
			expected: []byte{0x01, 0x80}, // should be clamped to -1.0 → -32767
			input:    -2.5,
		},
		{
			name:     "slightly above 1.0",
			input:    1.1,
			expected: []byte{0xFF, 0x7F},
		},
		{
			name:     "slightly below -1.0",
			input:    -1.1,
			expected: []byte{0x01, 0x80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := format.Encode(format.Quantize(tt.input))
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestPCM16_LittleEndian verifies byte ordering is correct
func TestPCM16_LittleEndian(t *testing.T) {
	format := PCM16{}

	// Convert a known value and verify byte order
	// 0.5 * 32767 = 16383 = 0x3FFF
	// Little-endian: low byte first (0xFF), high byte second (0x3F)
	result := format.Encode(format.Quantize(0.5))
	require.Equal(t, byte(0xFF), result[0], "low byte should be first (little-endian)")
	require.Equal(t, byte(0x3F), result[1], "high byte should be second (little-endian)")

	// Reconstruct the int16 value to verify
	reconstructed := int16(result[0]) | int16(result[1])<<8
	require.Equal(t, int16(16383), reconstructed)
}

// TestPCM32_BitDepth verifies PCM32 reports correct bit depth
func TestPCM32_BitDepth(t *testing.T) {
	format := PCM32{}
	require.Equal(t, 32, format.BitDepth())
}

// TestPCM32_ConvertSample tests PCM32 sample conversion
func TestPCM32_ConvertSample(t *testing.T) {
	format := PCM32{}

	tests := []struct {
		name     string
		expected []byte
		input    float64
	}{
		{
			name:     "zero value",
			input:    0.0,
			expected: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "maximum positive value",
			input:    1.0,
			expected: []byte{0xFF, 0xFF, 0xFF, 0x7F}, // 2147483647 in little-endian
		},
		{
			name:     "maximum negative value",
			input:    -1.0,
			expected: []byte{0x01, 0x00, 0x00, 0x80}, // -2147483647 in little-endian
		},
		{
			name:     "positive mid value",
			input:    0.5,
			expected: []byte{0xFF, 0xFF, 0xFF, 0x3F}, // ~1073741823 in little-endian
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := format.Encode(format.Quantize(tt.input))
			require.Equal(t, tt.expected, result, "bytes mismatch for input %f", tt.input)
			require.Len(t, result, 4, "PCM32 should produce 4 bytes")
		})
	}
}

// TestPCM32_LittleEndian verifies byte ordering is correct
func TestPCM32_LittleEndian(t *testing.T) {
	format := PCM32{}

	// Convert a known value: 0.5 * 2147483647 = 1073741823 = 0x3FFFFFFF
	// Little-endian: 0xFF 0xFF 0xFF 0x3F
	result := format.Encode(format.Quantize(0.5))
	require.Equal(t, byte(0xFF), result[0], "byte 0 should be 0xFF")
	require.Equal(t, byte(0xFF), result[1], "byte 1 should be 0xFF")
	require.Equal(t, byte(0xFF), result[2], "byte 2 should be 0xFF")
	require.Equal(t, byte(0x3F), result[3], "byte 3 should be 0x3F")

	// Reconstruct the int32 value to verify
	reconstructed := int32(result[0]) | int32(result[1])<<8 | int32(result[2])<<16 | int32(result[3])<<24
	require.Equal(t, int32(1073741823), reconstructed)
}

// TestFloat64_BitDepth verifies Float64 reports correct bit depth
func TestFloat64_BitDepth(t *testing.T) {
	format := Float64{}
	require.Equal(t, 64, format.BitDepth())
}

// TestFloat64_ConvertSample tests Float64 sample conversion
func TestFloat64_ConvertSample(t *testing.T) {
	format := Float64{}

	tests := []struct {
		name  string
		input float64
	}{
		{name: "zero", input: 0.0},
		{name: "one", input: 1.0},
		{name: "negative one", input: -1.0},
		{name: "pi", input: math.Pi},
		{name: "e", input: math.E},
		{name: "small positive", input: 0.123456789},
		{name: "small negative", input: -0.987654321},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := format.Encode(format.Quantize(tt.input))
			require.Len(t, result, 8, "Float64 should produce 8 bytes")

			// Verify we can reconstruct the original value
			bits := uint64(result[0]) | uint64(result[1])<<8 | uint64(result[2])<<16 | uint64(result[3])<<24 |
				uint64(result[4])<<32 | uint64(result[5])<<40 | uint64(result[6])<<48 | uint64(result[7])<<56
			reconstructed := math.Float64frombits(bits)
			require.Equal(t, tt.input, reconstructed, "round-trip conversion failed")
		})
	}
}

// TestFloat64_RoundTrip verifies Float64 can round-trip any float64 value
func TestFloat64_RoundTrip(t *testing.T) {
	format := Float64{}

	testValues := []float64{
		0.0,
		1.0,
		-1.0,
		0.5,
		-0.5,
		math.MaxFloat64,
		math.SmallestNonzeroFloat64,
		-math.MaxFloat64,
		math.Pi,
		math.E,
		1e-10,
		1e10,
	}

	for _, value := range testValues {
		t.Run("", func(t *testing.T) {
			bytes := format.Encode(format.Quantize(value))

			// Reconstruct from bytes
			bits := uint64(bytes[0]) | uint64(bytes[1])<<8 | uint64(bytes[2])<<16 | uint64(bytes[3])<<24 |
				uint64(bytes[4])<<32 | uint64(bytes[5])<<40 | uint64(bytes[6])<<48 | uint64(bytes[7])<<56
			reconstructed := math.Float64frombits(bits)

			require.Equal(t, value, reconstructed, "failed to round-trip value %v", value)
		})
	}
}

// TestFloat64_SpecialValues tests Float64 with special IEEE 754 values
func TestFloat64_SpecialValues(t *testing.T) {
	format := Float64{}

	tests := []struct {
		name  string
		input float64
	}{
		{name: "positive infinity", input: math.Inf(1)},
		{name: "negative infinity", input: math.Inf(-1)},
		{name: "NaN", input: math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := format.Encode(format.Quantize(tt.input))
			require.Len(t, bytes, 8)

			// Reconstruct
			bits := uint64(bytes[0]) | uint64(bytes[1])<<8 | uint64(bytes[2])<<16 | uint64(bytes[3])<<24 |
				uint64(bytes[4])<<32 | uint64(bytes[5])<<40 | uint64(bytes[6])<<48 | uint64(bytes[7])<<56
			reconstructed := math.Float64frombits(bits)

			// Special handling for NaN (NaN != NaN)
			if math.IsNaN(tt.input) {
				require.True(t, math.IsNaN(reconstructed), "expected NaN")
			} else {
				require.Equal(t, tt.input, reconstructed)
			}
		})
	}
}

// TestAllFormats_ConsistentBehavior ensures all formats handle common cases consistently
func TestAllFormats_ConsistentBehavior(t *testing.T) {
	formats := []struct {
		format AudioFormat
		name   string
	}{
		{PCM16{}, "PCM16"},
		{PCM32{}, "PCM32"},
		{Float64{}, "Float64"},
	}

	for _, f := range formats {
		t.Run(f.name, func(t *testing.T) {
			// All formats should produce bytes
			result := f.format.Encode(f.format.Quantize(0.5))
			require.NotNil(t, result)
			require.Greater(t, len(result), 0)

			// Bit depth should match byte length
			expectedBytes := f.format.BitDepth() / 8
			require.Equal(t, expectedBytes, len(result),
				"bit depth %d should produce %d bytes", f.format.BitDepth(), expectedBytes)
		})
	}
}
