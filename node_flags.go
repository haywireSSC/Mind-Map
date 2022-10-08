package main
import (
  "github.com/gen2brain/raylib-go/raylib"
)

type NodeFlags struct {
  IsNested bool
  IsCode bool
}
func (s *Node) DrawFlags() {
  text := ""
  if s.Flags.IsNested {
    text += "l"
  }
  if s.Flags.IsCode {
    text += "p"
  }

  size := rl.MeasureTextEx(SYMBOL_FONT, text, 32, 4)

  var pos rl.Vector2
  pos.X = s.Pos.X// - size.X
  pos.Y = s.Pos.Y - size.Y
  rl.DrawTextEx(SYMBOL_FONT, text, pos, 32, 4, s.Theme.BG)
}
