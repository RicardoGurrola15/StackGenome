'use client'
import { useEffect, useRef } from 'react'
import * as d3Force from 'd3-force'
import { select } from 'd3-selection'
import type { NodeDTO, EdgeDTO } from '@/types/stackgenome'

interface GraphNode extends d3Force.SimulationNodeDatum {
  id: string
  type: string
  name: string
}

interface GraphLink extends d3Force.SimulationLinkDatum<GraphNode> {
  source: string | GraphNode
  target: string | GraphNode
  relation: string
}

const TYPE_COLORS: Record<string, string> = {
  primary: '#3fb950',
  satellite: '#238636',
  framework: '#58a6ff',
  database: '#bc8cff',
  tooling: '#f0883e',
  infrastructure: '#8b949e',
}

export default function StackGraph({ nodes, edges }: { nodes: NodeDTO[], edges: EdgeDTO[] }) {
  const svgRef = useRef<SVGSVGElement>(null)

  useEffect(() => {
    if (!svgRef.current || nodes.length === 0) return

    const width = svgRef.current.clientWidth || 800
    const height = svgRef.current.clientHeight || 500

    // Prepare data
    const graphNodes: GraphNode[] = nodes.map(n => ({ ...n }))
    const graphLinks: GraphLink[] = edges.map(e => ({ ...e, source: e.source_id, target: e.target_id }))

    // Setup simulation
    const simulation = d3Force.forceSimulation<GraphNode>(graphNodes)
      .force('link', d3Force.forceLink<GraphNode, GraphLink>(graphLinks).id(d => d.id).distance(100))
      .force('charge', d3Force.forceManyBody().strength(-300))
      .force('center', d3Force.forceCenter(width / 2, height / 2))
      .force('collision', d3Force.forceCollide().radius(30))

    const svg = select(svgRef.current)
    svg.selectAll('*').remove() // Clear previous render

    // Define arrow markers
    svg.append('defs').append('marker')
      .attr('id', 'arrow')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 20)
      .attr('refY', 0)
      .attr('markerWidth', 6)
      .attr('markerHeight', 6)
      .attr('orient', 'auto')
      .append('path')
      .attr('d', 'M0,-5L10,0L0,5')
      .attr('fill', 'var(--border)')

    // Draw links
    const link = svg.append('g')
      .selectAll('line')
      .data(graphLinks)
      .join('line')
      .attr('stroke', 'var(--border)')
      .attr('stroke-width', 1.5)
      .attr('marker-end', 'url(#arrow)')

    // Draw nodes
    const node = svg.append('g')
      .selectAll('g')
      .data(graphNodes)
      .join('g')
      
    // Node circles
    node.append('circle')
      .attr('r', 12)
      .attr('fill', d => TYPE_COLORS[d.type] || 'var(--text-secondary)')
      .attr('stroke', 'var(--bg-primary)')
      .attr('stroke-width', 2)

    // Node labels
    node.append('text')
      .text(d => d.name)
      .attr('x', 16)
      .attr('y', 4)
      .attr('fill', 'var(--text-primary)')
      .attr('font-size', '12px')
      .attr('font-family', 'var(--font-sans)')

    // Simulation tick updates
    simulation.on('tick', () => {
      link
        .attr('x1', d => (d.source as GraphNode).x!)
        .attr('y1', d => (d.source as GraphNode).y!)
        .attr('x2', d => (d.target as GraphNode).x!)
        .attr('y2', d => (d.target as GraphNode).y!)

      node
        .attr('transform', d => `translate(${d.x},${d.y})`)
    })

    return () => {
      simulation.stop()
    }
  }, [nodes, edges])

  return (
    <div style={{ width: '100%', height: '500px', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border)', overflow: 'hidden' }}>
      <svg ref={svgRef} style={{ width: '100%', height: '100%' }} />
    </div>
  )
}
