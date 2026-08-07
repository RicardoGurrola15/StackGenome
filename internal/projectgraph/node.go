package projectgraph

import "stackgenome/internal/evidence"

// NodeType represents the kind of entity.
type NodeType string

const (
	TypeModule         NodeType = "module"
	TypeLanguage       NodeType = "language"
	TypeFramework      NodeType = "framework"
	TypeTool           NodeType = "tool"
	TypePlatform       NodeType = "platform"
	TypeInfrastructure NodeType = "infrastructure"
	TypeDependency     NodeType = "dependency"
	TypeWorkspace      NodeType = "workspace"
	TypeCICD           NodeType = "cicd"
	TypeEditor         NodeType = "editor"
)

// NodeRole represents the role of the entity in the project.
type NodeRole string

const (
	RolePrimary   NodeRole = "primary"
	RoleSatellite NodeRole = "satellite"
	RoleGenerated NodeRole = "generated"
	RoleVendored  NodeRole = "vendored"
	RoleUnknown   NodeRole = "unknown"
)

// NodeScope represents the dependency scope/lifecycle of the node.
type NodeScope string

const (
	ScopeRuntime     NodeScope = "runtime"
	ScopeDevelopment NodeScope = "development"
	ScopeBuild       NodeScope = "build"
	ScopeTest        NodeScope = "test"
	ScopeOptional    NodeScope = "optional"
	ScopePlatform    NodeScope = "platform"
	ScopeUnknown     NodeScope = "unknown"
)

// Node represents an entity in the project graph.
type Node struct {
	ID         string
	Type       NodeType
	Role       NodeRole
	Scope      NodeScope // dependency lifecycle scope
	Name       string
	Version    string // declared constraint (e.g. "^2.5.1")
	Resolved   string // locked/resolved version from lockfile
	Evidences  []evidence.Evidence
	Confidence float64 // 0.0 to 1.0
	Properties map[string]string
}
