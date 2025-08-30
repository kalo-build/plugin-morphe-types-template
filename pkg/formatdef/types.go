package formatdef

// Type represents a type in the target format
// TODO: Replace with your target format's type system
type Type interface {
	// GetName returns the type name in the target format
	GetName() string
	// IsNullable returns whether the type can be null/nil/undefined
	IsNullable() bool
	// TODO: Add format-specific methods as needed
}

// BasicType represents a basic/primitive type
type BasicType struct {
	Name     string
	Nullable bool
}

func (t BasicType) GetName() string {
	return t.Name
}

func (t BasicType) IsNullable() bool {
	return t.Nullable
}

// ArrayType represents an array/list type
type ArrayType struct {
	ElementType Type
	// TODO: Add format-specific array properties
}

func (t ArrayType) GetName() string {
	// TODO: Implement format-specific array syntax
	return "Array<" + t.ElementType.GetName() + ">"
}

func (t ArrayType) IsNullable() bool {
	return false
}

// TODO: Define common basic types for your format
// Replace these with your format's actual type names
var (
	TypeString  = BasicType{Name: "string"}
	TypeInteger = BasicType{Name: "integer"}
	TypeFloat   = BasicType{Name: "float"}
	TypeBoolean = BasicType{Name: "boolean"}
	TypeDate    = BasicType{Name: "date"}
	TypeJSON    = BasicType{Name: "json"}
)
