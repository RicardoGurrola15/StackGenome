'use client'
import { useEffect, useRef, useState } from 'react'
import * as d3Force from 'd3-force'
import { select } from 'd3-selection'
import * as d3Drag from 'd3-drag'
import * as d3Zoom from 'd3-zoom'
import type { NodeDTO, EdgeDTO } from '@/types/stackgenome'

interface GraphNode extends d3Force.SimulationNodeDatum {
  id: string
  type: string
  name: string
  version?: string
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
  dependency: '#6e7681',
}

export default function StackGraph({ nodes, edges }: { nodes: NodeDTO[], edges: EdgeDTO[] }) {
  const svgRef = useRef<SVGSVGElement>(null)
  const [showDeps, setShowDeps] = useState(false)

  useEffect(() => {
    if (!svgRef.current || nodes.length === 0) return

    const width = svgRef.current.clientWidth || 800
    const height = svgRef.current.clientHeight || 500

    const visibleNodes = showDeps ? nodes : nodes.filter(n => n.type !== 'dependency')
    const visibleIds = new Set(visibleNodes.map(n => n.id))
    const visibleEdges = edges.filter(e => visibleIds.has(e.source_id) && visibleIds.has(e.target_id))

    const graphNodes: GraphNode[] = visibleNodes.map(n => ({ ...n }))
    const graphLinks: GraphLink[] = visibleEdges.map(e => ({ ...e, source: e.source_id, target: e.target_id }))

    const simulation = d3Force.forceSimulation<GraphNode>(graphNodes)
      .force('link', d3Force.forceLink<GraphNode, GraphLink>(graphLinks).id(d => d.id).distance(120))
      .force('charge', d3Force.forceManyBody().strength(-400))
      .force('center', d3Force.forceCenter(width / 2, height / 2))
      .force('collision', d3Force.forceCollide().radius(40))

    const svg = select(svgRef.current)
    svg.selectAll('*').remove() // Clear

    const g = svg.append('g') // Zoomable container

    const zoom = d3Zoom.zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.1, 4])
      .on('zoom', (event) => {
        g.attr('transform', event.transform)
      })

    svg.call(zoom)

    // Arrows
    svg.append('defs').append('marker')
      .attr('id', 'arrow')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 22)
      .attr('refY', 0)
      .attr('markerWidth', 6)
      .attr('markerHeight', 6)
      .attr('orient', 'auto')
      .append('path')
      .attr('d', 'M0,-5L10,0L0,5')
      .attr('fill', 'rgba(255,255,255,0.2)')

    // Links
    const link = g.append('g')
      .selectAll('line')
      .data(graphLinks)
      .join('line')
      .attr('stroke', 'rgba(255,255,255,0.1)')
      .attr('stroke-width', 1.5)
      .attr('marker-end', 'url(#arrow)')

    // Nodes
    const node = g.append('g')
      .selectAll('g')
      .data(graphNodes)
      .join('g')
      .call(d3Drag.drag<SVGGElement, GraphNode>()
        .on('start', (event, d) => {
          if (!event.active) simulation.alphaTarget(0.3).restart()
          d.fx = d.x
          d.fy = d.y
        })
        .on('drag', (event, d) => {
          d.fx = event.x
          d.fy = event.y
        })
        .on('end', (event, d) => {
          if (!event.active) simulation.alphaTarget(0)
          d.fx = null
          d.fy = null
        })
      )
      
    // Circles with glow
    node.append('circle')
      .attr('r', d => d.type === 'primary' ? 16 : 12)
      .attr('fill', d => TYPE_COLORS[d.type] || '#6e7681')
      .attr('stroke', '#0d1117')
      .attr('stroke-width', 2)
      .style('filter', d => `drop-shadow(0 0 8px ${TYPE_COLORS[d.type] || '#6e7681'}40)`)
      .style('cursor', 'grab')

    // Labels
    node.append('text')
      .text(d => d.name)
      .attr('x', d => (d.type === 'primary' ? 20 : 16))
      .attr('y', 4)
      .attr('fill', '#e6edf3')
      .attr('font-size', '13px')
      .attr('font-weight', d => d.type === 'primary' ? '600' : '400')
      .attr('font-family', 'var(--font-sans)')
      .style('text-shadow', '0 2px 4px rgba(0,0,0,0.8)')
      .style('pointer-events', 'none')

    // Tooltip
    node.append('title')
      .text(d => `${d.name} (${d.type})${d.version ? `\nv: ${d.version}` : ''}`)

    simulation.on('tick', () => {
      link
        .attr('x1', d => (d.source as GraphNode).x!)
        .attr('y1', d => (d.source as GraphNode).y!)
        .attr('x2', d => (d.target as GraphNode).x!)
        .attr('y2', d => (d.target as GraphNode).y!)

      node.attr('transform', d => `translate(${d.x},${d.y})`)
    })

    return () => {
      simulation.stop()
    }
  }, [nodes, edges, showDeps])

  const depCount = nodes.filter(n => n.type === 'dependency').length

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-3)', width: '100%' }}>
      {depCount > 0 && (
        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <label className="badge badge--gray" style={{ cursor: 'pointer', padding: 'var(--space-2) var(--space-4)' }}>
            <input 
              type="checkbox" 
              checked={showDeps} 
              onChange={e => setShowDeps(e.target.checked)} 
              style={{ marginRight: 'var(--space-2)' }} 
            />
            Mostrar dependencias de runtime ({depCount})
          </label>
        </div>
      )}
      <div style={{ 
        width: '100%', 
        height: '600px', 
        background: 'var(--bg-secondary)', 
        backgroundImage: 'radial-gradient(rgba(255,255,255,0.05) 1px, transparent 1px)',
        backgroundSize: '20px 20px',
        borderRadius: 'var(--radius-xl)', 
        border: '1px solid var(--border)', 
        overflow: 'hidden',
        boxShadow: 'inset 0 0 60px rgba(0,0,0,0.5)'
      }}>
        <svg ref={svgRef} style={{ width: '100%', height: '100%' }} />
      </div>
      <p style={{ fontSize: 'var(--text-xs)', color: 'var(--text-muted)', textAlign: 'center' }}>
        Usa la rueda del ratón para hacer zoom. Arrastra los nodos para reorganizarlos.
      </p>
    </div>
  )
}
