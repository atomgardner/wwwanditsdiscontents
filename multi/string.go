package multi

import "fmt"

type String []string

func (x String) String() string {
	return fmt.Sprintf("%#v", x)
}

func (x *String) Set(val string) error {
	*x = append(*x, val)
	return nil
}
