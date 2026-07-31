package cicd

import "stackgenome/internal/detectors"

func init() {
	detectors.Register(&CICDDetector{})
}
