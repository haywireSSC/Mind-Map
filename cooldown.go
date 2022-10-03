package main
import (
  "github.com/gen2brain/raylib-go/raylib"
)

type Cooldown struct {
  Gap float64
  lasttime float64
}


func (s *Cooldown) tick() bool {
  time := rl.GetTime()
  if time - s.lasttime > s.Gap {
    s.lasttime = time
    return true
  }
  return false
}
