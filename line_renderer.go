package main

type LineRenderer struct {
}

func NewLineRenderer() (inst LineRenderer) {
	return
}

type List struct {
	isAlignX bool
	alignX   float32

	isAlignY bool
	alignY   float32

	root *Node
}

func (s *List) Item(i int) *Node {
	return s.GetItem(i, 0, s.root)
}

func (s *List) GetItem(target, depth int, node *Node) (out *Node) {
	if target == depth {
		out = node
	} else {
		out = s.GetItem(target, depth+1, node.listNext)
	}
	return
}

func GetListID() int {
	MAX_LISTID += 1
	LISTS[MAX_LISTID] = &List{}
	return MAX_LISTID
}

type line struct {
	colour, outline string
}

func (s *LineRenderer) Draw() {
	ids := []int{}
	listids := []int{}
	NODES[ROOT].TotalIDs(&ids, &listids)

	cols := make(map[int]Outline)
	outlines := make(map[int]Outline)

	i := 0
	pv := 0
	for _, v := range ids {
		if pv != v {
			i += 1
		}
		cols[v] = Outline{LineColours[i%len(LineColours)], "White"}
		pv = v
	}

	i = 0
	pv = 0
	for _, v := range listids {
		if pv != v {
			i += 1
		}
		outlines[v] = OutlineColours[i%len(OutlineColours)]
		pv = v
	}

	NODES[ROOT].DrawLines(cols, outlines)
}

func (s *Node) DrawLines(cols map[int]Outline, outlines map[int]Outline) {
	if s.ID != ROOT {
		//colour
		var l Outline
		if s.isList {
			l = outlines[s.ListID]
		} else {
			l = cols[s.ID]
		}

		if s.Parent.ID != ROOT {
			DrawPath(FindEdge(s, s.Parent), FindEdge(s.Parent, s), s.invertedLine, PALETTE[l.Inner], s.isList, PALETTE[l.Outer], false)
		}
		// soft links
		if s.slDragging {
			DrawPath(CAMERA.MousePos, FindEdgeToPoint(s, CAMERA.MousePos), s.invertedLine, PALETTE[l.Inner], s.isList, PALETTE[l.Outer], true)
		}
		for _, v := range s.SoftLinks {
			l = cols[v.ID]
			DrawPath(FindEdge(v, s), FindEdge(s, v), s.invertedLine, PALETTE[l.Inner], s.isList, PALETTE[l.Outer], true)
		}
	}
	if !s.Flags.IsNested || s.ID == ROOT {
		for _, v := range s.Childs {
			v.DrawLines(cols, outlines)
		}
	}
}

func (s *Node) TotalIDs(ids *[]int, listids *[]int) {
	if s.isList {
		*listids = append(*listids, s.ListID)
	} else {
		*ids = append(*ids, s.ID)
	}
	for _, v := range s.Childs {
		v.TotalIDs(ids, listids)
	}
}
