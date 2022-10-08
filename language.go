package main

import (
	//TODO: change to gopher-lua and check licences(mit 4 go-lua)
	"fmt"
	"regexp"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

var L *lua.LState

func GetLuaAccesible(s *Node, key string) lua.LValue {
	switch key {
	case "x":
		return lua.LNumber(s.Pos.X)
	case "y":
		return lua.LNumber(s.Pos.Y)
	}
	return lua.LNil
}

func SetLuaAccesible(s *Node, key string, value lua.LValue) {
	switch key {
	case "x":
		s.Pos.X = float32(value.(lua.LNumber))
	case "y":
		s.Pos.Y = float32(value.(lua.LNumber))
	}
}

func SetPropertyL(l *lua.LState) int {
	nodeTable := l.Get(1).(*lua.LTable)
	idL := nodeTable.RawGet(lua.LString("id"))
	ID, _ := strconv.Atoi(lua.LVAsString(idL))

	if ID == -1 {
		l.Push(lua.LString("wrong node path!"))
	} else {
		property := l.ToString(2)
		SetLuaAccesible(NODES[ID], property, l.Get(3))
		l.Push(lua.LString("test")) //temp
	}

	return 1
}

func GetPropertyL(l *lua.LState) int { //gets a property by id, string name and returns it as string
	nodeTable := l.Get(1).(*lua.LTable)
	idL := nodeTable.RawGet(lua.LString("id"))
	ID, _ := strconv.Atoi(lua.LVAsString(idL))

	if ID == -1 {
		l.Push(lua.LString("wrong node path!"))
	} else {
		property := l.ToString(2)
		result := GetLuaAccesible(NODES[ID], property)
		l.Push(result)
	}
	return 1
}

func DoPathString(originID int, path string) (result string) {
	var id int
	if node := DoPath(NODES[originID], path); node != nil {
		id = node.ID
	} else {
		id = -1
	}
	return fmt.Sprintf("NODES[%d]", id)
}

func DoPath(origin *Node, path string) (result *Node) { // implement this, relative to a node, path starts with $ maybe idk
	//get path origin
	pathLen := len(path)
	var start int
	if path[0] == '$' {
		if pathLen > 1 && path[1] == '$' {
			result = NODES[ROOT]
			start = 2
			if pathLen == 2 {
				return
			}
		} else {
			result = origin
			start = 1
			if pathLen == 1 {
				return
			}
		}
	} else {
		panic("path string does not start with $ or $$: " + path)
	}
	path = path[start:]

	if v, err := strconv.Atoi(path); err == nil {
		result = NODES[v]
	} else {
		items := strings.Split(path, "/")
		for _, v := range items {
			if v == "*" { //could be ..
				result = result.Parent
			} else {
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

func RunLua(ID int) string { //make nodes both run and addlua(done), just fixing node deletion and error handling
	L.GetGlobal("NODES")
	lFunc := L.GetGlobal("FUNCS").(*lua.LTable).RawGetInt(ID) //STORE FUNCTION IN NODE
	err := L.CallByParam(lua.P{Fn: lFunc, NRet: 1, Protect: true})
	var str string
	if err != nil {
		str = err.Error()
	} else {
		lStr, _ := L.Get(-1).(lua.LString)
		str = lStr.String()
		L.Pop(1)
	}

	return str
}

func RemoveLua(ID int) { //store the table
	L.GetGlobal("FUNCS").(*lua.LTable).RawSetInt(ID, lua.LNil)
	L.GetGlobal("NODES").(*lua.LTable).RawSetInt(ID, lua.LNil)
}

func StartLua() {
	L = lua.NewState()
	ResetLua()
	L.SetGlobal("get", L.NewFunction(GetPropertyL))
	L.SetGlobal("set", L.NewFunction(SetPropertyL))
}
func ResetLua() {
	L.DoString("FUNCS = {}; NODES = {}")
}

func AddLua(n *Node, text string) {
	code, hasCode := ParseCode(n.ID, text)
	n.Flags.IsCode = hasCode
	// if !hasCode {
	// 	return
	// }

	precode := fmt.Sprintf(`NODES[%[1]d] = {id = %[1]d}
		setmetatable(NODES[%[1]d], NODES[%[1]d])

		NODES[%[1]d].__index = function(table, key)
		value = get(NODES[%[1]d], key)
		if value == nil then
		return rawget(table, key)
		else
		return value
		end
		end

		NODES[%[1]d].__newindex = function(table, key, value)
		if set(NODES[%[1]d], key, value) == nil then
		rawset(table, key, value)
		end
		end`, n.ID)

	L.DoString(precode)
	L.DoString(code)
}

func ParseCode(ID int, text string) (code string, hasCode bool) {
	code = fmt.Sprintf(`
    FUNCS[%[1]d] = function()
    `, ID)

	var proprepl = func(in string) string { return DoPathString(ID, in) }
	re := regexp.MustCompile(`\$[\w/\*]*`) //replace paths with NODES[id]
	ptext := text
	text = re.ReplaceAllStringFunc(text, proprepl)
	if ptext != text {
		hasCode = true
	}

	extraCode := &[]string{}
	var coderepl = func(in string) string { *extraCode = append(*extraCode, in[2:len(in)-2]); return "" }
	re = regexp.MustCompile(`\|\|.*?\|\|`) //remove the ||expr||
	ptext = text
	text = re.ReplaceAllStringFunc(text, coderepl)
	if ptext != text {
		hasCode = true
	}

	re = regexp.MustCompile(`{{(.*?)}}|(NODES\S*)`) //do the tostring
	ptext = text
	text = re.ReplaceAllString(text, `]] .. tostring($1$2) .. [[`)
	if ptext != text {
		hasCode = true
	}
	text = "[[" + text + "]]"

	code += " " + strings.Join(*extraCode, " ") + " "

	code += fmt.Sprintf(`return %s end`, text) //add text to code

	return
}
