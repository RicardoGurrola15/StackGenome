package infra

import "stackgenome/internal/detectors"

func init() {
	detectors.Register(&InfraDetector{})
}
