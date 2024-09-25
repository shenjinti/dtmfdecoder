package dtmfdecoder

import (
	"math"
)

type Goertzel struct {
	frequencies    []int
	sampleRate     int
	coefficient    map[int]float64
	firstPrevious  map[int]float64
	secondPrevious map[int]float64
	filterLength   map[int]int
	totalPower     map[int]float64
	energies       map[int]float64
}

func NewGoertzel(frequencies []int, sampleRate int) *Goertzel {
	g := &Goertzel{
		frequencies:    frequencies,
		sampleRate:     sampleRate,
		coefficient:    make(map[int]float64),
		firstPrevious:  make(map[int]float64),
		secondPrevious: make(map[int]float64),
		filterLength:   make(map[int]int),
		totalPower:     make(map[int]float64),
		energies:       make(map[int]float64),
	}
	g.initialize()
	g.refresh()
	return g
}

func (g *Goertzel) initialize() {
	for _, frequency := range g.frequencies {
		normalizedFrequency := float64(frequency) / float64(g.sampleRate)
		omega := 2 * math.Pi * normalizedFrequency
		cosine := math.Cos(omega)
		g.coefficient[frequency] = 2 * cosine
	}
}

func (g *Goertzel) processSample(sample float64) {
	for _, frequency := range g.frequencies {
		g.getEnergyOfFrequency(sample, frequency)
	}
}

func (g *Goertzel) refresh() {
	g.firstPrevious = make(map[int]float64)
	g.secondPrevious = make(map[int]float64)
	g.filterLength = make(map[int]int)
	g.totalPower = make(map[int]float64)
	g.energies = make(map[int]float64)
}

func (g *Goertzel) getEnergyOfFrequency(sample float64, frequency int) {
	f1 := g.firstPrevious[frequency]
	f2 := g.secondPrevious[frequency]
	coefficient := g.coefficient[frequency]
	sine := sample + (coefficient * f1) - f2
	f2 = f1
	f1 = sine
	g.filterLength[frequency] += 1
	power := (f2 * f2) + (f1 * f1) - (coefficient * f1 * f2)
	totalPower := g.totalPower[frequency] + sample*sample
	if totalPower == 0 {
		totalPower = 1
	}
	g.energies[frequency] = power / totalPower / float64(g.filterLength[frequency])
	g.firstPrevious[frequency] = f1
	g.secondPrevious[frequency] = f2
	g.totalPower[frequency] = totalPower
}
