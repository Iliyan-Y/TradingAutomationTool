package orders

type Orders interface {
	Buy() bool	// todo
	Sell() bool // error ?
}