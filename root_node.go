package main
import (
  "github.com/gen2brain/raylib-go/raylib"
)

type RootNode struct {
  node *Node
  EditedNode *Node
}

type NodeTheme struct {
  Text rl.Color
  EditText rl.Color
  BG rl.Color
  EditBG rl.Color

  Circle bool
  Rounded bool

  FontSize float32
  FontSpacing float32

  Radius float32
  Margin int32
}

func NewRootNode() (inst RootNode) {//root node not got id or added proerly, redo root node where no proper rootnode
  inst.node = NewNode("node", 0, 0, -1)
  inst.node.ID = -1

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

  inst.node.Theme = theme
  return
}

func (s *RootNode) AddChild() {//change to use newnodeex
  node := NewNode(s.node.Name, CAMERA.MousePos.X, CAMERA.MousePos.Y, -1)
  node.GetID()
  ParseCode(node.ID, node.Name)
  node.Parent = s.node
  node.EnableEditing()
  node.Theme = ROOT.Theme
  s.node.Childs = append(s.node.Childs, node)
}

func (s *RootNode) Update() {
  if rl.IsKeyPressed(rl.KeyEnter) && CAMERA.shift{
    s.AddChild()
  }

  for _, v := range s.node.Childs {
    v.Update()
  }
}
