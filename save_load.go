package main
import  (
  "encoding/json"
  "io/ioutil"
  "fmt"
  "github.com/sqweek/dialog"
)

type jNode struct {
  X float32
  Y float32
  Text string
  ID int

  Childs []jNode
  Theme NodeTheme

  IsList bool
  ListNextID int
  ListID int



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
  open.Filter("*.json", "json")
  open.SetStartFile("mind_map.json")
  open.Title("Save Mind Map")
  path, err := open.Load()
  if err != nil {
    return
  }
  fmt.Println("SAVE!")
  j := toJnode(ROOT)
  l := toJLists(LISTS)
  save := save{Root:j, Lists:l}
  data, _ := json.MarshalIndent(save, "", " ")


  ioutil.WriteFile(path, data, 0644)
}
func Load() {
  open := dialog.File()
  open.Filter("*.json", "json")
  open.Title("Load Mind Map")
  path, err := open.Load()
  if err != nil {
    return
  }

  fmt.Println("LOAD!")

  //reset
  NODES = make(map[int]*Node)
  MAX_ID = 0
  MAX_LISTID = 0
  ResetLua()

  file, _ := ioutil.ReadFile(path)
  save := save{}
  json.Unmarshal([]byte(file), &save)
  fmt.Println(save.Lists[1])
  *ROOT = *toNode(save.Root, nil)
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

func toJnode(n *Node) (out jNode) {//maybe use pointers
  out.X = n.Pos.X
  out.Y = n.Pos.Y
  out.Text = n.E.Text
  out.Theme = n.Theme
  out.ID = n.ID

  out.IsList = n.isList
  if n.listNext != nil {
    out.ListNextID = n.listNext.ID
  }
  out.ListID = n.ListID

  for _, v := range n.Childs {
    out.Childs = append(out.Childs, toJnode(v))
  }

  return
}

func toNode(j jNode, p *Node) (out *Node){
  out = NewNode(j.Text, j.X, j.Y, j.ID)
  out.Parent = p
  out.Theme = j.Theme

  out.isList = j.IsList
  if p != nil{
    out.invertedLine = !p.invertedLine
  }else {
    out.invertedLine = true
  }

  for _, v := range j.Childs {
    out.Childs = append(out.Childs, toNode(v, out))
  }

  if out.isList {//after childs done
    out.ListID = j.ListID
    out.listNext = NODES[j.ListNextID]
  }

  return
}
