package format

import (
	"math"
)

type AudioFormat interface {
	BitDepth() int                // BitDepth return an integer representing the byte deph of the format.
	ConvertSample(float64) []byte // ConvertSample return the sample value in byte using byte shifting.
}

type PCM16 struct{}

func (f PCM16) BitDepth() int {
	return 16
}

// NOTE: We could have use std lib function for these functions but its way cooler to understand
// how it works under the hood.

// Clamp make sure that we never exceed any type limit
// which could lead to a value oveflow and generate unpredictable
// behaviour in the sound.
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (f PCM16) ConvertSample(sample float64) []byte {
	sample = Clamp(sample, -1.0, 1.0)
	// Scale the float64 sample (usually btw -1.0 and +1.0)
	// to the full int16 range (-32768 to 32767)
	//
	// Example: if we have 0.85 then int16(0.85) = 0 which is false.
	//
	// Why 2^(16) - 1 ? to avoid int overflow
	value := int16(sample * 32767.0)
	// Little-endian, cutting 16-bits to 2-bytes.
	return []byte{byte(value & 0xFF), byte((value >> 8) & 0xFF)}
}

type PCM32 struct{}

func (f PCM32) BitDepth() int {
	return 32
}

func (f PCM32) ConvertSample(sample float64) []byte {
	sample = Clamp(sample, -1.0, 1.0)
	// Scale the float64 sample to the full int32 range (-2147483648 to 2147483647)
	value := int32(sample * 2147483647.0)
	return []byte{
		byte(value & 0xFF), byte((value >> 8) & 0xFF),
		byte((value >> 16) & 0xFF), byte((value >> 24) & 0xFF),
	}
}

type Float64 struct{}

func (f Float64) BitDepth() int {
	return 64
}

func (f Float64) ConvertSample(sample float64) []byte {
	// (IEEE 754)
	value := math.Float64bits(sample)
	return []byte{
		byte(value & 0xFF), byte((value >> 8) & 0xFF),
		byte((value >> 16) & 0xFF), byte((value >> 24) & 0xFF),
		byte((value >> 32) & 0xFF), byte((value >> 40) & 0xFF),
		byte((value >> 48) & 0xFF), byte((value >> 56) & 0xFF),
	}
}
