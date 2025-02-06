package format

import "math"

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

func (f PCM16) ConvertSample(sample float64) []byte {
	value := int16(sample)
	// Little-endian
	return []byte{byte(value & 0xFF), byte((value >> 8) & 0xFF)}
}

type PCM32 struct{}

func (f PCM32) BitDepth() int {
	return 32
}

func (f PCM32) ConvertSample(sample float64) []byte {
	value := int32(sample)
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
