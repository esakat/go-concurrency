package addendum

//go:generate echo Hello, go generate!
func Foo() {
	return
}

//go:generate stringer -type=Suit
type Suit int
const (
	Spade Suit = iota
	Club
	Diamond
	Heart
)