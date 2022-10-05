package main
import (
  "github.com/gen2brain/raylib-go/raylib"
)

type EditorCamera struct {
  Cam rl.Camera2D
  MousePos rl.Vector2
  DoubleClick bool
  QuickClick bool

  quickClickTimeout float64
  lastQclick float64

  doubleClickTimeout float64
  lastClick float64

  ctrl bool
  shift bool
}

func NewEditorCamera() (inst EditorCamera){
  inst.Cam = rl.NewCamera2D(rl.Vector2{}, rl.Vector2{}, 0, 1)
  inst.doubleClickTimeout = 0.5

  inst.quickClickTimeout = 0.2

  return
}

func (s *EditorCamera) RefreshTarget() {
  s.Cam.Target = s.MousePos
}

func CheckSaveLoad() {
  if rl.IsKeyPressed(rl.KeyS) && CAMERA.ctrl {
    Save()
  }else if rl.IsKeyPressed(rl.KeyO) && CAMERA.ctrl {
    Load()
  }
}

func (s *EditorCamera) Update() {
  s.ctrl = rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)
  s.shift = rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

  CheckSaveLoad()

  s.MousePos = rl.GetScreenToWorld2D(rl.GetMousePosition(), s.Cam)
  wheelDelta := rl.GetMouseWheelMove()
  if wheelDelta != 0 {
    s.Cam.Zoom += wheelDelta
    if s.Cam.Zoom < 1 {
      s.Cam.Zoom = 1
    }
  }

  if rl.IsMouseButtonDown(2) {
    if rl.IsMouseButtonPressed(2) {
      s.RefreshTarget()
    }
    s.Cam.Offset = rl.GetMousePosition()
  }

  s.QuickClick = false
  if rl.IsMouseButtonPressed(0) {
    s.lastQclick = rl.GetTime()
  }else if rl.IsMouseButtonReleased(0) && (rl.GetTime() - s.lastQclick) < s.quickClickTimeout {
    s.QuickClick = true
  }

  s.DoubleClick = false
  if rl.IsMouseButtonPressed(0) {
    if (rl.GetTime() - s.lastClick) < s.doubleClickTimeout {
      s.DoubleClick = true
      s.lastClick = 0
    }else {
      s.lastClick = rl.GetTime()
    }
  }
}
