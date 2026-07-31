package projectgraph

// EdgeRelation describes the type of relationship between two nodes.
type EdgeRelation string

const (
	RelationContains  EdgeRelation = "contains"
	RelationDependsOn EdgeRelation = "depends_on"
	RelationUses      EdgeRelation = "uses"
)

// Edge represents a directed relationship between two nodes.
type Edge struct {
	SourceID string
	TargetID string
	Relation EdgeRelation
}
