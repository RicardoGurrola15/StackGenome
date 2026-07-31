package language

import (
	"path/filepath"
	"strings"

	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

var extensionToLanguage = map[string]string{
	".go":    "Go",
	".js":    "JavaScript",
	".ts":    "TypeScript",
	".py":    "Python",
	".rs":    "Rust",
	".java":  "Java",
	".c":     "C",
	".cpp":   "C++",
	".h":     "C/C++ Header",
	".cs":    "C#",
	".rb":    "Ruby",
	".php":   "PHP",
	".swift": "Swift",
	".kt":    "Kotlin",
	".sh":    "Shell",
	".bash":  "Shell",
	".html":  "HTML",
	".css":   "CSS",
}

// ExtensionDetector infers programming languages based on file extensions.
type ExtensionDetector struct{}

func (d *ExtensionDetector) Handles(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	_, exists := extensionToLanguage[ext]
	return exists
}

func (d *ExtensionDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	ext := strings.ToLower(filepath.Ext(relPath))
	langName, exists := extensionToLanguage[ext]
	if !exists {
		return nil, nil, nil
	}

	node := &projectgraph.Node{
		ID:         "lang_ext_" + strings.ToLower(langName),
		Type:       projectgraph.TypeLanguage,
		Name:       langName,
		Confidence: 0.3, // Low confidence since it's just an extension
		Evidences: []evidence.Evidence{
			{
				Kind:        "extension",
				Path:        relPath,
				Value:       ext,
				Sensitivity: evidence.LocalPath,
			},
		},
	}

	return []*projectgraph.Node{node}, nil, nil
}
