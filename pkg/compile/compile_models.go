package compile

import (
	"fmt"
	"sort"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-types-template/pkg/formatdef"
	"github.com/kalo-build/plugin-morphe-types-template/pkg/typemap"
)

// CompileModel converts a Morphe model to the target format
func CompileModel(model yaml.Model, r *registry.Registry) (*formatdef.Struct, error) {
	// Create the struct definition
	formatStruct := &formatdef.Struct{
		Name:   model.Name,
		Fields: make([]formatdef.Field, 0),
	}

	// Sort fields for consistent output
	var fieldNames []string
	for name := range model.Fields {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	// Process model fields
	for _, fieldName := range fieldNames {
		field := model.Fields[fieldName]
		fieldType := typemap.GetFieldType(field.Type)

		formatField := formatdef.Field{
			Name: fieldName,
			Type: fieldType,
		}
		formatStruct.Fields = append(formatStruct.Fields, formatField)
	}

	// Process related models (if any)
	if len(model.Related) > 0 {
		// Sort related for consistent output
		var relatedNames []string
		for name := range model.Related {
			relatedNames = append(relatedNames, name)
		}
		sort.Strings(relatedNames)

		// Add a comment field to indicate relationships exist
		// TODO: Implement actual relationship handling for your format
		for _, relatedName := range relatedNames {
			relation := model.Related[relatedName]

			// Example: Add a foreign key field
			relField := formatdef.Field{
				Name: relatedName + "ID",
				Type: formatdef.TypeString, // Assuming string IDs
			}
			formatStruct.Fields = append(formatStruct.Fields, relField)
			_ = relation // Use relation type to determine cardinality
		}
	}

	return formatStruct, nil
}

// CompileAllModels compiles all models and writes them using the writer
func CompileAllModels(config MorpheCompileConfig, r *registry.Registry, writer *MorpheWriter) error {
	modelContents := make(map[string][]byte)

	// Process each model in the registry
	for modelName, model := range r.GetAllModels() {
		// Compile the model
		compiledModel, err := CompileModel(model, r)
		if err != nil {
			return fmt.Errorf("failed to compile model %s: %w", modelName, err)
		}

		// Generate the content for this model
		content := generateModelContent(compiledModel, config.FormatConfig)
		modelContents[modelName] = content
	}

	// Write all model contents
	return writer.WriteAllModels(modelContents)
}

// generateModelContent generates the actual file content for a model
// DEFAULT IMPLEMENTATION: Generates a simple struct-like format
// TODO: Replace this with your target format's actual syntax
func generateModelContent(model *formatdef.Struct, config _FORMAT_Config) []byte {
	cb := formatdef.NewContentBuilder("  ")

	// Add header comment
	cb.Comment("Generated model: %s", model.Name)
	cb.Comment("Fields: %d", len(model.Fields))
	cb.Line("")

	// Example patterns for different languages:
	// Python:     @dataclass
	//            class User:
	//                name: str
	//                age: int
	//
	// TypeScript: export interface User {
	//               name: string;
	//               age: number;
	//             }
	//
	// Java:       public class User {
	//               private String name;
	//               private int age;
	//               // getters/setters...
	//             }

	// Default implementation (struct-like)
	cb.Line("type %s struct {", model.Name)
	cb.Indent()

	// Group fields by type for better readability
	for _, field := range model.Fields {
		// Capitalize field names for export (common pattern)
		fieldName := formatdef.ToPascalCase(field.Name)

		// Add field with type
		cb.Line("%s %s", fieldName, field.Type.GetName())
	}

	cb.Dedent()
	cb.Line("}")
	cb.Line("")

	// Add constructor/factory function (common pattern)
	cb.Comment("Constructor for %s", model.Name)
	cb.Line("func New%s() *%s {", model.Name, model.Name)
	cb.Indent()
	cb.Line("return &%s{}", model.Name)
	cb.Dedent()
	cb.Line("}")

	return cb.Build()
}
