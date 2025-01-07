import React, { useCallback, useEffect, useState } from 'react';
import {
  Background,
  ReactFlow,
  useNodesState,
  useEdgesState,
  addEdge,
  MiniMap,
  Controls,
  MarkerType,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

interface Edge {
  From: string;
  To: string;
}

interface Node {
  Name: string;
  Type: string;
}

interface Graph {
  edges: Edge[];
  nodes: Node[];
}

// const hide = (hidden: boolean) => (nodeOrEdge: any) => {
//   return {
//     ...nodeOrEdge,
//     hidden,
//   };
// };

const Flow = () => {
  const [loading, setLoading] = useState(true);
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [hidden, setHidden] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  const onConnect = useCallback(
    // @ts-ignore
    (params: any) => setEdges((els) => addEdge(params, els)),
    []
  );

  async function fetchData() {
    try {
      const response = await fetch('http://localhost:39261/data');
      const data: Graph = await response.json();

      const initialNodes = data.nodes.map((node, idx) => {
        return {
          id: node.Name,
          data: { label: node.Name },
          position: { x: 100, y: 100 * idx },
          // targetPosition: 'left',
          // sourcePosition: 'right',
          // connectable: false,
          // deletable: false,
        };
      });

      const initialEdges = data.edges.map((edge) => {
        return {
          id: edge.From + '-' + edge.To,
          source: edge.From,
          target: edge.To,
          deletable: false,
          reconnectable: false,
          markerEnd: { type: MarkerType.Arrow, width: 24, height: 24 },
        };
      });
      // @ts-ignore
      setNodes(initialNodes);
      // @ts-ignore
      setEdges(initialEdges);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  if (loading) {
    return <div>Loading...</div>;
  }
  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      fitView
      style={{ backgroundColor: '#F7F9FB' }}
    >
      <MiniMap />
      <Controls />
      <div className="isHidden__button">
        <div>
          <label htmlFor="ishidden">
            isHidden
            <input
              id="ishidden"
              type="checkbox"
              checked={hidden}
              onChange={(event) => setHidden(event.target.checked)}
              className="react-flow__ishidden"
            />
          </label>
        </div>
      </div>
      <Background />
    </ReactFlow>
  );
};

const App = () => {
  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <Flow />
    </div>
  );
};

export default App;
