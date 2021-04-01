package gf2n

import (
	"errors"
)

// arithmetic over GF(2^n)
// n must be less than 64
// x^n + irre must be a irreducible polynomial over GF(2^n)
// where irre is a polynomial over GF(2) of degree less than n
type GF2nField struct {
	n    uint8
	irre uint64
}

// generate a binary extension field
// we do not check whether x^n + irre is irreducible
func NewGF2nField(n uint8, irre uint64) (*GF2nField, error) {
	if n >= 64 || irre > (1<<n) {
		return nil, errors.New("invalid parameters.")
	}
	return &GF2nField{n, irre}, nil
}

// return the degree of extension
func (f *GF2nField) ExtDegree() uint8 {
	return f.n
}

// return irreducible polynomial wo/ the highest order term
func (f *GF2nField) IrrePolyWithoutHighest() uint64 {
	return f.irre
}

func (f *GF2nField) Equal(e *GF2nField) bool {
	return f.n == e.n && f.irre == e.irre
}

// return an integer which is less than 2^n
func (f *GF2nField) normalize(v uint64) uint64 {
	var big uint64 = (1 << f.n) - 1
	var t uint64 = v >> f.n
	for t > 0 {
		v = (v & big) ^ (t * f.irre)
		t = v >> f.n
	}
	return v
}
