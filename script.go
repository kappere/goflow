package goflow

import (
	"encoding/json"
	"errors"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type Script interface {
	String() string
	Run(map[string]interface{}) (interface{}, error)
}

type LuaScript struct {
	Data string
}

func (s LuaScript) String() string {
	return s.Data
}

func (s LuaScript) Run(param map[string]interface{}) (result interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("unknown lua script error: %s, panic: %v", s.Data, p)
		}
	}()
	L := lua.NewState()
	defer L.Close()
	if err = L.DoString(s.Data); err != nil {
		return nil, err
	}
	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal("run"),
		NRet:    1,
		Protect: true,
	}, s.mapToTable(param))
	if err != nil {
		return nil, err
	}
	retValue := L.Get(-1)
	L.Pop(1)
	var r map[string]interface{}
	if retTable, ok := retValue.(*lua.LTable); ok {
		r = s.tableToMap(retTable)
	} else {
		r = map[string]interface{}{
			"success": false,
			"message": "invalid return value: " + retValue.String(),
		}
	}
	if !r["success"].(bool) {
		return nil, errors.New(r["message"].(string))
	}
	return r["data"], nil
}

func (s LuaScript) mapToTable(m map[string]interface{}) *lua.LTable {
	jsonBytes, _ := json.Marshal(m)
	m2 := map[string]interface{}{}
	json.Unmarshal(jsonBytes, &m2)
	t := &lua.LTable{}
	for k1, v1 := range m2 {
		s.getGoRetValueMap(k1, v1, t)
	}
	return t
}

func (s LuaScript) tableToMap(t *lua.LTable) map[string]interface{} {
	ret := map[string]interface{}{}
	t.ForEach(func(l1, l2 lua.LValue) {
		s.getLuaRetValueMap(l1, l2, ret)
	})
	return ret
}

func (s LuaScript) getLuaRetValueMap(k, v lua.LValue, m map[string]interface{}) {
	switch v2 := v.(type) {
	case *lua.LTable:
		m2 := make(map[string]interface{})
		m[k.String()] = m2
		v2.ForEach(func(k1, v1 lua.LValue) {
			s.getLuaRetValueMap(k1, v1, m2)
		})
	case lua.LString:
		m[k.String()] = string(v2)
	case lua.LBool:
		m[k.String()] = bool(v2)
	case lua.LNumber:
		m[k.String()] = float64(v2)
	case *lua.LNilType:
		m[k.String()] = nil
	}
}

func (s LuaScript) getGoRetValueMap(k string, v interface{}, m *lua.LTable) {
	switch v2 := v.(type) {
	case map[string]interface{}:
		m2 := &lua.LTable{}
		m.RawSet(lua.LString(k), m2)
		for k1, v1 := range v2 {
			s.getGoRetValueMap(k1, v1, m2)
		}
	case string:
		m.RawSet(lua.LString(k), lua.LString(v2))
	case bool:
		m.RawSet(lua.LString(k), lua.LBool(v2))
	case float64:
		m.RawSet(lua.LString(k), lua.LNumber(v2))
	case nil:
		m.RawSet(lua.LString(k), lua.LNil)
	}
}
