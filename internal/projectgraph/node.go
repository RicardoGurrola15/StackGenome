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

// Node represents an entity in the project graph.
type Node struct {
	ID         string
	Type       NodeType
	Role       NodeRole
	Name       string
	Version    string
	Evidences  []evidence.Evidence
	Confidence float64 // 0.0 to 1.0
	Properties map[string]string
}
