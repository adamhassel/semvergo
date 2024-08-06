package flags

import (
	"fmt"
	"strconv"
)

type String struct {
	set   bool
	value string
}

type Bool struct {
	set   bool
	value bool
}

func (sf *String) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *String) String() string {
	return sf.value
}

func (sf *String) IsSet() bool {
	return sf.set
}

func (bf *Bool) IsBoolFlag() bool {
	return true
}

func (bf *Bool) Set(x string) error {
	var err error
	bf.value, err = strconv.ParseBool(x)
	if err != nil {
		return err
	}
	bf.set = true
	return nil
}

func (bf *Bool) String() string {
	return fmt.Sprintf("%t", bf.value)
}

func (bf *Bool) Bool() bool {
	return bf.value
}

func (bf *Bool) IsSet() bool {
	return bf.set
}
