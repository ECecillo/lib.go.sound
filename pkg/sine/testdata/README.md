# Golden Test Files

This directory contains reference audio files used for regression testing. These files are generated with known-good parameters and serve as the "golden standard" for verifying that audio generation remains consistent over time.

## What Are Golden Files?

Golden files (also called snapshot tests) are pre-generated reference outputs that tests compare against. If your code changes cause the output to differ from these files, the tests will fail, alerting you to potential regressions.

## Files in This Directory

| File | Description |
|------|-------------|
| `440hz_1sec_pcm16.bin` | 440Hz sine wave, 1 second, PCM16 format |
| `440hz_1sec_pcm32.bin` | 440Hz sine wave, 1 second, PCM32 format |
| `440hz_1sec_float64.bin` | 440Hz sine wave, 1 second, Float64 format |
| `1000hz_500ms_pcm16.bin` | 1000Hz sine wave, 500ms, PCM16 format |
| `220hz_100ms_low_amp.bin` | 220Hz sine wave, 100ms, low amplitude (0.3) |
| `1hz_1sec_pcm16_low_sampling.bin` | 1Hz sine wave, 1 second, low sampling rate (100Hz) |
| `zero_amplitude.bin` | Sine wave with zero amplitude (silence) |
| `very_high_frequency.bin` | 10000Hz sine wave (high frequency edge case) |
| `very_low_frequency.bin` | 10Hz sine wave (low frequency edge case) |

## Running the Tests

### Normal Test Run (Comparison Mode)
```bash
go test ./pkg/sine/... -run TestGoldenFiles
```

This compares newly generated audio against the golden files. Tests will fail if there are any differences.

### Updating Golden Files
```bash
go test ./pkg/sine/... -run TestGoldenFiles -update-golden
```

Use this flag when you've made **intentional** changes to the audio generation logic and need to update the reference files.

⚠️ **Warning:** Only update golden files when you're certain the changes are correct! Review the differences carefully before updating.

## When to Update Golden Files

Update golden files in these scenarios:
1. You've fixed a bug that improves audio quality
2. You've made intentional changes to the sine generation algorithm
3. You've modified the format conversion logic (e.g., improved scaling)

## When NOT to Update Golden Files

Don't update if:
1. Tests fail unexpectedly - investigate the root cause first
2. You're not sure why the output changed
3. The changes might be regressions

## Example: Detecting a Regression

If you accidentally introduce a bug:
```go
// Bug: adding an offset to the sine wave
value := s.Amplitude * math.Sin(angle) + 0.001
```

The golden file test will catch it:
```
--- FAIL: TestGoldenFiles/440hz_1sec_pcm16 (0.00s)
    sine_test.go:235: Generated audio does not match golden file
    sine_test.go:235: Generated size: 88200 bytes, Golden size: 88200 bytes
    sine_test.go:235: First difference at byte 0: got 0x20, want 0x00
```

The test pinpoints the exact byte where the output diverged, making debugging easier.

## Best Practices

1. **Run tests before committing** - Ensure no regressions
2. **Review diffs carefully** - When updating golden files, understand what changed
3. **Keep golden files small** - Use short durations to minimize file sizes
4. **Test multiple scenarios** - Cover different frequencies, formats, and edge cases
5. **Version control** - Commit golden files to git so everyone uses the same references
