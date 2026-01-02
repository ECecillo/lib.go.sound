package format

import (
	"math"
	"testing"
)

// BenchmarkPCM16_ConvertSample benchmarks PCM16 format conversion
func BenchmarkPCM16_ConvertSample(b *testing.B) {
	format := PCM16{}
	sample := 0.75

	for b.Loop() {
		_ = format.ConvertSample(sample)
	}
}

// BenchmarkPCM16_ConvertSample_Clamping benchmarks with out-of-range values
func BenchmarkPCM16_ConvertSample_Clamping(b *testing.B) {
	format := PCM16{}
	sample := 2.5 // Out of range, will be clamped

	for b.Loop() {
		_ = format.ConvertSample(sample)
	}
}

// BenchmarkPCM32_ConvertSample benchmarks PCM32 format conversion
func BenchmarkPCM32_ConvertSample(b *testing.B) {
	format := PCM32{}
	sample := 0.75

	for b.Loop() {
		_ = format.ConvertSample(sample)
	}
}

// BenchmarkFloat64_ConvertSample benchmarks Float64 format conversion
func BenchmarkFloat64_ConvertSample(b *testing.B) {
	format := Float64{}
	sample := 0.75

	for b.Loop() {
		_ = format.ConvertSample(sample)
	}
}

// BenchmarkClamp benchmarks the Clamp utility function
func BenchmarkClamp(b *testing.B) {
	b.Run("InRange", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Clamp(0.5, -1.0, 1.0)
		}
	})

	b.Run("BelowMin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Clamp(-2.0, -1.0, 1.0)
		}
	})

	b.Run("AboveMax", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Clamp(2.0, -1.0, 1.0)
		}
	})
}

// BenchmarkAllFormats_ConvertSample compares all formats side-by-side
func BenchmarkAllFormats_ConvertSample(b *testing.B) {
	sample := 0.75

	b.Run("PCM16", func(b *testing.B) {
		format := PCM16{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = format.ConvertSample(sample)
		}
	})

	b.Run("PCM32", func(b *testing.B) {
		format := PCM32{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = format.ConvertSample(sample)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		format := Float64{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = format.ConvertSample(sample)
		}
	})
}

// BenchmarkBatchConversion benchmarks converting many samples
func BenchmarkBatchConversion(b *testing.B) {
	const numSamples = 1000
	samples := make([]float64, numSamples)
	for i := range samples {
		samples[i] = math.Sin(float64(i) * 0.1)
	}

	b.Run("PCM16", func(b *testing.B) {
		format := PCM16{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, sample := range samples {
				_ = format.ConvertSample(sample)
			}
		}
		b.ReportMetric(float64(numSamples*b.N)/b.Elapsed().Seconds(), "samples/sec")
	})

	b.Run("PCM32", func(b *testing.B) {
		format := PCM32{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, sample := range samples {
				_ = format.ConvertSample(sample)
			}
		}
		b.ReportMetric(float64(numSamples*b.N)/b.Elapsed().Seconds(), "samples/sec")
	})

	b.Run("Float64", func(b *testing.B) {
		format := Float64{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, sample := range samples {
				_ = format.ConvertSample(sample)
			}
		}
		b.ReportMetric(float64(numSamples*b.N)/b.Elapsed().Seconds(), "samples/sec")
	})
}
