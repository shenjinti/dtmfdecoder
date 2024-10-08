# A pure Golang implementation of the DTMF detector
> Ported from https://dtmf.netlify.app

## Usage
```shell
go get github.com/shenjinti/dtmfdecoder
```
## API
```go
type DTMFDecoder struct {
    // The minimum duration of a DTMF tone in seconds
    PressInterval time.Duration // default 200ms
}
func NewDTMFDecoder(pressInterval float64, sampleRate float64) *DTMFDecoder
```

### Example 

```go
decoder := dtmfdecoder.NewDTMFDecoder(0.32, 8000)
f, _ := os.Open("testdata/123456654321_8k_s16le.raw")
decoded := ""
for {
    data := make([]byte, 320) // 320 bytes = 160 samples = 20ms
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
```

## Thanks to the original author for the code.
- https://github.com/Ravenstine/goertzeljs
- https://dtmf.netlify.app
