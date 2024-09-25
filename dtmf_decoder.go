package dtmfdecoder

import (
	"math"
	"time"
)

var frequencyTable = map[int]map[int]rune{
	697: {1209: '1', 1336: '2', 1477: '3', 1633: 'A'},
	770: {1209: '4', 1336: '5', 1477: '6', 1633: 'B'},
	852: {1209: '7', 1336: '8', 1477: '9', 1633: 'C'},
	941: {1209: '*', 1336: '0', 477: '#', 1633: 'D'},
}

type DTMFDecoder struct {
	energyThreshold float64
	sampleRate      int
	lowFrequencies  []int
	highFrequencies []int
	duration        time.Duration
	lastKey         rune
	lastDuration    time.Duration
	goertzel        *Goertzel
	PressInterval   time.Duration
}

func NewDTMFDecoder(energyThreshold float64, sampleRate int) *DTMFDecoder {
	lowFrequencies := []int{697, 770, 852, 941}
	highFrequencies := []int{1209, 1336, 1477, 1633}
	frequencies := append(lowFrequencies, highFrequencies...)

	decoder := &DTMFDecoder{
		energyThreshold: energyThreshold,
		sampleRate:      sampleRate,
		lowFrequencies:  lowFrequencies,
		highFrequencies: highFrequencies,
		goertzel:        NewGoertzel(frequencies, sampleRate),
		PressInterval:   200 * time.Millisecond,
	}

	return decoder
}

func (d *DTMFDecoder) Decode(samples []float64) (string, bool) {
	frameDuration := time.Duration(len(samples)) * time.Second / time.Duration(d.sampleRate)
	d.duration += frameDuration
	r, ok := d.process(samples)
	if !ok {
		return "", false
	}
	if r == d.lastKey && d.duration-d.lastDuration < d.PressInterval {
		return "", false
	}
	d.lastKey = r
	d.lastDuration = d.duration
	return string(r), ok
}

func (d *DTMFDecoder) process(samples []float64) (rune, bool) {
	sampleSize := len(samples)
	for index, sample := range samples {
		d.goertzel.processSample(d.exactBlackmanWindow(sample, index, sampleSize))
	}
	character, ok := d.energyProfileToCharacter(d.goertzel.energies)
	d.goertzel.refresh()
	return character, ok
}

func (d *DTMFDecoder) exactBlackmanWindow(sample float64, sampleIndex, bufferSize int) float64 {
	return sample * (0.426591 -
		0.496561*math.Cos((2*math.Pi*float64(sampleIndex))/float64(bufferSize)) +
		0.076849*math.Cos((4*math.Pi*float64(sampleIndex))/float64(bufferSize)))
}

func (d *DTMFDecoder) energyProfileToCharacter(energies map[int]float64) (rune, bool) {
	var lowFrequency, highFrequency int
	var lowFrequencyEnergy, highFrequencyEnergy float64

	for _, f := range d.lowFrequencies {
		if energies[f] > lowFrequencyEnergy && energies[f] >= d.energyThreshold {
			lowFrequencyEnergy = energies[f]
			lowFrequency = f
		}
	}
	if lowFrequency == 0 {
		return 0, false
	}

	for _, f := range d.highFrequencies {
		if energies[f] > highFrequencyEnergy && energies[f] >= d.energyThreshold {
			highFrequencyEnergy = energies[f]
			highFrequency = f
		}
	}
	if highFrequency == 0 {
		return 0, false
	}
	return frequencyTable[lowFrequency][highFrequency], true
}
