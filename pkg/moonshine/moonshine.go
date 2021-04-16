package moonshine

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
	lua "github.com/yuin/gopher-lua"
)

type builder struct {
	s llb.State
}

func (b *builder) from(luaState *lua.LState) int {
	t := luaState.ToTable(1)

	base := t.RawGetString("base")
	b.s = llb.Image(base.String())

	t.ForEach(func(l1 lua.LValue, l2 lua.LValue) {
		switch l2.Type().String() {
		case "table":
			t := l2.(*lua.LTable)
			f := t.RawGetString("func")

			if f.String() == "run" {
				t.ForEach(func(k lua.LValue, v lua.LValue) {
					if k.String() != "func" {
						b.s = b.s.Run(
							llb.Shlex(v.String()),
						).Root()
					}
				})
			}
		case "string":
			// fmt.Println("string", l1.String(), l2.String())
		default:
			fmt.Println("unknown thing")
		}
	})

	return 0
}

func (b *builder) run(luaState *lua.LState) int {
	t := luaState.ToTable(1)

	t.RawSetString("func", lua.LString("run"))

	luaState.Push(t)

	return 1
}

func (b *builder) copy(luaState *lua.LState) int {
	t := luaState.ToTable(1)

	t.RawSetString("func", lua.LString("copy"))

	luaState.Push(t)

	return 1
}

func (b *builder) workdir(luaState *lua.LState) int {
	return 0
}

func DockerLuaToLLB(l string) (llb.State, error) {
	luaState := lua.NewState()
	defer luaState.Close()

	b := builder{}
	luaState.SetGlobal("from", luaState.NewFunction(b.from))
	luaState.SetGlobal("copy", luaState.NewFunction(b.copy))
	luaState.SetGlobal("run", luaState.NewFunction(b.run))
	luaState.SetGlobal("workdir", luaState.NewFunction(b.workdir))

	lines := strings.Split(l, "\n")
	if err := luaState.DoString(strings.Join(lines[1:], "\n")); err != nil {
		return llb.State{}, err
	}

	return b.s, nil
}
