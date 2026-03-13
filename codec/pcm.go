package codec

import (
	"encoding/binary"
	"math"
)

// Float32ToBytes converts float32 audio samples to bytes (16-bit PCM)
func Float32ToBytes(samples []float32) []byte {
	// Convert float32 (-1.0 to 1.0) to int16
	int16Samples := make([]int16, len(samples))
	for i, s := range samples {
		// Clamp to [-1.0, 1.0]
		if s < -1.0 {
			s = -1.0
		} else if s > 1.0 {
			s = 1.0
		}
		int16Samples[i] = int16(s * 32767.0)
	}

	// Convert to bytes (little-endian)
	bytes := make([]byte, len(int16Samples)*2)
	for i, s := range int16Samples {
		binary.LittleEndian.PutUint16(bytes[i*2:], uint16(s))
	}

	return bytes
}

// BytesToFloat32 converts bytes (16-bit PCM) to float32 audio samples
func BytesToFloat32(data []byte) []float32 {
	// Must have even number of bytes
	if len(data)%2 != 0 {
		return nil
	}

	samples := make([]float32, len(data)/2)
	for i := 0; i < len(samples); i++ {
		int16Val := int16(binary.LittleEndian.Uint16(data[i*2:]))
		samples[i] = float32(int16Val) / 32767.0
	}

	return samples
}

// Interleave converts separate mono channels to interleaved stereo
// Takes two mono buffers (left and right) and interleaves them
func Interleave(left, right []float32) []float32 {
	if len(left) != len(right) {
		return nil
	}

	stereo := make([]float32, len(left)*2)
	for i := 0; i < len(left); i++ {
		stereo[i*2] = left[i]    // Left channel
		stereo[i*2+1] = right[i] // Right channel
	}

	return stereo
}

// Deinterleave converts interleaved stereo to separate mono channels
func Deinterleave(stereo []float32) (left, right []float32) {
	if len(stereo)%2 != 0 {
		return nil, nil
	}

	frames := len(stereo) / 2
	left = make([]float32, frames)
	right = make([]float32, frames)

	for i := 0; i < frames; i++ {
		left[i] = stereo[i*2]
		right[i] = stereo[i*2+1]
	}

	return left, right
}

// MixMonoToStereo duplicates mono audio to both left and right channels
func MixMonoToStereo(mono []float32) []float32 {
	stereo := make([]float32, len(mono)*2)
	for i := 0; i < len(mono); i++ {
		stereo[i*2] = mono[i]
		stereo[i*2+1] = mono[i]
	}
	return stereo
}

// MixStereoToMono averages left and right channels to mono
func MixStereoToMono(stereo []float32) []float32 {
	if len(stereo)%2 != 0 {
		return nil
	}

	mono := make([]float32, len(stereo)/2)
	for i := 0; i < len(mono); i++ {
		mono[i] = (stereo[i*2] + stereo[i*2+1]) / 2.0
	}
	return mono
}

// VolumeAdjust multiplies all samples by a volume factor
func VolumeAdjust(samples []float32, volume float32) []float32 {
	result := make([]float32, len(samples))
	for i, s := range samples {
		result[i] = s * volume
		// Clamp to [-1.0, 1.0]
		if result[i] > 1.0 {
			result[i] = 1.0
		} else if result[i] < -1.0 {
			result[i] = -1.0
		}
	}
	return result
}

// RMS calculates the RMS (root mean square) of audio samples
func RMS(samples []float32) float64 {
	if len(samples) == 0 {
		return 0.0
	}

	var sum float64
	for _, s := range samples {
		sum += float64(s) * float64(s)
	}

	return math.Sqrt(sum / float64(len(samples)))
}
