// Code generated by "go generate gonum.org/v1/gonum/unit”; DO NOT EDIT.

// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf8"
)

// Inductance represents an electrical inductance in Henry.
type Inductance float64

const Henry Inductance = 1

// Unit converts the Inductance to a *Unit
func (i Inductance) Unit() *Unit {
	return New(float64(i), Dimensions{
		CurrentDim: -2,
		LengthDim:  2,
		MassDim:    1,
		TimeDim:    -2,
	})
}

// Inductance allows Inductance to implement a Inductancer interface
func (i Inductance) Inductance() Inductance {
	return i
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension
func (i *Inductance) From(u Uniter) error {
	if !DimensionsMatch(u, Henry) {
		*i = Inductance(math.NaN())
		return errors.New("Dimension mismatch")
	}
	*i = Inductance(u.Unit().Value())
	return nil
}

func (i Inductance) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", i, float64(i))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		const unit = " H"
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), p, float64(i))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(i))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), float64(i))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(i))
		}
		fmt.Fprint(fs, unit)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g H)", c, i, float64(i))
	}
}
