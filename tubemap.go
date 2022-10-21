package main
import (
  "github.com/gen2brain/raylib-go/raylib"
  "fmt"
)
var (
  CAMERA EditorCamera
  ROOT int
  NODES map[int]*Node = make(map[int]*Node)
  HANDLER RootNode
  LINE_RENDERER LineRenderer
  FONT rl.Font
  SYMBOL_FONT rl.Font
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

  rl.SetExitKey(0)

  FONT = rl.LoadFont("johnston-itc-std-bold.otf")
  SYMBOL_FONT = rl.LoadFont("RailwayAlternate.otf")
  HANDLER = NewRootNode()
  ROOT = HANDLER.root.ID
  CAMERA = NewEditorCamera()
  LINE_RENDERER = NewLineRenderer()


  for !rl.WindowShouldClose() {

    CAMERA.Update()



    rl.BeginDrawing()

    HANDLER.Update()

    rl.EndDrawing()
  }

  rl.CloseWindow()
}
