package main
import (
  "github.com/gen2brain/raylib-go/raylib"
  "fmt"
)
var (
  CAMERA EditorCamera
  ROOT *Node
  NODES map[int]*Node = make(map[int]*Node)
  HANDLER RootNode
  LINE_RENDERER LineRenderer
  FONT rl.Font
  MAX_ID int
  MAX_LISTID int
  LISTS map[int]*List = make(map[int]*List)
)

func main() {
  StartLua()
  defer L.Close()

  rl.SetConfigFlags(rl.FlagMsaa4xHint)
  fmt.Println("start")
  rl.InitWindow(500,500, "mind map thing")
  rl.SetTargetFPS(60)

  SetupClipboard()

  FONT = rl.LoadFont("johnston-itc-std-bold.otf")
  HANDLER = NewRootNode()
  ROOT = HANDLER.node
  CAMERA = NewEditorCamera()
  LINE_RENDERER = NewLineRenderer()


  for !rl.WindowShouldClose() {

    CAMERA.Update()



    rl.BeginDrawing()
    rl.ClearBackground(rl.White)

    rl.BeginMode2D(CAMERA.Cam)
    LINE_RENDERER.Draw()
    HANDLER.Update()
    rl.EndMode2D()

    rl.EndDrawing()
  }

  rl.CloseWindow()
}
