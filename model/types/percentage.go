package types

import (
	"database/sql/driver"
	"errors"
	"math/big"
)

var (
	errInvalidPercentage = errors.New("invalid percentage")
)

type Percentage struct {
	*big.Float
}

// NewPercentage returns a new percentage from the given string
func NewPercentage(src string) Percentage {
	if src == "" {
		src = "0"
	}

	n := new(big.Float)
	n.SetString(src)

	return Percentage{Float: n}
}

// NewFloat64Percentage returns a new percentage for the given float64 value
func NewFloat64Percentage(val float64) Percentage {
	n := big.NewFloat(val)
	return Percentage{Float: n}
}

// Value returns a serialized value
func (a Percentage) Value() (driver.Value, error) {
	if a.Float != nil {
		return a.Float.String(), nil
	}
	return nil, nil
}

func (a Percentage) String() string {
	if a.Float == nil {
		return ""
	}
	return a.Float.String()
}

// Add adds two numbers
func (a Percentage) Add(b Percentage) Percentage {
	n := new(big.Float)
	n = n.Add(a.Float, b.Float)
	return Percentage{n}
}

// Sub substitutes a given percentage amount from the current one
func (a Percentage) Sub(b Percentage) Percentage {
	n := new(big.Float)
	n = n.Sub(a.Float, b.Float)
	return Percentage{n}
}

// Mul multiplies two numbers
func (a Percentage) Mul(b Percentage) Percentage {
	n := new(big.Float)
	n = n.Mul(a.Float, b.Float)
	return Percentage{n}
}

// Quo divides two numbers
func (a Percentage) Quo(b Percentage) Percentage {
	n := new(big.Float)
	n = n.Quo(a.Float, b.Float)
	return Percentage{n}
}

// Scan assigns the value from interface
func (a *Percentage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	n := new(big.Float)

	switch v := value.(type) {
	case float64:
		a.Float = n.SetFloat64(v)
	case string:
		n, ok := n.SetString(v)
		if !ok {
			return errInvalidPercentage
		}
		a.Float = n
	case []byte:
		n, ok := n.SetString(string(value.([]byte)))
		if !ok {
			return errInvalidPercentage
		}
		a.Float = n
	}

	return nil
}
