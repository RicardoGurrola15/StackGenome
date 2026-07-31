package editor

import "stackgenome/internal/detectors"

func init() {
	detectors.Register(&EditorDetector{})
}
