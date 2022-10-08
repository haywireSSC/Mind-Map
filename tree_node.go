package main

import (
	"fmt"
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/slices"
)

type NodeEditor struct {
	IsEditing bool
	Pos       int
	Text      string
}

type Node struct {
	Rect   rl.Rectangle
	Pos    rl.Vector2
	Center rl.Vector2
	Name   string

	Childs []*Node
	Parent *Node

	isList       bool
	listNext     *Node
	ListID       int
	invertedLine bool
	hasImg       bool
	codeCooldown Cooldown
	ID           int

	Tex rl.Texture2D

	E     NodeEditor
	Theme NodeTheme

	dragging   bool
	dragOffset rl.Vector2
	startDrag  rl.Vector2

	Flags NodeFlags
}

func NewNode(name string, x, y float32, id int) (inst *Node) {
	inst = &Node{}
	inst.Pos.X = x
	inst.Pos.Y = y
	inst.Name = name
	inst.E.Text = name
	//inst.Flags = NodeFlags{true, true}
	inst.codeCooldown = Cooldown{0, 0}
	if id == -1 {
		MAX_ID += 1
		NODES[MAX_ID] = inst
		inst.ID = MAX_ID
	} else {
		inst.ID = id
		if id > MAX_ID {
			MAX_ID = id
		}

		NODES[id] = inst
	}
	AddLua(inst, inst.Name)
	return inst
}

func (s *Node) DoAlignment() {
	if s.isList {
		l := LISTS[s.ListID]

		isRoot := s == l.root

		if l.isAlignX {
			if isRoot {
				if s.Parent.ID != ROOT {
					l.alignX = s.Parent.Pos.X
				} else {
					l.alignX = s.Pos.X
				}
			}
			s.Pos.X = l.alignX
		} else if l.isAlignY {
			if isRoot {
				if s.Parent.ID != ROOT {
					l.alignY = s.Parent.Pos.Y
				} else {
					l.alignY = s.Pos.Y
				}
			}
			s.Pos.Y = l.alignY
		}

	}
}

func (s *Node) ToggleAlign() {
	l := LISTS[s.ListID]

	if l.isAlignX || l.isAlignY {
		l.isAlignX = false
		l.isAlignY = false
		return
	}

	xdiff := math.Abs(float64(s.Pos.X - s.Parent.Pos.X))
	ydiff := math.Abs(float64(s.Pos.Y - s.Parent.Pos.Y))

	if xdiff < ydiff {
		l.isAlignX = true
		l.alignX = l.root.Pos.X
	} else {
		l.isAlignY = true
		l.alignY = l.root.Pos.Y
	}
}

func NewNodeEx() *Node {
	node := NewNode(NODES[ROOT].Name, 0, 0, -1)
	node.Theme = NODES[ROOT].Theme
	//node.CenterOn(CAMERA.MousePos)
	//node.StartDrag()
	return node
}

