package language

import (
	"stackgenome/internal/detectors"
)

func init() {
	detectors.Register(&GoDetector{})
	detectors.Register(&NodeDetector{})
	detectors.Register(&PythonDetector{})
	detectors.Register(&RustDetector{})
	detectors.Register(&JVMDetector{})
	detectors.Register(&SwiftDetector{})
	detectors.Register(&DotNetDetector{})
	detectors.Register(&DartDetector{})
	detectors.Register(&PHPDetector{})
	detectors.Register(&RubyDetector{})
	detectors.Register(&CppDetector{})
	detectors.Register(&ExtensionDetector{})
}
