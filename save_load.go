package main
import  (
  "encoding/json"
  "os"
  "fmt"
  "github.com/sqweek/dialog"
  "image"
  "image/png"
  "archive/zip"
  "github.com/gen2brain/raylib-go/raylib"
  "strconv"
  "io"
)

type jNode struct {
  X float32
  Y float32
  Text string
  ID int

  Childs []jNode
  Theme NodeTheme
  Flags NodeFlags

  IsList bool
  ListNextID int
  ListID int

  ImageIdx int



  //Tex rl.Texture2D - add saving for pictures
}

type jList struct {
  IsAlignX bool
  AlignX float32

  IsAlignY bool
  AlignY float32

  RootID int
}

type save struct {
  Root jNode
  Lists map[int]jList
}

func Save() {//error logging, add extracting theme and if list node
  open := dialog.File()
  open.Filter("*.zip", "zip")
  open.SetStartFile("test.zip")
  open.Title("Save Mind Map")
  path, err := open.Save()
  if err != nil {
    return
  }
  fmt.Println("SAVE!")
  j, images := toJnode(NODES[ROOT], []image.Image{})
  l := toJLists(LISTS)
  save := save{Root:j, Lists:l}
  data, _ := json.MarshalIndent(save, "", " ")


  f, _ := os.Create(path)
  ar := zip.NewWriter(f)
  file, _ := ar.Create("nodes.json")
  file.Write(data)
  for i, img := range images {
    out, _ := ar.Create(strconv.Itoa(i) + ".png")
    png.Encode(out, img)
  }
  ar.Close()
}


func Load() {
  open := dialog.File()
  open.Filter("*.zip", "zip")
  open.Title("Load Mind Map")
  path, err := open.Load()
  if err != nil {
    return
  }

  fmt.Println("LOAD!")

  //reset
  NODES = make(map[int]*Node)//doesnt unload textures
  MAX_ID = 0
  MAX_LISTID = 0
  ResetLua()


  ar, _ := zip.OpenReader(path)
  defer ar.Close()

  var data []byte
  images := []*rl.Image{}

  for _, v := range ar.File {
    file, _ := v.Open()
    if v.FileHeader.Name == "nodes.json" {
      data, _ = io.ReadAll(file)
    }else {
      img, _ := png.Decode(file)
      images = append(images, rl.NewImageFromImage(img))
    }
    file.Close()
  }

  save := save{}
  json.Unmarshal(data, &save)
  HANDLER.root = toNode(save.Root, nil, images)
  ROOT = HANDLER.root.ID
  LISTS = toLists(save.Lists)
}

func toJLists(lists map[int]*List) (out map[int]jList){
  out = make(map[int]jList)
  for k, v := range lists {
    l := jList{}
    l.IsAlignX = v.isAlignX
    l.IsAlignY = v.isAlignY
    l.AlignX = v.alignX
    l.AlignY = v.alignY
    l.RootID = v.root.ID
    out[k] = l
  }
  return
}

func toLists(lists map[int]jList) (out map[int]*List){
  out = make(map[int]*List)
  for k, v := range lists {
    if k > MAX_LISTID {
      MAX_LISTID = k
    }
    l := List{}
    l.isAlignX = v.IsAlignX
    l.isAlignY = v.IsAlignY
    l.alignX = v.AlignX
    l.alignY = v.AlignY
    l.root = NODES[v.RootID]
    out[k] = &l
  }
  return
}

func toJnode(n *Node, im []image.Image) (out jNode, outIm []image.Image) {//maybe use pointers
  outIm = im
  out.X = n.Pos.X
  out.Y = n.Pos.Y
  out.Text = n.E.Text
  out.Theme = n.Theme
  out.ID = n.ID
  out.Flags = n.Flags

  out.IsList = n.isList
  if n.listNext != nil {
    out.ListNextID = n.listNext.ID
  }
  out.ListID = n.ListID

  if n.hasImg {
    outIm = append(outIm, rl.LoadImageFromTexture(n.Tex).ToImage())
    out.ImageIdx = len(outIm)-1
  }else {
    out.ImageIdx = -1
  }

  for _, v := range n.Childs {
    c, newIm := toJnode(v, outIm)
    out.Childs = append(out.Childs, c)
    outIm = append(outIm, newIm...)
  }

  return
}

func toNode(j jNode, p *Node, images []*rl.Image) (out *Node){
  out = NewNode(j.Text, j.X, j.Y, j.ID)
  out.Parent = p
  out.Theme = j.Theme
  out.Flags = j.Flags

  out.isList = j.IsList
  if p != nil{
    out.invertedLine = !p.invertedLine
  }else {
    out.invertedLine = true
  }

  if j.ImageIdx != -1 {
    out.Tex = rl.LoadTextureFromImage(images[j.ImageIdx])
    out.hasImg = true
  }

  for _, v := range j.Childs {
    out.Childs = append(out.Childs, toNode(v, out, images))
  }

  if out.isList {//after childs done
    out.ListID = j.ListID
    out.listNext = NODES[j.ListNextID]
  }

  return
}