func (s *Node) AddChild() *Node {
	fmt.Println("added child")
	node := NewNodeEx()
	node.Parent = s
	node.EnableEditing()
	node.CenterOn(CAMERA.MousePos)
	node.StartDrag()

	node.EnableEditing()
	node.invertedLine = !s.invertedLine
	s.Childs = append(s.Childs, node)
	return node
}
func (s *Node) GetChildByName(name string) *Node {
	for _, v := range s.Childs {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (s *Node) CenterOn(center rl.Vector2) {
	s.UpdateRect()
	s.Pos = center
	s.Pos.X -= s.Rect.Width / 2
	s.Pos.Y -= s.Rect.Height / 2
}

func (s *Node) MakeList(id int) {
	//recurse back to root
	s.ListID = id
	s.isList = true
	if s.Parent.ID != ROOT && !s.Parent.isList {
		s.Parent.listNext = s
		s.Parent.MakeList(id)
	} else {
		LISTS[s.ListID].root = s
	}
}

func (s *Node) SplitList(id int, first bool) {
	//recurse forward on same id replacing id

	for _, v := range s.Childs {
		if v.isList && v.ListID == s.ListID {
			if s.Parent.listNext == nil {
				LISTS[id].root = s
			}
			v.SplitList(id, false)
			if !first {
				s.listNext = v
			}
		}
	}
	if !first {
		s.ListID = id
	}
}

func (s *Node) AddToList() (id int) {
	if s.isList {
		if len(s.Childs) == 1 {
			id = s.ListID
		} else {
			//split off branch and return new id
			id = GetListID()
			s.SplitList(GetListID(), true)
		}
	} else {
		if s.Parent.ID != ROOT {
			id = s.Parent.AddToList()
			//been branched or not
			if s.Parent.listNext == nil {
				LISTS[id].root = s
			} else {
				s.Parent.listNext = s
			}
			s.ListID = id
			s.isList = true
		}
	}
	return
}

func (s *Node) EnableEditing() {
	HANDLER.EditedNode = s
	s.E.Pos = len(s.E.Text)
	s.E.IsEditing = true
}

func (s *Node) DisableEditing() {
	if s.E.IsEditing {
		AddLua(s, s.E.Text)
		s.Name = RunLua(s.ID)
		s.E.IsEditing = false
	}
}

func (s *Node) UpdateRect() {
	if s.Theme.Circle {
		s.Rect = rl.NewRectangle(s.Pos.X, s.Pos.Y, 0, 0)
		s.Center = s.Pos
	} else {
		width := float32(s.Theme.Margin * 2)
		height := float32(s.Theme.Margin * 2)
		if s.hasImg {
			width += float32(s.Tex.Width)
			height += float32(s.Tex.Height)
		} else {
			v := rl.MeasureTextEx(FONT, s.Name, s.Theme.FontSize, s.Theme.FontSpacing)
			width += v.X
			height += v.Y
		}
		s.Rect = rl.NewRectangle(s.Pos.X, s.Pos.Y, width, height)
		s.Center.X = s.Pos.X + s.Rect.Width/2
		s.Center.Y = s.Pos.Y + s.Rect.Height/2
	}
}

func (s *Node) StartDrag() {
	s.dragging = true
	s.dragOffset = rl.Vector2Subtract(s.Pos, CAMERA.MousePos)
	s.startDrag = s.Pos
}

func ManhatLength(pos rl.Vector2) float64 {
	return math.Abs(float64(pos.X)) + math.Abs(float64(pos.Y))
}

func (s *Node) Update() {
	if HANDLER.EditedNode != s {
		s.DisableEditing()
	}
	//dragging
	var mouseHover bool
	if s.Theme.Circle {
		mouseHover = rl.CheckCollisionPointCircle(CAMERA.MousePos, s.Center, s.Theme.Radius)
	} else {
		mouseHover = rl.CheckCollisionPointRec(CAMERA.MousePos, s.Rect)
	}
	if mouseHover && rl.IsMouseButtonPressed(0) {
		s.StartDrag()
	}
	if s.dragging {
		s.Pos = rl.Vector2Add(CAMERA.MousePos, s.dragOffset)
		s.UpdateRect()
	}
	if (rl.IsMouseButtonReleased(0) || rl.IsMouseButtonReleased(1)) && s.dragging {
		s.dragging = false
	}

	s.DoAlignment()
	//text edit
	if s.E.IsEditing {
		s.Editor()
	}

	//hover actions
	if mouseHover {
		//ctrl := rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)
		shift := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

		if rl.IsMouseButtonReleased(0) && ManhatLength(rl.Vector2Subtract(s.Pos, s.startDrag)) < 10 {
			if s.E.IsEditing {
				HANDLER.EditedNode = nil
				s.DisableEditing()
			} else {
				s.EnableEditing()
			}

		} else if rl.IsMouseButtonPressed(1) && !s.dragging {
			if !s.Flags.IsNested || s.ID == ROOT {
				child := s.AddChild()
				if shift {
					child.AddToList()
				}
			}

		}
	}

	if s.codeCooldown.tick() && !s.E.IsEditing {
		s.Name = RunLua(s.ID)
	}
	if !s.Flags.IsNested || s.ID == ROOT {
		for _, v := range s.Childs {
			v.Update()
		}
	}
	s.Draw()
}

//editor functions it uses
func (s *Node) AddLetter(c string) {
	s.E.Text = s.E.Text[:s.E.Pos] + c + s.E.Text[s.E.Pos:]
	s.E.Pos += 1
}

func (s *Node) MoveBackWord() {
	for i := s.E.Pos - 1; i > 0; i-- {
		if s.E.Text[i] != ' ' && s.E.Text[i-1] == ' ' {
			s.E.Pos = i
			return
		}
	}
	s.E.Pos = 0
}

func (s *Node) MoveFwdWord() {
	for i := s.E.Pos + 1; i < len(s.E.Text); i++ {
		if s.E.Text[i-1] != ' ' && s.E.Text[i] == ' ' {
			s.E.Pos = i
			return
		}
	}
	s.E.Pos = len(s.E.Text)
}

func (s *Node) DeleteWordBack() {
	for i := s.E.Pos - 1; i > 0; i-- {
		if s.E.Text[i] != ' ' && s.E.Text[i-1] == ' ' {
			s.E.Text = s.E.Text[:i] + s.E.Text[s.E.Pos:]
			s.E.Pos = i
			return
		}
	}
	s.E.Text = s.E.Text[s.E.Pos:]
	s.E.Pos = 0
}

func (s *Node) DeleteWordFwd() {
	for i := s.E.Pos + 1; i < len(s.E.Text); i++ {
		if s.E.Text[i-1] != ' ' && s.E.Text[i] == ' ' {
			s.E.Text = s.E.Text[:s.E.Pos] + s.E.Text[i:]
			return
		}
	}
	s.E.Text = s.E.Text[:s.E.Pos]
}

func (s *Node) MoveUp() {
	before := strings.LastIndex(s.E.Text[:s.E.Pos], "\n")
	if before == -1 {
		s.E.Pos = 0
		return
	}
	diff := s.E.Pos - before
	next := strings.LastIndex(s.E.Text[:before], "\n") // + before
	s.E.Pos = next + diff
}

func (s *Node) MoveDown() {
	before := strings.LastIndex(s.E.Text[:s.E.Pos], "\n")
	diff := s.E.Pos - before
	next := strings.Index(s.E.Text[s.E.Pos:], "\n")
	if next == -1 {
		s.E.Pos = len(s.E.Text)
		return
	}
	s.E.Pos += next + diff
}

type nestedIdx struct {
	idx  int
	iIdx int
}

func (s *Node) Editor() {
	key := rl.GetCharPressed()
	ctrl := CAMERA.ctrl
	shift := CAMERA.shift
	//pasting
	if rl.IsKeyPressed(rl.KeyV) && ctrl {
		if text := TextFromClipboard(); text == "" {
			rl.UnloadTexture(s.Tex)
			s.Tex = TextureFromClipboard()
			s.hasImg = true
			s.Theme.Rounded = false
			s.Name = "image"
		} else {
			s.E.Text = s.E.Text[:s.E.Pos] + text + s.E.Text[s.E.Pos:]
			s.E.Pos += len(text)
		}
	}

	if rl.IsKeyPressed(rl.KeyEnter) && !shift {
		if ctrl {
			s.DisableEditing()
			ROOT = s.ID
		} else {
			s.AddLetter("\n")
		}
	}
	//typing
	if key != 0 && !ctrl {
		s.AddLetter(string(key))
		//backspacing
	} else if rl.IsKeyPressed(rl.KeyBackspace) {
		if shift {
			s.hasImg = false
			s.Theme.Rounded = NODES[ROOT].Theme.Rounded
			s.E.Text = ""
			s.E.Pos = 0
		} else if s.E.Pos-1 >= 0 {
			if ctrl {
				s.DeleteWordBack()
			} else {
				s.E.Text = s.E.Text[:s.E.Pos-1] + s.E.Text[s.E.Pos:]
				s.E.Pos -= 1
			}
		}

	} else if rl.IsKeyPressed(rl.KeyDelete) {
		if shift {
			s.Destroy()
			return
		} else if s.E.Pos+1 <= len(s.E.Text) {
			if ctrl {
				s.DeleteWordFwd()
			} else {
				s.E.Text = s.E.Text[:s.E.Pos] + s.E.Text[s.E.Pos+1:]
			}
		}

		//navigation
	} else if rl.IsKeyPressed(rl.KeyRight) && s.E.Pos+1 <= len(s.E.Text) {
		if !ctrl {
			s.E.Pos += 1
		} else {
			s.MoveFwdWord()
		}
	} else if rl.IsKeyPressed(rl.KeyLeft) && s.E.Pos-1 >= 0 {
		if !ctrl {
			s.E.Pos -= 1
		} else {
			s.MoveBackWord()
		}
	} else if rl.IsKeyPressed(rl.KeyUp) {
		s.MoveUp()
	} else if rl.IsKeyPressed(rl.KeyDown) && s.E.Pos+2 <= len(s.E.Text) {
		s.MoveDown()
	} else if rl.IsKeyPressed(rl.KeyEnd) {
		if ctrl {
			s.E.Pos = len(s.E.Text)
		} else {
			if val := strings.Index(s.E.Text[s.E.Pos:], "\n"); val != -1 {
				s.E.Pos = val + s.E.Pos
			} else {
				s.E.Pos = len(s.E.Text)
			}
		}
	} else if rl.IsKeyPressed(rl.KeyHome) {
		if ctrl {
			s.E.Pos = 0
		} else {
			if val := strings.LastIndex(s.E.Text[:s.E.Pos], "\n"); val != -1 {
				s.E.Pos = val + 1
			} else {
				s.E.Pos = 0
			}
		}
	} else if ctrl && rl.IsKeyPressed(rl.KeyR) {
		s.Flags.IsNested = !s.Flags.IsNested
	} else if rl.IsKeyPressed(rl.KeyTab) {
		if shift {
			s.AddToList()
		} else if !s.isList {
			s.MakeList(GetListID())
		} else {
			s.ToggleAlign()
		}
	}

	//add cursor
	s.Name = s.E.Text[:s.E.Pos] + "|" + s.E.Text[s.E.Pos:]
}

func (s *Node) Draw() {

	//create rectangle
	s.UpdateRect()

	//set rect colour
	var colour rl.Color
	if s.E.IsEditing {
		colour = s.Theme.EditBG
	} else {
		colour = s.Theme.BG
	}

	//render text and rectangle
	var x, y int32
	if s.Theme.Circle {
		rl.DrawCircleV(s.Center, s.Theme.Radius, colour)
		rl.DrawCircleV(s.Center, s.Theme.Radius-5, PALETTE["White"])
		var offset int32 = 10
		x = int32(s.Center.X) + offset
		y = int32(s.Center.Y) + offset
	} else {
		x = int32(s.Pos.X) + s.Theme.Margin
		y = int32(s.Pos.Y) + s.Theme.Margin
		if !s.Theme.Rounded {
			rl.DrawRectangleRec(s.Rect, colour)
		} else {
			rl.DrawRectangleRounded(s.Rect, 0.4, 1, colour)
		}
	}
	if !s.hasImg {
		//colours
		var textColour rl.Color
		if s.E.IsEditing {
			textColour = s.Theme.EditText
		} else {
			textColour = s.Theme.Text
		}

		rl.DrawTextEx(FONT, s.Name, rl.Vector2{float32(x), float32(y)}, s.Theme.FontSize, s.Theme.FontSpacing, textColour)
	} else {
		rl.DrawTexture(s.Tex, x, y, rl.White)
	}

	s.DrawFlags()
}

func (s *Node) Destroy() {
	if s.Parent.ID != ROOT {
		s.Parent.EnableEditing()
	}
	//remove self from parent
	i := slices.Index(s.Parent.Childs, s)
	delete(NODES, s.ID) //free'd id's still unusable
	RemoveLua(s.ID)
	if s.isList {
		delete(LISTS, s.ListID)
		if s.Parent.isList && s.Parent.ListID == s.ListID {
			s.Parent.listNext = s.listNext
		}
	}
	if s.hasImg {
		rl.UnloadTexture(s.Tex)
	}
	if i == -1 {
		fmt.Println("huh, not in parents childs?")
	} else {
		s.Parent.Childs = slices.Delete(s.Parent.Childs, i, i+1)
	}
	// add back in all the children as childs of parent
	for _, v := range s.Childs {
		v.Parent = s.Parent
	}
	s.Parent.Childs = append(s.Parent.Childs, s.Childs...)
}
