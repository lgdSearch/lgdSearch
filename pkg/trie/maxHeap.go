package trie

func (node1 *heapNode) maxHeapCompare(node2 *heapNode) bool {
	return node1.deep > node2.deep
}

// MaxHeap is a max-heap of *heapNode.
type MaxHeap []*heapNode

func (h MaxHeap) Len() int { return len(h) }

func (h MaxHeap) Less(i, j int) bool {
	return h[i].maxHeapCompare(h[j])
}

func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(*heapNode))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
