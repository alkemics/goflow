package nodes

// IgnoredType have all it's methods ignored.
// node:ignore
type IgnoredType struct{}

func NewIgnoredType() *IgnoredType {
	return &IgnoredType{}
}

func (n IgnoredType) Run() {}

type IgnoredMethod struct{}

func NewIgnoredMethod() IgnoredMethod {
	return IgnoredMethod{}
}

func (n IgnoredMethod) Create() {}

// Delete is an ignored method.
// node:ignore
func (n IgnoredMethod) Delete() {}

// IgnoredFunction is an ignored function.
// node:ignore
func IgnoredFunction() {}
