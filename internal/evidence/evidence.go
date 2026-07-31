package evidence

// Sensitivity defines how sensitive a piece of evidence is.
// This is used by the privacy sanitizer to determine if the evidence can be sent remotely.
type Sensitivity string

const (
	// PublicMetadata means the evidence does not contain personal or proprietary information.
	// Example: version of React in package.json.
	PublicMetadata Sensitivity = "public-metadata"

	// LocalPath means the evidence contains file paths that might be revealing but are relative to the project.
	LocalPath Sensitivity = "local-path"

	// Secret means the evidence must NEVER leave the local machine.
	Secret Sensitivity = "secret"
)

// Evidence captures the exact origin of an assertion made by StackGenome.
type Evidence struct {
	// Kind describes the type of evidence (e.g., "manifest", "rule", "heuristic").
	Kind string `json:"kind"`

	// Path is the relative path from the project root where this evidence was found.
	// It MUST NEVER be an absolute path.
	Path string `json:"path"`

	// Selector is a query or key that points to the exact finding (e.g., "dependencies.react").
	Selector string `json:"selector,omitempty"`

	// Value is the data found (e.g., "18.2.0").
	Value string `json:"value,omitempty"`

	// Sensitivity indicates if this data can be sent remotely or not.
	Sensitivity Sensitivity `json:"sensitivity"`
}
