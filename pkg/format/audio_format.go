package format

import (
	"math"
)

type AudioFormat interface {
	// BitDepth return an integer representing the byte deph of the format.
	BitDepth() int
	// ConvertSample combine Quantize and Encode process to
	// return a sample value in byte using byte shifting.
	ConvertSample(float64) []byte
}

type PCM16 struct{}

func (f PCM16) BitDepth() int {
	return 16
}

// Clamp bound value to the given min and max so we encounter a value
// oveflow during our Quantization potentially creating unpredictable
// behavior.
func Clamp(value, minVal, maxVal float64) float64 {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

func (f PCM16) ConvertSample(sample float64) []byte {
	value := f.Quantize(sample)
	return f.Encode(value)
}

// Quantize scale the float64 sample to the full int16 range (-32768 to 32767)
func (f PCM16) Quantize(sample float64) int16 {
	sample = Clamp(sample, -1.0, 1.0)
	return int16(sample * 32767.0)
}

func (f PCM16) Encode(value int16) []byte {
	return []byte{byte(value & 0xFF), byte((value >> 8) & 0xFF)}
}

type PCM32 struct{}

func (f PCM32) BitDepth() int {
	return 32
}

func (f PCM32) ConvertSample(sample float64) []byte {
	value := f.Quantize(sample)
	return f.Encode(value)
}

// Quantize scale the float64 sample to the full int32 range (-2147483648 to 2147483647)
func (f PCM32) Quantize(sample float64) int32 {
	sample = Clamp(sample, -1.0, 1.0)
	return int32(sample * 2147483647.0)
}

func (f PCM32) Encode(value int32) []byte {
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
	value := f.Quantize(sample)
	return f.Encode(value)
}

// Quantize converts the float64 sample to IEEE 754 binary representation
func (f Float64) Quantize(sample float64) uint64 {
	return math.Float64bits(sample)
}

func (f Float64) Encode(value uint64) []byte {
	return []byte{
		byte(value & 0xFF), byte((value >> 8) & 0xFF),
		byte((value >> 16) & 0xFF), byte((value >> 24) & 0xFF),
		byte((value >> 32) & 0xFF), byte((value >> 40) & 0xFF),
		byte((value >> 48) & 0xFF), byte((value >> 56) & 0xFF),
	}
}

var (
	_ AudioFormat = new(PCM16)
	_ AudioFormat = new(PCM32)
	_ AudioFormat = new(Float64)
)
