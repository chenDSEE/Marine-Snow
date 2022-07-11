package framework

import (
	"errors"
	"fmt"
	"strings"
)

type rTree struct {
	name string
	root *node

	// TODO: add lock to support hot update route tree
}

func newRouteTree(name string) *rTree {
	return &rTree{
		name: name,
		root: newNode(),
	}
}

func (rtree *rTree) addRoute(url string, handlers []handlerFuncEntry) error {
	/* if uri existed, return error */
	if n := rtree.root.matchNode(url); n != nil {
		return errors.New("Route conflict: " + url + " with [" + n.fullUrl + "]")
	}

	segments := strings.Split(url, "/") // "" as root node

	/* find the right place to create and insert current node for each segment */
	currNode := rtree.root

nextRound:
	for index, segment := range segments { // each round for a segment place search
		// make uri case-insensitive
		if isWildSegment(segment) == false {
			segment = strings.ToUpper(segment)
		}

		isEnd := (index == len(segments)-1)

		// try to attach current segment to an existed node
		children := currNode.getValidChildNode(segment)
		if len(children) > 0 {
			for _, child := range children {
				if child.segment == segment {
					currNode = child // for next round
					continue nextRound
				}
			}
		}

		// can not attach current segment into any child node, create a new node
		segmentNode := newNode()
		segmentNode.segment = segment
		if isEnd == true {
			segmentNode.isEnd = true
			segmentNode.fullUrl = url
			segmentNode.handlerEntryList = handlers
		}

		currNode.children = append(currNode.children, segmentNode) // register into parent node
		currNode = segmentNode                                     // for next round
	}

	return nil
}

type routeEntry struct {
	url         string
	handlerName string
}

func (rtree *rTree) printRouteTree() {
	// perform a DFS
	entries := travelRouteTable(rtree.root, make([]routeEntry, 0))

	fmt.Printf("===== %s tire tree dump =====\n", rtree.name)
	for i := 0; i < len(entries); i++ {
		fmt.Printf("+ [%s] --> [%s]\n", entries[i].url, entries[i].handlerName)
	}
}

func travelRouteTable(n *node, entries []routeEntry) []routeEntry {
	// perform a DFS search
	if n.isEnd == true {
		entry := routeEntry{
			url:         n.fullUrl,
			handlerName: n.handlerEntryList[len(n.handlerEntryList)-1].funName,
		}
		entries = append(entries, entry)
	}

	for i := 0; i < len(n.children); i++ {
		entries = travelRouteTable(n.children[i], entries)
	}

	return entries
}

func (rtree *rTree) FindHandlerEntryList(url string) ([]handlerFuncEntry, string) {
	matchNode := rtree.root.matchNode(url)
	if matchNode == nil {
		return nil, ""
	}

	return matchNode.handlerEntryList, matchNode.fullUrl
}

type node struct {
	segment          string
	fullUrl          string
	isEnd            bool
	handlerEntryList []handlerFuncEntry
	children         []*node
}

func newNode() *node {
	return &node{
		isEnd:    false,
		segment:  "", // as default value for root node
		children: []*node{},
	}
}

func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// get the valid children node of current n by segment
func (n *node) getValidChildNode(segment string) []*node {
	if len(n.children) == 0 {
		// current node has not children node
		return nil
	}

	if isWildSegment(segment) == true {
		// only used when addRoute()
		return n.children
	}

	// try to find valid child node by segment and childrenNode.segment
	nodes := make([]*node, 0, len(n.children))
	for _, child := range n.children {
		if isWildSegment(child.segment) == true || segment == child.segment {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

// match node by URL, perform a BFS search
func (n *node) matchNode(url string) *node {
	// only get current segment from url
	segments := strings.SplitN(url, "/", 2)
	segment := segments[0]

	if isWildSegment(segment) == false { // for addRoute
		segment = strings.ToUpper(segment)
	}

	/* check current level node segment */
	children := n.getValidChildNode(segment)
	if children == nil || len(children) == 0 {
		// can not match any node by this url segment[0]
		return nil
	}

	if len(segments) == 1 { // url segment is end, check node
		for _, child := range children {
			if child.isEnd == true {
				// already check segment == child.segment in node.getValidChildNode()
				return child // match a node by this url
			}
		}

		return nil // can not match any node by this url
	}

	// further BFS
	for _, child := range children {
		matchNode := child.matchNode(segments[1])
		if matchNode != nil {
			return matchNode // match success and return
		}
	}

	return nil // can not match any node by this url
}
