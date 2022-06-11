package trie

type heapNode struct {
	node *Node
	deep int32 // deep represents the relationship with the pattern string
	size int32
	son  rune
}

func (node1 *heapNode) compare(node2 *heapNode) bool {
	if node1.deep != node2.deep {
		return node1.deep < node2.deep
	} else {
		return node1.node.count < node2.node.count
	}
}

// Heap is a min-heap of *heapNode.
type Heap []*heapNode

func (h Heap) Len() int { return len(h) }

func (h Heap) Less(i, j int) bool {
	return h[i].compare(h[j])
}

func (h Heap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(*heapNode))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
