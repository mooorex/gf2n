package gf2n

import (
	"testing"
)

func TestGF2nElement(t *testing.T) {
	c := uint8(16)
	poly := uint64(1<<12 + 1<<3 + 1<<1 + 1)
	// poly = x^12 + x^3 + x + 1
	// x^c + poly may be irreducible
	field, err := NewGF2nField(c, poly)
	if err != nil {
		t.Fatal(err)
	}

	t1 := field.NewGF2nElement(1<<2 + 1<<3)
	t2 := field.NewGF2nElement(1<<15 + 1<<7)
	t3 := field.NewGF2nElement(1<<2 + 1<<3 + 1<<16 + 1<<17 + 1<<63)
	t4 := field.NewGF2nElement(1<<2 + 1<<3 + poly*(1<<0+1<<1+1<<47))

	if new(GF2nElement).Add(t1, t1).v != uint64(0) ||
		new(GF2nElement).Add(t3, t4).v != uint64(0) ||
		new(GF2nElement).Add(t1, t2).v != uint64(32908) {
		t.Fatalf("add err!")
	}

	if !t4.Equal(t3) {
		t.Fatalf("compare err!")
	}

	// (x^2 + x^3) * (x^15 + x^7)
	//= x^17 + x^10 + x^18 + x^9
	//= x + x^3 + x^4 + x^5 + x^9 + x^10 + x^13 + x^14
	//= 26170
	if new(GF2nElement).Mul(t1, t2).v != 26170 {
		t.Fatalf("mul err!")
	}

	var power uint64 = 100
	f00 := field.NewGF2nElement(1)
	for i := uint64(0); i < power; i++ {
		f00 = new(GF2nElement).Mul(f00, t1)
	}
	f01 := t1.Pow(power)

	if f00.v != f01.v {
		t.Fatalf("pow err!")
	}

	// (f0 + f1)*f2
	//=f0*f2 + f1*f2
	f0 := field.NewGF2nElement(1<<2 + 1<<3)
	f1 := field.NewGF2nElement(1<<5 + 1<<7)
	f2 := field.NewGF2nElement(1<<13 + 1<<7)

	f0f2 := new(GF2nElement).Mul(f0, f2)
	f1f2 := new(GF2nElement).Mul(f1, f2)

	rhs := new(GF2nElement).Mul(new(GF2nElement).Add(f0, f1), f2)
	lhs := new(GF2nElement).Add(f0f2, f1f2)
	if lhs.v != rhs.v {
		t.Fatalf("distribution law err: %d \t %d", lhs.v, rhs.v)
	}

	poly1 := uint64(1<<1 + 1)
	// poly =  x + 1
	// x^c + poly1 may be irreducible
	field1, err := NewGF2nField(c, poly1)
	if err != nil {
		t.Fatal(err)
	}
	ff0 := field1.NewGF2nElement(1<<2 + 1<<3)
	ret := new(GF2nElement).Add(f0, ff0)
	if ret != nil {
		t.Fatalf("field check err!")
	}
}

func BenchmarkGF2nElementAdd(b *testing.B) {
	c := uint8(16)
	poly := uint64(4107)
	field, _ := NewGF2nField(c, poly)
	f0 := field.NewGF2nElement(1<<2 + 1<<3)
	f1 := field.NewGF2nElement(1<<15 + 1<<7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f0.Add(f0, f1)
	}
}

func BenchmarkGF2nElementMul(b *testing.B) {
	c := uint8(16)
	poly := uint64(4107)
	field, _ := NewGF2nField(c, poly)
	f0 := field.NewGF2nElement(1<<2 + 1<<3)
	f1 := field.NewGF2nElement(1<<15 + 1<<7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f0.Mul(f0, f1)
	}
}

func BenchmarkGF2nElementPow(b *testing.B) {
	c := uint8(16)
	poly := uint64(4107)
	field, _ := NewGF2nField(c, poly)
	f0 := field.NewGF2nElement(1<<2 + 1<<3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f0.Pow(uint64(10000000))
	}
}

func BenchmarkGF2nPolyEval(b *testing.B) {
	c := uint8(16)
	poly := uint64(4107)
	field, _ := NewGF2nField(c, poly)

	f0 := field.NewGF2nElement(1<<2 + 1<<3)
	f1 := field.NewGF2nElement(1<<5 + 1<<7)
	f2 := field.NewGF2nElement(1<<13 + 1<<7)
	cf := []*GF2nElement{f0, f1, f2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.NewGF2nPoly(cf).Eval(f0)
	}
}
