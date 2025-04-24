package message

type Address struct {
	name string
	addr string
}

func NewAddress(name, addr string) Address {
	return Address{
		name: name,
		addr: addr,
	}
}

func (a Address) String() string {
	if a.name == "" {
		return a.addr
	}
	return a.name + " <" + a.addr + ">"
}

func (a Address) Addr() string {
	return a.addr
}

func (a Address) Name() string {
	return a.name
}
