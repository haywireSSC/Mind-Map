package main
import (
  //TODO: change to gopher-lua and check licences(mit 4 go-lua)
  "github.com/yuin/gopher-lua"
  "fmt"
  "reflect"
  "strings"
  "regexp"
  "strconv"
)

var L *lua.LState

func MakeLuaAccesible(s *Node, key string) lua.LValue{
  switch key {
  case "x":
    return lua.LNumber(s.Pos.X)
  case "y":
    return lua.LNumber(s.Pos.Y)
  default:
    return lua.LString(fmt.Sprint(getAttr(s, key)))
  }
}
func getAttr(obj interface{}, fieldName string) reflect.Value {
    pointToStruct := reflect.ValueOf(obj) // addressable
    curStruct := pointToStruct.Elem()
    if curStruct.Kind() != reflect.Struct {
        panic("not struct")
    }
    curField := curStruct.FieldByName(fieldName) // type: reflect.Value
    if !curField.IsValid() {
        return reflect.ValueOf("value doesnt exist!")
    }
    return curField
}

func GetPropertyL(l *lua.LState) int{//gets a property by id, string name and returns it as string
  nodeTable := l.Get(1).(*lua.LTable)
  idL := nodeTable.RawGet(lua.LString("id"))
  ID, _ := strconv.Atoi(lua.LVAsString(idL))


  if ID == -1 {
    l.Push(lua.LString("wrong node path!"))
  }else {
    property := l.ToString(2)
    result := MakeLuaAccesible(NODES[ID], property)
    l.Push(result)
  }
  return 1
}

func DoPathString(originID int, path string) (result string) {
  var id int
  if node := DoPath(NODES[originID], path); node != nil {
    id = node.ID
  }else {
    id = -1
  }
  return fmt.Sprintf("NODES[%d]", id)
}

func DoPath(origin *Node, path string) (result *Node){// implement this, relative to a node, path starts with $ maybe idk
  //get path origin
  pathLen := len(path)
  var start int
  if path[0] == '$' {
    if pathLen > 1 && path[1] == '$' {
      result = ROOT
      start = 2
      if pathLen == 2 {
        return
      }
    }else {
      result = origin
      start = 1
      if pathLen == 1 {
        return
      }
    }
  }else {
    panic("path string does not start with $ or $$: " + path)
  }
  path = path[start:]

  if v, err := strconv.Atoi(path); err == nil {
    result = NODES[v]
  }else {
    items := strings.Split(path, "/")
    for _, v := range items {
      if v == "*" {//could be ..
        result = result.Parent
      }else {
        result = result.GetChildByName(v)
      }
      if result == nil {
        return
      }
    }
  }
  return
}

//first it compiles paths to nodes[some_id], or better is just some_id
//it also needs a get function to get normal properties in lua
//so you can do stuff like
//self.a = get($child/grandchild, property_name)
//issue is paths outside get need to be table, and inside is ID

//either a special char like raw$child/grandchild, (probs best)
//then another node can do
//$*/*.a

//WATER PLANTS

func RunLua(ID int) string{//make nodes both run and addlua(done), just fixing node deletion and error handling
  lFunc := L.GetGlobal("FUNCS").(*lua.LTable).RawGetInt(ID)//STORE FUNCTION IN NODE
  err := L.CallByParam(lua.P{Fn:lFunc, NRet:1, Protect:true})
  var str string
  if err != nil {
    str = err.Error()
  }else {
    lStr, _ := L.Get(-1).(lua.LString)
    str = lStr.String()
    L.Pop(1)
  }

  return str
}

func StartLua() {
  L = lua.NewState()
  L.DoString("FUNCS = {}; NODES = {}")
  L.SetGlobal("get", L.NewFunction(GetPropertyL))
}

func ParseCode(ID int, text string){
  code := fmt.Sprintf(`
    FUNCS[%[1]d] = function()
    if NODES[%[1]d] == nil then NODES[%[1]d] = {["id"] = %[1]d} end
    `, ID, NODES[ID].Name)


  var proprepl = func (in string) string {return DoPathString(ID, in)}
  re := regexp.MustCompile(`\$[\w/\*]*`)//replace paths with NODES[id]
  text = re.ReplaceAllStringFunc(text, proprepl)

  extraCode := &[]string{}
  var coderepl = func(in string) string {*extraCode = append(*extraCode, in[2:len(in)-2]); return ""}
  re = regexp.MustCompile(`\|\|.*?\|\|`)//remove the ||expr||
  text = re.ReplaceAllStringFunc(text, coderepl)


  re = regexp.MustCompile(`{{(.*?)}}|(NODES\S*)`)//do the tostring
  text = re.ReplaceAllString(text, `]] .. tostring($1$2) .. [[`)
  text = "[[" + text + "]]"

  code += " " + strings.Join(*extraCode, " ") + " "

  code += fmt.Sprintf(`return %s end`, text)//add text to code


  //fmt.Println(code)
  L.DoString(code)
}
