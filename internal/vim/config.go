package vim

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// LoadConfig executes the file at path as Lua 5.1 and applies bindings to engine.
// Returns nil if path does not exist. Logs warnings for unknown action names (non-fatal).
func LoadConfig(path string, engine *Engine) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	L := lua.NewState()
	defer L.Close()

	// Register vim global table
	vimTable := setupVimTable(L, engine)
	L.SetGlobal("vim", vimTable)

	if err := L.DoFile(path); err != nil {
		return fmt.Errorf("lua config %s: %w", path, err)
	}
	return nil
}

// LoadConfigString executes a Lua string and applies bindings. Useful for testing.
func LoadConfigString(src string, engine *Engine) error {
	L := lua.NewState()
	defer L.Close()

	vimTable := setupVimTable(L, engine)
	L.SetGlobal("vim", vimTable)

	if err := L.DoString(src); err != nil {
		return fmt.Errorf("lua config: %w", err)
	}
	return nil
}

func setupVimTable(L *lua.LState, engine *Engine) *lua.LTable {
	vimTable := L.NewTable()

	// vim.g with __newindex that intercepts mapleader
	gTable := L.NewTable()
	gMeta := L.NewTable()
	L.SetField(gMeta, "__newindex", L.NewFunction(func(L *lua.LState) int {
		// args: table, key, value
		key := L.CheckString(2)
		val := L.CheckString(3)
		if key == "mapleader" {
			engine.SetLeader(val)
		}
		// Silently accept all other writes
		return 0
	}))
	L.SetMetatable(gTable, gMeta)
	L.SetField(vimTable, "g", gTable)

	// vim.keymap table
	keymapTable := L.NewTable()

	// vim.keymap.set(mode, lhs, rhs)
	// mode: "n", "i", "v" or array of those
	// lhs: key sequence string (may contain <leader>)
	// rhs: action name string
	L.SetField(keymapTable, "set", L.NewFunction(func(L *lua.LState) int {
		// Determine modes
		var modes []string
		arg1 := L.Get(1)
		switch v := arg1.(type) {
		case lua.LString:
			modes = []string{string(v)}
		case *lua.LTable:
			v.ForEach(func(_, val lua.LValue) {
				if s, ok := val.(lua.LString); ok {
					modes = append(modes, string(s))
				}
			})
		default:
			L.ArgError(1, "expected string or table of strings for mode")
			return 0
		}

		lhs := L.CheckString(2)
		rhs := L.CheckString(3)

		// Expand <leader> in lhs
		expandedLHS := expandLeader(lhs, engine.parser.leaderKey)

		// Resolve action name
		act, ok := ActionName[rhs]
		if !ok {
			log.Printf("vim config: unknown action %q (binding %q ignored)", rhs, lhs)
			return 0
		}

		// Apply to each mode
		for _, modeStr := range modes {
			switch strings.ToLower(modeStr) {
			case "n":
				engine.BindKey(ModeNormal, expandedLHS, act)
			case "i":
				engine.BindKey(ModeInsert, expandedLHS, act)
			case "v":
				engine.BindKey(ModeVisual, expandedLHS, act)
				engine.BindKey(ModeVisualLine, expandedLHS, act)
				engine.BindKey(ModeVisualBlock, expandedLHS, act)
			default:
				log.Printf("vim config: unknown mode %q for binding %q", modeStr, lhs)
			}
		}
		return 0
	}))
	L.SetField(vimTable, "keymap", keymapTable)

	return vimTable
}

// expandLeader replaces literal "<leader>" in lhs with the leader string.
func expandLeader(lhs, leader string) string {
	if leader == "" {
		return lhs
	}
	return strings.ReplaceAll(lhs, "<leader>", leader)
}
