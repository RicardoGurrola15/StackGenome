package platform

import "stackgenome/internal/detectors"

func init() {
	detectors.Register(&PlatformDetector{})
}
