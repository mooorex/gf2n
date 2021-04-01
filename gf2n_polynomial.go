package gf2n

type GF2nPoly struct {
	Coeff []*GF2nElement
	Field *GF2nField
}

func (f *GF2nField) NewGF2nPoly(coeff []*GF2nElement) *GF2nPoly {
	return &GF2nPoly{coeff, f}
}

func (f *GF2nPoly) Eval(x *GF2nElement) *GF2nElement {
	res := x.f.NewGF2nElement(0)
	for i, v := range f.Coeff {
		xi := new(GF2nElement).Mul(x.Pow(uint64(i)), v)
		res = res.Add(res, xi)
	}

	return res
}
