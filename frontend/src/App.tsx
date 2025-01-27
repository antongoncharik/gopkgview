import { useEffect, useState } from 'react';
import {
  Background,
  ReactFlow,
  useNodesState,
  useEdgesState,
  MiniMap,
  Controls,
  MarkerType,
  useReactFlow,
  ReactFlowProvider,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import ELK from 'elkjs/lib/elk.bundled.js';

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

interface InitialNode {
  id: string;
  data: {
    label: string;
  };
  position: {
    x: number;
    y: number;
  };
  targetPosition: string;
  sourcePosition: string;
  connectable: boolean;
  deletable: boolean;
}

interface InitialEdge {
  id: string;
  source: string;
  target: string;
  deletable: boolean;
  reconnectable: boolean;
  markerEnd: {
    type: MarkerType;
    width: number;
    height: number;
  };
}

const elk = new ELK();

const elkOptions = {
  'elk.algorithm': 'layered',
  'elk.direction': 'RIGHT',
  'elk.edgeRouting': 'SPLINES',
  'elk.layered.edgeRouting.splines.mode': 'CONSERVATIVE',
  'elk.layered.layering.strategy': 'NETWORK_SIMPLEX',
  'elk.layered.spacing.nodeNodeBetweenLayers': 200,
  'elk.spacing.nodeNodeBetweenLayers': 100,
  'elk.spacing.nodeNode': 30,
  'elk.spacing.edgeEdge': 20,
  'elk.spacing.edgeNode': 50,
};

const Flow = () => {
  const [loading, setLoading] = useState(true);
  const [initialNodes, setInitialNodes] = useState<InitialNode[]>([]);
  const [initialEdges, setInitialEdges] = useState<InitialEdge[]>([]);
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const { fitView } = useReactFlow();

  useEffect(() => {
    (async () => {
      const data = await fetchData();
      if (!data) {
        return;
      }

      const { initialNodes, initialEdges } = data;

      setInitialNodes(initialNodes);
      setInitialEdges(initialEdges);
      setNodes(initialNodes as never[]);
      setEdges(initialEdges as never[]);
    })();
  }, []);

  useEffect(() => {
    (async () => {
      const layout = await getLayoutedElements(initialNodes, initialEdges);

      setNodes(layout.nodes);
      setEdges(layout.edges);

      setTimeout(() => {
        window.requestAnimationFrame(() => fitView());
      }, 20);
    })();
  }, [initialNodes, initialEdges]);

  const getLayoutedElements = async (nodes: any, edges: any) => {
    const elkGraph = {
      id: 'root',
      layoutOptions: elkOptions,
      edges,
      children: nodes,
    };

    try {
      // @ts-ignore
      const { children, edges: layoutedEdges } = await elk.layout(elkGraph);
      // @ts-ignore
      const layoutedNodes = children.map(({ x, y, ...node }) => ({
        ...node,
        position: { x, y },
      }));

      return { nodes: layoutedNodes, edges: layoutedEdges };
    } catch (err) {
      console.error('ELK layout failed:', err);
      return { nodes, edges };
    }
  };

  const fetchData = async (): Promise<
    | {
        initialNodes: InitialNode[];
        initialEdges: InitialEdge[];
      }
    | undefined
  > => {
    try {
      const response = await fetch('/data');
      const data: Graph = await response.json();

      const initialNodes = data.nodes.map((node) => {
        return {
          id: node.Name,
          data: { label: node.Name },
          position: { x: 0, y: 0 },
          targetPosition: 'left',
          sourcePosition: 'right',
          connectable: false,
          deletable: false,
          width: Math.max(80, node.Name.length * 8),
          height: 40,
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

      return { initialNodes, initialEdges };
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }
  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      fitView
      style={{ backgroundColor: '#F7F9FB' }}
      minZoom={0.1}
      maxZoom={3}
      fitViewOptions={{ padding: 0.5 }}
    >
      <MiniMap />
      <Controls />
      <Background />
    </ReactFlow>
  );
};

const App = () => {
  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlowProvider>
        <Flow />
      </ReactFlowProvider>
    </div>
  );
};

export default App;
