package sine

import (
	"bytes"
	"testing"
	"time"

	"github.com/ECecillo/lib.go.sound/pkg/format"
)

// BenchmarkCalculateSampleValue benchmarks the core sine calculation
func BenchmarkCalculateSampleValue(b *testing.B) {
	sine := NewSine(440.0, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sine.calculateSampleValue(i % 44100)
	}
}

// BenchmarkGenerate benchmarks full sine wave generation
func BenchmarkGenerate(b *testing.B) {
	b.Run("100ms_44100Hz", func(b *testing.B) {
		sine := NewSine(440.0, 100*time.Millisecond, WithSamplingRate(44100.0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = sine.Generate()
		}
		// Report samples per second
		totalSamples := int(44100 * 0.1 * float64(b.N))
		b.ReportMetric(float64(totalSamples)/b.Elapsed().Seconds(), "samples/sec")
	})

	b.Run("1sec_44100Hz", func(b *testing.B) {
		sine := NewSine(440.0, time.Second, WithSamplingRate(44100.0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = sine.Generate()
		}
		totalSamples := int(44100 * float64(b.N))
		b.ReportMetric(float64(totalSamples)/b.Elapsed().Seconds(), "samples/sec")
	})

	b.Run("1sec_48000Hz", func(b *testing.B) {
		sine := NewSine(440.0, time.Second, WithSamplingRate(48000.0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = sine.Generate()
		}
		totalSamples := int(48000 * float64(b.N))
		b.ReportMetric(float64(totalSamples)/b.Elapsed().Seconds(), "samples/sec")
	})

	b.Run("10sec_44100Hz", func(b *testing.B) {
		sine := NewSine(440.0, 10*time.Second, WithSamplingRate(44100.0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = sine.Generate()
		}
		totalSamples := int(44100 * 10 * float64(b.N))
		b.ReportMetric(float64(totalSamples)/b.Elapsed().Seconds(), "samples/sec")
	})
}

// BenchmarkGenerate_DifferentFrequencies benchmarks different frequencies
func BenchmarkGenerate_DifferentFrequencies(b *testing.B) {
	frequencies := []struct {
		name string
		freq float64
	}{
		{"LowFreq_20Hz", 20.0},
		{"MidFreq_440Hz", 440.0},
		{"MidFreq_1000Hz", 1000.0},
		{"HighFreq_10000Hz", 10000.0},
	}

	for _, tc := range frequencies {
		b.Run(tc.name, func(b *testing.B) {
			sine := NewSine(tc.freq, time.Second, WithSamplingRate(44100.0))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = sine.Generate()
			}
		})
	}
}

// BenchmarkWriteTo benchmarks the full pipeline including format conversion
func BenchmarkWriteTo(b *testing.B) {
	b.Run("PCM16_1sec", func(b *testing.B) {
		sine := NewSine(440.0, time.Second, WithFormat(format.PCM16{}))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
		totalSamples := int(44100 * float64(b.N))
		b.ReportMetric(float64(totalSamples)/b.Elapsed().Seconds(), "samples/sec")
		b.ReportMetric(float64(totalSamples*2)/b.Elapsed().Seconds(), "bytes/sec")
	})

	b.Run("PCM32_1sec", func(b *testing.B) {
		sine := NewSine(440.0, time.Second, WithFormat(format.PCM32{}))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
		totalSamples := int(44100 * float64(b.N))
		b.ReportMetric(float64(totalSamples)/b.Elapsed().Seconds(), "samples/sec")
		b.ReportMetric(float64(totalSamples*4)/b.Elapsed().Seconds(), "bytes/sec")
	})

	b.Run("Float64_1sec", func(b *testing.B) {
		sine := NewSine(440.0, time.Second, WithFormat(format.Float64{}))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
		totalSamples := int(44100 * float64(b.N))
		b.ReportMetric(float64(totalSamples)/b.Elapsed().Seconds(), "samples/sec")
		b.ReportMetric(float64(totalSamples*8)/b.Elapsed().Seconds(), "bytes/sec")
	})
}

// BenchmarkWriteTo_DifferentDurations benchmarks different audio durations
func BenchmarkWriteTo_DifferentDurations(b *testing.B) {
	durations := []struct {
		name     string
		duration time.Duration
	}{
		{"10ms", 10 * time.Millisecond},
		{"100ms", 100 * time.Millisecond},
		{"500ms", 500 * time.Millisecond},
		{"1sec", time.Second},
		{"5sec", 5 * time.Second},
	}

	for _, tc := range durations {
		b.Run(tc.name, func(b *testing.B) {
			sine := NewSine(440.0, tc.duration, WithFormat(format.PCM16{}))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var buf bytes.Buffer
				_, _ = sine.WriteTo(&buf)
			}
		})
	}
}

// BenchmarkWriteTo_DifferentAmplitudes benchmarks different amplitudes
func BenchmarkWriteTo_DifferentAmplitudes(b *testing.B) {
	amplitudes := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	for _, amp := range amplitudes {
		b.Run("", func(b *testing.B) {
			sine := NewSine(440.0, 100*time.Millisecond, WithAmplitude(amp))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var buf bytes.Buffer
				_, _ = sine.WriteTo(&buf)
			}
		})
	}
}

// BenchmarkFullPipeline benchmarks realistic audio generation scenarios
func BenchmarkFullPipeline(b *testing.B) {
	b.Run("BeepSound_440Hz_100ms_PCM16", func(b *testing.B) {
		// Typical short beep sound
		sine := NewSine(440.0, 100*time.Millisecond,
			WithAmplitude(0.8),
			WithSamplingRate(44100.0),
			WithFormat(format.PCM16{}),
		)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
	})

	b.Run("ToneTest_1000Hz_1sec_PCM32", func(b *testing.B) {
		// Standard test tone
		sine := NewSine(1000.0, time.Second,
			WithAmplitude(1.0),
			WithSamplingRate(48000.0),
			WithFormat(format.PCM32{}),
		)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
	})

	b.Run("LowFreqDrone_60Hz_5sec_Float64", func(b *testing.B) {
		// Low frequency ambient sound
		sine := NewSine(60.0, 5*time.Second,
			WithAmplitude(0.5),
			WithSamplingRate(44100.0),
			WithFormat(format.Float64{}),
		)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
	})
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("Generate_1sec", func(b *testing.B) {
		sine := NewSine(440.0, time.Second)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = sine.Generate()
		}
	})

	b.Run("WriteTo_1sec", func(b *testing.B) {
		sine := NewSine(440.0, time.Second)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
	})
}

// BenchmarkParallelGeneration benchmarks concurrent generation
func BenchmarkParallelGeneration(b *testing.B) {
	sine := NewSine(440.0, 100*time.Millisecond)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = sine.Generate()
		}
	})
}

// BenchmarkParallelWriteTo benchmarks concurrent writing
func BenchmarkParallelWriteTo(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sine := NewSine(440.0, 100*time.Millisecond)
		for pb.Next() {
			var buf bytes.Buffer
			_, _ = sine.WriteTo(&buf)
		}
	})
}
