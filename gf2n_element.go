package gf2n

type GF2nElement struct {
	v uint64     // value
	f *GF2nField // field
}

func (f *GF2nField) NewGF2nElement(v uint64) *GF2nElement {
	return &GF2nElement{f.normalize(v), f}
}

func (z *GF2nElement) SetFromValues(v uint64, field *GF2nField) *GF2nElement {
	if z == nil {
		z = &GF2nElement{field.normalize(v), field}
	} else {
		z.v = field.normalize(v)
		z.f = field
	}
	return z
}

func (z *GF2nElement) SetFromElement(x *GF2nElement) *GF2nElement {
	if z == nil {
		z = &GF2nElement{x.v, x.f}
	} else {
		z.v = x.v
		z.f = x.f
	}
	return z
}

func (z *GF2nElement) Value() uint64 {
	return z.v
}

func (z *GF2nElement) Field() *GF2nField {
	return z.f
}

// Add will add two field elements and return their sum
func (z *GF2nElement) Add(x, y *GF2nElement) *GF2nElement {
	if !x.f.Equal(y.f) {
		return nil
	}
	if z == nil {
		z = &GF2nElement{x.v ^ y.v, x.f}
	} else {
		z.v = x.v ^ y.v
		z.f = x.f
	}
	return z
}

// Substract a field element from another element and return their difference
func (z *GF2nElement) Sub(x, y *GF2nElement) *GF2nElement {
	return z.Add(x, y)
}

// Multiply two field elements and return their product
func (z *GF2nElement) Mul(x, y *GF2nElement) *GF2nElement {
	if !x.f.Equal(y.f) {
		return nil
	}
	// irre = a_0 + a_1*z + ... + a_{n-1}*e^{n-1}
	// e^n = irre (mod e^n + irre)
	// e^(n+1) = a_{n-1} * irre (mod e^n + irre)
	var f uint64 = x.v
	var big uint64 = 1<<x.f.n - 1
	var ret uint64 = 0
	for i := uint8(0); i < x.f.n; i++ {
		t0 := f * ((y.v >> i) & 1)
		ret = ret ^ t0
		t1 := (f << 1) & big
		t2 := 1 & (f >> (x.f.n - 1))
		f = t1 ^ (t2 * x.f.irre)
	}
	return z.SetFromValues(ret, x.f)
}

// Return the power of a field element to a power
func (z *GF2nElement) Pow(n uint64) *GF2nElement {
	res := z.f.NewGF2nElement(1)
	a := new(GF2nElement).SetFromElement(z)
	for n > 0 {
		if n&1 == 1 {
			res = res.Mul(res, a)
		}
		a.Mul(a, a)
		n = n >> 1
	}
	return res
}

//return whether two elements are equal
func (z *GF2nElement) Equal(x *GF2nElement) bool {
	return z.v == x.v && z.f.Equal(x.f)
}
