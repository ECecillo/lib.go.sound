package format

import "math"

type AudioFormat interface {
	BitDepth() int                // Retourne la profondeur de bits (16, 32, 64)
	IsFloatingPoint() bool        // Indique si le format est en virgule flottante
	ConvertSample(float64) []byte // Convertit une valeur d'échantillon en bytes pour ce format
}

type PCM16 struct{}

func (f PCM16) BitDepth() int {
	return 16
}

func (f PCM16) IsFloatingPoint() bool {
	return false
}

func (f PCM16) ConvertSample(sample float64) []byte {
	value := int16(sample)                                       // Convertir en int16
	return []byte{byte(value & 0xFF), byte((value >> 8) & 0xFF)} // Little-endian
}

type PCM32 struct{}

func (f PCM32) BitDepth() int {
	return 32
}

func (f PCM32) IsFloatingPoint() bool {
	return false
}

func (f PCM32) ConvertSample(sample float64) []byte {
	value := int32(sample) // Convertir en int32
	return []byte{
		byte(value & 0xFF), byte((value >> 8) & 0xFF),
		byte((value >> 16) & 0xFF), byte((value >> 24) & 0xFF),
	}
}

type Float64 struct{}

func (f Float64) BitDepth() int {
	return 64
}

func (f Float64) IsFloatingPoint() bool {
	return true
}

func (f Float64) ConvertSample(sample float64) []byte {
	// Convertir en float64 (IEEE 754)
	value := math.Float64bits(sample)
	return []byte{
		byte(value & 0xFF), byte((value >> 8) & 0xFF),
		byte((value >> 16) & 0xFF), byte((value >> 24) & 0xFF),
		byte((value >> 32) & 0xFF), byte((value >> 40) & 0xFF),
		byte((value >> 48) & 0xFF), byte((value >> 56) & 0xFF),
	}
}
