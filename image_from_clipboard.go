package main
import (
  "github.com/gen2brain/raylib-go/raylib"
  "golang.design/x/clipboard"
  "bytes"
  //"github.com/h2non/filetype"
  "io"
  "fmt"
)

func StreamToByte(stream io.Reader) []byte {
  buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

// func TextureFromClipboard() rl.Texture2D {
//   r, _ := clipboard.ReadFromClipboard()
//   data := StreamToByte(r)
//   ext, _ := filetype.Get(data)
//   img := rl.LoadImageFromMemory(ext.Extension, data, int32(len(data)))
//   fmt.Println(img)
//   return rl.LoadTextureFromImage(img)
// }

func SetupClipboard() {
  err := clipboard.Init()
  if err != nil {
        panic(err)
  }
  fmt.Println("setup clipboard")
}

func TextureFromClipboard() (tex rl.Texture2D) {
  data := clipboard.Read(clipboard.FmtImage)
  if len(data) > 0 {
    img := rl.LoadImageFromMemory(".png", data, int32(len(data)))
    tex = rl.LoadTextureFromImage(img)
  }
  return
}
