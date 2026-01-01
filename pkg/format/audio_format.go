package format

import (
	"math"
)

type AudioFormat interface {
	BitDepth() int            // BitDepth return an integer representing the byte deph of the format.
	Quantize(float64) int     // Quantize converts a float64 sample to an integer representation.
	Encode(int) []byte        // Encode converts the integer value to bytes using byte shifting.
}

type PCM16 struct{}

func (f PCM16) BitDepth() int {
	return 16
}

// NOTE: We could have use std lib function for these functions but its way cooler to understand
// how it works under the hood.

// Clamp make sure that we never exceed any type limit
// which could lead to a value oveflow and generate unpredictable
// behavior in the sound.
func Clamp(value, minVal, maxVal float64) float64 {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

func (f PCM16) Quantize(sample float64) int {
	sample = Clamp(sample, -1.0, 1.0)
	// Scale the float64 sample (usually btw -1.0 and +1.0)
	// to the full int16 range (-32768 to 32767)
	//
	// Example: if we have 0.85 then int16(0.85) = 0 which is false.
	//
	// Why 2^(16) - 1 ? to avoid int overflow
	return int(sample * 32767.0)
}

func (f PCM16) Encode(value int) []byte {
	// Convert to int16 and encode as little-endian, cutting 16-bits to 2-bytes.
	val16 := int16(value)
	return []byte{byte(val16 & 0xFF), byte((val16 >> 8) & 0xFF)}
}

type PCM32 struct{}

func (f PCM32) BitDepth() int {
	return 32
}

func (f PCM32) Quantize(sample float64) int {
	sample = Clamp(sample, -1.0, 1.0)
	// Scale the float64 sample to the full int32 range (-2147483648 to 2147483647)
	return int(sample * 2147483647.0)
}

func (f PCM32) Encode(value int) []byte {
	// Convert to int32 and encode as little-endian bytes
	val32 := int32(value)
	return []byte{
		byte(val32 & 0xFF), byte((val32 >> 8) & 0xFF),
		byte((val32 >> 16) & 0xFF), byte((val32 >> 24) & 0xFF),
	}
}

type Float64 struct{}

func (f Float64) BitDepth() int {
	return 64
}

func (f Float64) Quantize(sample float64) int {
	// For Float64, we convert to IEEE 754 bit representation
	return int(math.Float64bits(sample))
}

func (f Float64) Encode(value int) []byte {
	// Convert back to uint64 for byte encoding (IEEE 754)
	val64 := uint64(value)
	return []byte{
		byte(val64 & 0xFF), byte((val64 >> 8) & 0xFF),
		byte((val64 >> 16) & 0xFF), byte((val64 >> 24) & 0xFF),
		byte((val64 >> 32) & 0xFF), byte((val64 >> 40) & 0xFF),
		byte((val64 >> 48) & 0xFF), byte((val64 >> 56) & 0xFF),
	}
}
