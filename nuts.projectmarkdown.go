package gonuts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/lexers"
	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

type MarkdownGeneratorConfig struct {
	BaseDir         *string  `yaml:"baseDir"`
	Includes        []string `yaml:"includes"`
	Excludes        []string `yaml:"excludes"`
	BaseHeaderLevel *int     `yaml:"baseHeaderLevel"`
	PrependText     *string  `yaml:"prependText"`
	Title           *string  `yaml:"title"`
	OutputFile      *string  `yaml:"outputFile"`
}

var (
	defaultExcludes = []string{
		"**/vendor/**",
		"**/node_modules/**",
		"**/.git/**",
		"**/.vault/**",
	}

	defaultIncludes = []string{
		"**/*.go", "**/*.js", "**/*.py", "**/*.java", "**/*.c", "**/*.cpp", "**/*.h",
		"**/*.cs", "**/*.rb", "**/*.php", "**/*.swift", "**/*.kt", "**/*.rs",
		"**/*.ts", "**/*.jsx", "**/*.tsx", "**/*.vue", "**/*.scala", "**/*.groovy",
		"**/*.sh", "**/*.bash", "**/*.zsh", "**/*.fish",
		"**/*.sql", "**/*.md", "**/*.yaml", "**/*.yml", "**/*.json", "**/*.xml",
		"**/*.html", "**/*.css", "**/*.scss", "**/*.sass", "**/*.less",
		"**/Dockerfile", "**/Makefile", "**/Jenkinsfile", "**/Gemfile",
		"**/.gitignore", "**/.dockerignore",
	}

	defaultBaseHeaderLevel = 2
	defaultPrependText     = "This is a generated markdown file containing code from the project."
	defaultTitle           = "Project Code Documentation"
)

// // EXAMPLE: Load config from YAML
// config, err := LoadConfigFromYAML("markdown_config.yaml")
// if err != nil {
//     log.Fatalf("Error loading config: %v", err)
// }

// // Generate markdown
// err = GenerateMarkdownFromFiles(config)
// if err != nil {
//     log.Fatalf("Error generating markdown: %v", err)
// }

// // markdown_config.yaml
// baseDir: /path/to/your/project
// includes:
//   - "**/*.go"
//   - "**/*.md"
// excludes:
//   - "**/test/**"

func LoadConfigFromYAML(filename string) (*MarkdownGeneratorConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config MarkdownGeneratorConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return &config, nil
}

func ApplyDefaults(config *MarkdownGeneratorConfig) error {
	if config.BaseDir == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %w", err)
		}
		config.BaseDir = &cwd
	}

	if config.BaseHeaderLevel == nil {
		config.BaseHeaderLevel = &defaultBaseHeaderLevel
	}

	if len(config.Excludes) == 0 {
		config.Excludes = defaultExcludes
	}

	if len(config.Includes) == 0 {
		config.Includes = defaultIncludes
	}

	if config.PrependText == nil {
		config.PrependText = &defaultPrependText
	}

	if config.Title == nil {
		config.Title = &defaultTitle
	}

	if config.OutputFile == nil {
		defaultOutputFile := filepath.Join(*config.BaseDir, "project_code.md")
		config.OutputFile = &defaultOutputFile
	}

	return nil
}

// GenerateMarkdownFromFiles generates a markdown file containing code snippets from files in a directory.
//
//	config := &gonuts.MarkdownGeneratorConfig{
//		BaseDir: gonuts.Ptr("/path/to/your/project"),
//		Includes: []string{
//				"**/*.go",
//				"**/*.md",
//				"**/*.yaml",
//		},
//		Excludes: []string{
//				"**/vendor/**",
//				"**/test/**",
//		},
//		Title: gonuts.Ptr("My Custom Project Documentation"),
//		BaseHeaderLevel: gonuts.Ptr(3),
//	}
func GenerateMarkdownFromFiles(config *MarkdownGeneratorConfig) error {
	err := ApplyDefaults(config)
	if err != nil {
		return fmt.Errorf("error applying defaults: %w", err)
	}

	// Validate config
	if *config.BaseDir == "" {
		return fmt.Errorf("base directory is required")
	}
	if *config.OutputFile == "" {
		return fmt.Errorf("output file is required")
	}

	// Add output file to excludes
	config.Excludes = append(config.Excludes, *config.OutputFile)

	// Get all files matching include patterns
	var files []string
	for _, pattern := range config.Includes {
		matches, err := doublestar.Glob(os.DirFS(*config.BaseDir), pattern)
		if err != nil {
			return fmt.Errorf("error matching include pattern %s: %w", pattern, err)
		}
		files = append(files, matches...)
	}

	// Remove excluded files
	files = filterExcludedFiles(files, config.Excludes)

	// Sort files based on the order of include patterns
	sort.Slice(files, func(i, j int) bool {
		return getPatternIndex(files[i], config.Includes) < getPatternIndex(files[j], config.Includes)
	})

	// Generate markdown content
	content := generateMarkdownContent(files, config)

	// Write markdown file
	err = os.WriteFile(*config.OutputFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing markdown file: %w", err)
	}

	return nil
}

func filterExcludedFiles(files, excludes []string) []string {
	var result []string
	for _, file := range files {
		excluded := false
		for _, pattern := range excludes {
			match, err := doublestar.Match(pattern, file)
			if err == nil && match {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, file)
		}
	}
	return result
}

func getPatternIndex(file string, patterns []string) int {
	for i, pattern := range patterns {
		match, _ := doublestar.Match(pattern, file)
		if match {
			return i
		}
	}
	return len(patterns)
}

func generateMarkdownContent(files []string, config *MarkdownGeneratorConfig) string {
	var sb strings.Builder

	// Add title
	sb.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", *config.BaseHeaderLevel), *config.Title))

	// Add prepend text
	sb.WriteString(*config.PrependText + "\n\n")

	// Add file contents
	for _, file := range files {
		relPath, _ := filepath.Rel(*config.BaseDir, file)
		sb.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", *config.BaseHeaderLevel+1), relPath))

		content, err := ioutil.ReadFile(filepath.Join(*config.BaseDir, file))
		if err != nil {
			sb.WriteString(fmt.Sprintf("Error reading file: %s\n\n", err))
			continue
		}

		language := inferLanguage(file, string(content))
		sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", language, string(content)))
	}

	return sb.String()
}

func inferLanguage(filename, content string) string {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Analyse(content)
	}
	if lexer == nil {
		return ""
	}
	return lexer.Config().Name
}
