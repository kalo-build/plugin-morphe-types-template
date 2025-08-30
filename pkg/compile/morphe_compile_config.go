package compile

import (
	"path"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
)

// MorpheCompileConfig contains all configuration for compiling Morphe to the target format
type MorpheCompileConfig struct {
	// Registry loading configuration
	rcfg.MorpheLoadRegistryConfig

	// Output path for generated files
	OutputPath string

	// Format-specific configuration
	FormatConfig _FORMAT_Config
}

// _FORMAT_Config contains format-specific configuration options
// TODO: Replace _FORMAT_ with your target format name
type _FORMAT_Config struct {
	// TODO: Add format-specific configuration options
	// Examples:
	// - PackagePrefix string (for Java/Python packages)
	// - ModuleName string (for JavaScript/TypeScript modules)
	// - IndentSize int
	// - UseTabsForIndent bool
	// - GenerateComments bool
	// - FileExtension string
}

// DefaultMorpheCompileConfig creates a default configuration
func DefaultMorpheCompileConfig(
	yamlRegistryPath string,
	baseOutputDirPath string,
) MorpheCompileConfig {
	return MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      path.Join(yamlRegistryPath, "enums"),
			RegistryModelsDirPath:     path.Join(yamlRegistryPath, "models"),
			RegistryStructuresDirPath: path.Join(yamlRegistryPath, "structures"),
			RegistryEntitiesDirPath:   path.Join(yamlRegistryPath, "entities"),
		},
		OutputPath:   baseOutputDirPath,
		FormatConfig: _FORMAT_Config{
			// TODO: Set default format-specific configuration
		},
	}
}

// Validate checks if the configuration is valid
func (config MorpheCompileConfig) Validate() error {
	// Validate registry paths
	if err := config.MorpheLoadRegistryConfig.Validate(); err != nil {
		return err
	}

	// TODO: Add format-specific validation
	// Examples:
	// - Check if package prefix is valid
	// - Verify indent size is positive
	// - Ensure file extension starts with "."

	return nil
}
