package compile

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-types-template/pkg/formatdef"
)

// CompileStructure converts a Morphe structure to the target format
func CompileStructure(structure yaml.Structure, r *registry.Registry) (*formatdef.Struct, error) {
	// Create the struct definition
	formatStruct := &formatdef.Struct{
		Name:   structure.Name,
		Fields: make([]formatdef.Field, 0),
	}

	// Structures are flexible containers - add common fields
	// Most formats will want at least an ID and data field

	// Add ID field
	idField := formatdef.Field{
		Name: "ID",
		Type: formatdef.TypeString,
	}
	formatStruct.Fields = append(formatStruct.Fields, idField)

	// Add data field for flexible storage
	dataField := formatdef.Field{
		Name: "Data",
		Type: formatdef.TypeJSON,
	}
	formatStruct.Fields = append(formatStruct.Fields, dataField)

	// Add timestamp fields (common pattern)
	createdField := formatdef.Field{
		Name: "CreatedAt",
		Type: formatdef.TypeDate,
	}
	formatStruct.Fields = append(formatStruct.Fields, createdField)

	updatedField := formatdef.Field{
		Name: "UpdatedAt",
		Type: formatdef.TypeDate,
	}
	formatStruct.Fields = append(formatStruct.Fields, updatedField)

	return formatStruct, nil
}

// CompileAllStructures compiles all structures and writes them using the writer
func CompileAllStructures(config MorpheCompileConfig, r *registry.Registry, writer *MorpheWriter) error {
	structureContents := make(map[string][]byte)

	// Process each structure in the registry
	for structureName, structure := range r.GetAllStructures() {
		// Compile the structure
		compiledStructure, err := CompileStructure(structure, r)
		if err != nil {
			return fmt.Errorf("failed to compile structure %s: %w", structureName, err)
		}

		// Generate the content for this structure
		content := generateStructureContent(compiledStructure, config.FormatConfig)
		structureContents[structureName] = content
	}

	// Write all structure contents
	return writer.WriteAllStructures(structureContents)
}

// generateStructureContent generates the actual file content for a structure
// DEFAULT IMPLEMENTATION: Generates a flexible container type
// TODO: Replace this with your target format's actual syntax
func generateStructureContent(structure *formatdef.Struct, config _FORMAT_Config) []byte {
	cb := formatdef.NewContentBuilder("  ")

	// Add documentation
	cb.BlockComment(
		fmt.Sprintf("%s is a flexible data structure", structure.Name),
		"It can store arbitrary data in a type-safe manner",
		"",
		"Generated from Morphe structure definition",
	)
	cb.Line("")

	// Example patterns for different languages:
	// Python:     class MyStructure:
	//               def __init__(self, data: dict):
	//                 self.id = str(uuid.uuid4())
	//                 self.data = data
	//                 self.created_at = datetime.now()
	//
	// TypeScript: export class MyStructure {
	//               id: string;
	//               data: Record<string, any>;
	//               createdAt: Date;
	//
	//               constructor(data: Record<string, any>) {
	//                 this.id = generateId();
	//                 this.data = data;
	//                 this.createdAt = new Date();
	//               }
	//             }

	// Default implementation
	cb.Line("type %s struct {", structure.Name)
	cb.Indent()

	for _, field := range structure.Fields {
		cb.Line("%s %s", field.Name, field.Type.GetName())
	}

	cb.Dedent()
	cb.Line("}")
	cb.Line("")

	// Add helper methods
	cb.Comment("NewStructure creates a new %s instance", structure.Name)
	cb.Line("func New%s(data map[string]interface{}) *%s {", structure.Name, structure.Name)
	cb.Indent()
	cb.Line("return &%s{", structure.Name)
	cb.Indent()
	cb.Line("ID:        generateID(),")
	cb.Line("Data:      data,")
	cb.Line("CreatedAt: time.Now(),")
	cb.Line("UpdatedAt: time.Now(),")
	cb.Dedent()
	cb.Line("}")
	cb.Dedent()
	cb.Line("}")
	cb.Line("")

	// Add getter/setter example
	cb.Comment("Get retrieves a value from the structure's data")
	cb.Line("func (s *%s) Get(key string) (interface{}, bool) {", structure.Name)
	cb.Indent()
	cb.Line("val, ok := s.Data[key]")
	cb.Line("return val, ok")
	cb.Dedent()
	cb.Line("}")

	return cb.Build()
}
