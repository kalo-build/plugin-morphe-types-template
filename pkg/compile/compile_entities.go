package compile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-types-template/pkg/formatdef"
	"github.com/kalo-build/plugin-morphe-types-template/pkg/typemap"
)

// CompileEntity converts a Morphe entity to the target format
func CompileEntity(entity yaml.Entity, r *registry.Registry) (*formatdef.Struct, error) {
	// Create the struct definition
	formatStruct := &formatdef.Struct{
		Name:   entity.Name,
		Fields: make([]formatdef.Field, 0),
	}

	// Sort fields for consistent output
	var fieldNames []string
	for name := range entity.Fields {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	// Process entity fields
	for _, fieldName := range fieldNames {
		field := entity.Fields[fieldName]
		fieldType, err := resolveEntityFieldType(field.Type, r)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve field type for %s: %w", fieldName, err)
		}

		formatField := formatdef.Field{
			Name: fieldName,
			Type: fieldType,
		}
		formatStruct.Fields = append(formatStruct.Fields, formatField)
	}

	// Sort and process relationships
	if len(entity.Related) > 0 {
		var relatedNames []string
		for name := range entity.Related {
			relatedNames = append(relatedNames, name)
		}
		sort.Strings(relatedNames)

		for _, relatedName := range relatedNames {
			relation := entity.Related[relatedName]

			// Add foreign key field
			fkField := formatdef.Field{
				Name: relatedName + "ID",
				Type: formatdef.TypeString,
			}
			formatStruct.Fields = append(formatStruct.Fields, fkField)

			// Add navigation field based on relation type
			var navType formatdef.Type
			switch relation.Type {
			case "HasMany", "ForMany":
				navType = formatdef.ArrayType{
					ElementType: formatdef.BasicType{Name: relatedName},
				}
			default:
				navType = formatdef.BasicType{Name: relatedName}
			}

			navField := formatdef.Field{
				Name: relatedName,
				Type: navType,
			}
			formatStruct.Fields = append(formatStruct.Fields, navField)
		}
	}

	return formatStruct, nil
}

// resolveEntityFieldType resolves a model field path to a concrete type
func resolveEntityFieldType(fieldPath yaml.ModelFieldPath, r *registry.Registry) (formatdef.Type, error) {
	// Split the path (e.g., "User.email" -> ["User", "email"])
	parts := strings.Split(string(fieldPath), ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid field path: %s", fieldPath)
	}

	// Get the root model
	modelName := parts[0]
	model, err := r.GetModel(modelName)
	if err != nil {
		return nil, fmt.Errorf("model %s not found", modelName)
	}

	// Navigate through the path to find the terminal field
	fieldName := parts[len(parts)-1]
	field, exists := model.Fields[fieldName]
	if !exists {
		return nil, fmt.Errorf("field %s not found in model %s", fieldName, modelName)
	}

	// Return the appropriate type
	return typemap.GetFieldType(field.Type), nil
}

// CompileAllEntities compiles all entities and writes them using the writer
func CompileAllEntities(config MorpheCompileConfig, r *registry.Registry, writer *MorpheWriter) error {
	entityContents := make(map[string][]byte)

	// Process each entity in the registry
	for entityName, entity := range r.GetAllEntities() {
		// Compile the entity
		compiledEntity, err := CompileEntity(entity, r)
		if err != nil {
			return fmt.Errorf("failed to compile entity %s: %w", entityName, err)
		}

		// Generate the content for this entity
		content := generateEntityContent(compiledEntity, entity, config.FormatConfig)
		entityContents[entityName] = content
	}

	// Write all entity contents
	return writer.WriteAllEntities(entityContents)
}

// generateEntityContent generates the actual file content for an entity
// DEFAULT IMPLEMENTATION: Generates an entity with relationships and identifiers
// TODO: Replace this with your target format's actual syntax
func generateEntityContent(entity *formatdef.Struct, morpheEntity yaml.Entity, config _FORMAT_Config) []byte {
	cb := formatdef.NewContentBuilder("  ")

	// Add comprehensive documentation
	cb.BlockComment(
		fmt.Sprintf("%s entity", entity.Name),
		"",
		"This entity aggregates data from multiple models",
		fmt.Sprintf("Identifiers: %d", len(morpheEntity.Identifiers)),
		fmt.Sprintf("Relationships: %d", len(morpheEntity.Related)),
	)
	cb.Line("")

	// Example patterns for different languages:
	// Python:     @entity
	//            class UserEntity:
	//                id: str  # primary identifier
	//                email: str  # unique identifier
	//                profile: Profile  # relationship
	//
	// TypeScript: export interface UserEntity {
	//               id: string; // primary
	//               email: string; // unique
	//               profile?: Profile; // optional relationship
	//             }
	//
	// Java:       @Entity
	//            public class UserEntity {
	//                @Id
	//                private String id;
	//
	//                @Unique
	//                private String email;
	//
	//                @OneToOne
	//                private Profile profile;
	//             }

	// Default implementation
	cb.Line("type %s struct {", entity.Name)
	cb.Indent()

	// Add fields with identifier annotations
	for _, field := range entity.Fields {
		fieldName := formatdef.ToPascalCase(field.Name)

		// Check if this field is part of an identifier
		var identifierType string
		for idName, identifier := range morpheEntity.Identifiers {
			for _, idField := range identifier.Fields {
				if idField == field.Name {
					identifierType = idName
					break
				}
			}
		}

		if identifierType != "" {
			cb.Line("%s %s // %s identifier", fieldName, field.Type.GetName(), identifierType)
		} else {
			cb.Line("%s %s", fieldName, field.Type.GetName())
		}
	}

	cb.Dedent()
	cb.Line("}")
	cb.Line("")

	// Add identifier methods
	if len(morpheEntity.Identifiers) > 0 {
		cb.Comment("Identifier methods")
		cb.Line("")

		// Generate method for primary identifier
		if primary, hasPrimary := morpheEntity.Identifiers["primary"]; hasPrimary && len(primary.Fields) > 0 {
			cb.Line("func (e *%s) GetID() string {", entity.Name)
			cb.Indent()
			cb.Line("return e.%s", formatdef.ToPascalCase(primary.Fields[0]))
			cb.Dedent()
			cb.Line("}")
			cb.Line("")
		}
	}

	// Add relationship helper methods
	if len(morpheEntity.Related) > 0 {
		cb.Comment("Relationship loaders")
		cb.Line("")

		for relName, relation := range morpheEntity.Related {
			switch relation.Type {
			case "HasMany", "ForMany":
				cb.Line("func (e *%s) Load%s() ([]%s, error) {", entity.Name, formatdef.ToPascalCase(relName), relName)
				cb.Indent()
				cb.Line("// TODO: Implement lazy loading")
				cb.Line("return nil, nil")
				cb.Dedent()
				cb.Line("}")
			default:
				cb.Line("func (e *%s) Load%s() (*%s, error) {", entity.Name, formatdef.ToPascalCase(relName), relName)
				cb.Indent()
				cb.Line("// TODO: Implement lazy loading")
				cb.Line("return nil, nil")
				cb.Dedent()
				cb.Line("}")
			}
			cb.Line("")
		}
	}

	return cb.Build()
}
