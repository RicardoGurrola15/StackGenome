package projectgraph

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	schemav1 "stackgenome/pkg/schema/v1"
)

// Graph manages nodes and edges in a deterministically ordered way.
type Graph struct {
	nodes map[string]*Node
	edges []Edge
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
		edges: make([]Edge, 0),
	}
}

// AddNode adds a new node to the graph.
// If a node with the same ID already exists, it merges the evidences and returns nil.
func (g *Graph) AddNode(n *Node) error {
	if n.ID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}

	for _, ev := range n.Evidences {
		if filepath.IsAbs(ev.Path) {
			return fmt.Errorf("evidence path must be relative: %s", ev.Path)
		}
	}

	if existing, exists := g.nodes[n.ID]; exists {
		// Merge evidences
		existing.Evidences = append(existing.Evidences, n.Evidences...)
		return nil
	}

	g.nodes[n.ID] = n
	return nil
}

// AddEdge adds a directed edge between two existing nodes.
func (g *Graph) AddEdge(sourceID, targetID string, relation EdgeRelation) error {
	if _, exists := g.nodes[sourceID]; !exists {
		return fmt.Errorf("source node %q does not exist", sourceID)
	}
	if _, exists := g.nodes[targetID]; !exists {
		return fmt.Errorf("target node %q does not exist", targetID)
	}

	g.edges = append(g.edges, Edge{
		SourceID: sourceID,
		TargetID: targetID,
		Relation: relation,
	})
	return nil
}

// ToDTO converts the graph into the v1 versioned schema struct, maintaining deterministic order.
func (g *Graph) ToDTO() *schemav1.ProjectGraphDTO {
	dto := &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes:   make([]schemav1.NodeDTO, 0, len(g.nodes)),
		Edges:   make([]schemav1.EdgeDTO, 0, len(g.edges)),
	}

	// Sort node IDs to guarantee determinism
	var nodeIDs []string
	for id := range g.nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	for _, id := range nodeIDs {
		n := g.nodes[id]

		// Map properties stably (though Go JSON marshals maps in sorted order, the struct maps themselves are fine)
		var props map[string]string
		if len(n.Properties) > 0 {
			props = make(map[string]string)
			for k, v := range n.Properties {
				props[k] = v
			}
		}

		dto.Nodes = append(dto.Nodes, schemav1.NodeDTO{
			ID:         n.ID,
			Type:       string(n.Type),
			Role:       string(n.Role),
			Name:       n.Name,
			Version:    n.Version,
			Evidences:  n.Evidences,
			Confidence: n.Confidence,
			Properties: props,
		})
	}

	// Sort edges deterministically
	sortedEdges := make([]Edge, len(g.edges))
	copy(sortedEdges, g.edges)
	sort.Slice(sortedEdges, func(i, j int) bool {
		cmpSource := strings.Compare(sortedEdges[i].SourceID, sortedEdges[j].SourceID)
		if cmpSource != 0 {
			return cmpSource < 0
		}
		cmpTarget := strings.Compare(sortedEdges[i].TargetID, sortedEdges[j].TargetID)
		if cmpTarget != 0 {
			return cmpTarget < 0
		}
		return strings.Compare(string(sortedEdges[i].Relation), string(sortedEdges[j].Relation)) < 0
	})

	for _, e := range sortedEdges {
		dto.Edges = append(dto.Edges, schemav1.EdgeDTO{
			SourceID: e.SourceID,
			TargetID: e.TargetID,
			Relation: string(e.Relation),
		})
	}

	return dto
}
