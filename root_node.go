package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RootNode struct {
	root       *Node
	EditedNode *Node
}

type NodeTheme struct {
	Text     rl.Color
	EditText rl.Color
	BG       rl.Color
	EditBG   rl.Color

	Circle  bool
	Rounded bool

	FontSize    float32
	FontSpacing float32

	Radius float32
	Margin int32
}

func NewRootNode() (inst RootNode) { //root node not got id or added proerly, redo root node where no proper rootnode
	inst.root = NewNode("node", 0, 0, -1)
	inst.root.Flags.IsNested = true

	//Text, EditText, BG, EditBG rl.Color
	theme := NodeTheme{}

	theme.Text = PALETTE["White"]
	theme.EditText = PALETTE["Northern"]
	theme.BG = PALETTE["Northern"]
	theme.EditBG = PALETTE["Circle"]

	theme.Circle = false
	theme.Rounded = true

	theme.FontSize = 16
	theme.FontSpacing = 1

	theme.Radius = 10
	theme.Margin = 5

	inst.root.Theme = theme
	return
}

func (s *RootNode) AddChild() { //change to use newnodeex
	node := NewNode(NODES[ROOT].Name, CAMERA.MousePos.X, CAMERA.MousePos.Y, -1)
	AddLua(node, node.Name)
	node.Parent = NODES[ROOT]
	node.EnableEditing()
	node.Theme = NODES[ROOT].Theme
	NODES[ROOT].Childs = append(NODES[ROOT].Childs, node)
}

func (s *Node) GetPath() string {
	if s.Parent == nil {
		return "/"
	}
	path := s.Parent.GetPath()
	path += s.Name + "/"

	return path
}

func FindOuter(s *Node) *Node {
	if s.Parent == nil {
		return s
	} else if s.Parent.Flags.IsNested {
		return s.Parent
	} else {
		return FindOuter(s.Parent)
	}
}

func (s *RootNode) Update() {

	rl.ClearBackground(PALETTE["White"])

	rl.BeginMode2D(CAMERA.Cam)

	LINE_RENDERER.Draw()

	if rl.IsKeyPressed(rl.KeyEnter) && CAMERA.shift {
		s.AddChild()
	} else if rl.IsKeyPressed(rl.KeyEscape) {
		proot := NODES[ROOT]
		ROOT = FindOuter(NODES[ROOT]).ID

		if s.root.ID != proot.ID {
			proot.EnableEditing()
		}
	}

	for _, v := range NODES[ROOT].Childs {
		v.Update()
	}

	rl.EndMode2D()

	path := "«" + NODES[ROOT].GetPath() + "»" //problems with symbols
	var tPos rl.Vector2
	tSize := rl.MeasureTextEx(FONT, path, 32, 4)
	tPos.X = float32(rl.GetScreenWidth())/2 - tSize.X/2
	tPos.Y = float32(rl.GetScreenHeight()) - tSize.Y
	//fmt.Println(path)
	rl.DrawTextEx(FONT, path, tPos, 32, 4, PALETTE["Northern"])
}
