package nodes

// DocumentedType is a type with documentation in the main package.
type DocumentedType struct{}

// NewDocumentedType creates a new DocumentedType.
func NewDocumentedType() *DocumentedType {
	return &DocumentedType{}
}

// MainMethod has documentation in the main package.
func (d *DocumentedType) MainMethod() {}

// DocumentedFunction is a function with documentation.
func DocumentedFunction() {}
