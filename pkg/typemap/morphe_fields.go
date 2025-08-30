package typemap

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-types-template/pkg/formatdef"
)

// MorpheModelFieldToFormatType maps Morphe field types to target format types
// TODO: Rename this variable to match your format (e.g., MorpheModelFieldToPythonType)
// TODO: Update the type mappings to match your target format's type system
var MorpheModelFieldToFormatType = map[yaml.ModelFieldType]formatdef.Type{
	// String types
	yaml.ModelFieldTypeString:    formatdef.TypeString,
	yaml.ModelFieldTypeUUID:      formatdef.TypeString,
	yaml.ModelFieldTypeProtected: formatdef.TypeString,
	yaml.ModelFieldTypeSealed:    formatdef.TypeString,

	// Numeric types
	yaml.ModelFieldTypeInteger:       formatdef.TypeInteger,
	yaml.ModelFieldTypeAutoIncrement: formatdef.TypeInteger,
	yaml.ModelFieldTypeFloat:         formatdef.TypeFloat,

	// Boolean type
	yaml.ModelFieldTypeBoolean: formatdef.TypeBoolean,

	// Date/Time types
	yaml.ModelFieldTypeTime: formatdef.TypeDate,
	yaml.ModelFieldTypeDate: formatdef.TypeDate,

	// TODO: Add mappings for any custom field types used in your Morphe schemas
}

// GetFieldType returns the format type for a given Morphe field type
// TODO: Add any format-specific type conversion logic
func GetFieldType(fieldType yaml.ModelFieldType) formatdef.Type {
	if formatType, exists := MorpheModelFieldToFormatType[fieldType]; exists {
		return formatType
	}
	// Default to string type for unknown field types
	// TODO: Consider whether to return an error instead
	return formatdef.TypeString
}
