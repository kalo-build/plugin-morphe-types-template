package compile

import (
	"fmt"
	"sort"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-types-template/pkg/formatdef"
)

// CompileEnum converts a Morphe enum to the target format
func CompileEnum(enum yaml.Enum) (*formatdef.Enum, error) {
	// Create the enum definition
	formatEnum := &formatdef.Enum{
		Name:    enum.Name,
		Type:    mapEnumType(enum.Type),
		Entries: make([]formatdef.EnumEntry, 0, len(enum.Entries)),
	}

	// Sort entries for consistent output
	var entryNames []string
	for name := range enum.Entries {
		entryNames = append(entryNames, name)
	}
	sort.Strings(entryNames)

	// Convert each enum entry
	for _, entryName := range entryNames {
		entry := formatdef.EnumEntry{
			Name:  entryName,
			Value: enum.Entries[entryName],
		}
		formatEnum.Entries = append(formatEnum.Entries, entry)
	}

	return formatEnum, nil
}

// mapEnumType maps Morphe enum types to format-specific types
func mapEnumType(morpheType yaml.EnumType) formatdef.Type {
	switch morpheType {
	case yaml.EnumTypeInteger:
		return formatdef.TypeInteger
	case yaml.EnumTypeFloat:
		return formatdef.TypeFloat
	case yaml.EnumTypeString:
		return formatdef.TypeString
	default:
		return formatdef.TypeString
	}
}

// CompileAllEnums compiles all enums and writes them using the writer
func CompileAllEnums(config MorpheCompileConfig, r *registry.Registry, writer *MorpheWriter) error {
	enumContents := make(map[string][]byte)

	// Process each enum in the registry
	for enumName, enum := range r.GetAllEnums() {
		// Compile the enum
		compiledEnum, err := CompileEnum(enum)
		if err != nil {
			return fmt.Errorf("failed to compile enum %s: %w", enumName, err)
		}

		// Generate the content for this enum
		content := generateEnumContent(compiledEnum, config.FormatConfig)
		enumContents[enumName] = content
	}

	// Write all enum contents
	return writer.WriteAllEnums(enumContents)
}

// generateEnumContent generates the actual file content for an enum
// DEFAULT IMPLEMENTATION: Generates a simple JSON-like format
// TODO: Replace this with your target format's actual syntax
func generateEnumContent(enum *formatdef.Enum, config _FORMAT_Config) []byte {
	cb := formatdef.NewContentBuilder("  ")

	// Example patterns for different languages (in comments):
	// Python:   class MyEnum(Enum):
	// TypeScript: export enum MyEnum {
	// Java:     public enum MyEnum {
	// Go:       type MyEnum int

	// Default implementation (JSON-like for demonstration)
	cb.Comment("Generated enum: %s", enum.Name)
	cb.Comment("Type: %s", enum.Type.GetName())
	cb.Line("")
	cb.Line("enum %s {", enum.Name)
	cb.Indent()

	for i, entry := range enum.Entries {
		// Different formats for different types
		switch enum.Type.GetName() {
		case "string":
			cb.Line("%s = %q", entry.Name, entry.Value)
		default:
			cb.Line("%s = %v", entry.Name, entry.Value)
		}

		// Add comma for all but last entry (common pattern)
		if i < len(enum.Entries)-1 {
			cb.AppendToLastLine(",")
		}
	}

	cb.Dedent()
	cb.Line("}")

	return cb.Build()
}
