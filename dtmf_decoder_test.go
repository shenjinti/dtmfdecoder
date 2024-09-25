package dtmfdecoder

import (
	"os"
	"testing"
)

func TestGoertzel(t *testing.T) {
	g := NewGoertzel([]int{697, 770, 852, 941}, 44100)
	g.processSample(42)
	g.processSample(84)
	energies := g.energies
	if energies[697] != 0.8980292970055115 {
		t.Error("Expected", 0.8980292970055115, "got", energies[697])
	}
	if energies[770] != 0.8975953139667142 {
		t.Error("Expected", 0.8975953139667142, "got", energies[770])
	}
	if energies[852] != 0.8970565383230518 {
		t.Error("Expected", 0.8970565383230518, "got", energies[852])
	}
	if energies[941] != 0.8964104403348228 {
		t.Error("Expected", 0.8964104403348228, "got", energies[941])
	}
}
func TestDecodeDTMF(t *testing.T) {
	energyThreshold := 0.032 // ???
	t.Run("TestDecodeDTMF 8k", func(t *testing.T) {
		decoder := NewDTMFDecoder(energyThreshold, 8000)
		f, _ := os.Open("testdata/123456654321_8k_s16le.raw")
		decoded := ""
		for {
			data := make([]byte, 320)
			n, err := f.Read(data)
			if n == 0 || err != nil {
				break
			}
			samples := make([]float64, 0)
			for i := 0; i < n; i += 2 {
				sample := int16(data[i]) | int16(data[i+1])<<8
				samples = append(samples, float64(sample)/32768.0)
			}

			current, ok := decoder.Decode(samples)
			if ok {
				decoded += current
			}
		}
		t.Log("[123456654321_8k_s16le] Detected DTMF tone:", decoded)
		if decoded != "123456654321" {
			t.Error("Expected 123456654321, got", decoded)
		}
	})
	t.Run("TestDecodeDTMF mix 16k", func(t *testing.T) {
		decoder := NewDTMFDecoder(energyThreshold, 16000)
		f, _ := os.Open("testdata/mix_1_16k_s16le.raw")
		decoded := ""

		for {
			data := make([]byte, 640)
			n, err := f.Read(data)
			if n == 0 || err != nil {
				break
			}
			samples := make([]float64, 0)
			for i := 0; i < n; i += 2 {
				sample := int16(data[i]) | int16(data[i+1])<<8
				samples = append(samples, float64(sample)/32768.0)
			}
			current, ok := decoder.Decode(samples)
			if ok {
				decoded += current
			}
		}
		t.Log("[mix_1_16k_s16le] Detected DTMF tone:", decoded)
		if decoded != "1" {
			t.Error("Expected 1, got", decoded)
		}
	})
}
