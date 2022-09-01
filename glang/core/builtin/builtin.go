package buitin

import (
	"almeng.com/glang/core/syntax"
	"almeng.com/glang/global"
	"github.com/almenglee/general"
	"github.com/llir/llvm/ir"
)

var GlobalGroup = &syntax.Group{}
var NewLine *ir.Global

func InitTypes(list *general.List[syntax.Decl]) {
	for _, v := range ITypes {
		list.Append(InitialTypeNode(v.N))
	}
}

func InitConsts(list *general.List[syntax.Decl]) {
	for _, v := range Consts {
		list.Append(_constNode(v))
	}
}

func _constNode(c _const) syntax.Decl {
	return &syntax.VarDecl{
		Group:    GlobalGroup,
		NameList: &syntax.Name{Value: c.Name},
		Type:     &syntax.Name{Value: c.Type},
		Values:   nil,
	}
}

func InitialTypeNode(name string) syntax.Decl {
	t := &syntax.Name{Value: name}
	return &syntax.TypeDecl{
		Group: GlobalGroup,
		Name:  t,
		Alias: false,
		Type:  t,
	}
}

func InitModule(m *ir.Module) {
	for _, v := range ITypes {
		m.NewTypeDef(v.N, v.T)
	}
	for _, v := range Consts {
		m.NewGlobalDef(v.Name, v.IConst)
	}

	//zero := constant.NewInt(Int, 0)
	//seed := m.NewGlobalDef("seed", zero)
	NewLine = global.NewGlobalCharArrayConstant("\n")
	Println = _println()
	Print = _print()
	RegisterFunc(m, Println)
	RegisterFunc(m, Print)
	for _, v := range Funcs {
		RegisterFunc(m, v)
	}

	//main

	//	a := constant.NewInt(Int, 0x15A4E35) // multiplier of the PRNG.
	//	c := constant.NewInt(Int, 1)
	//
	//	main := m.NewFunc(
	//		"main",
	//		types.I32)
	//	entry := main.NewBlock("")
	//
	//	tmp1 := entry.NewLoad(Int, seed)
	//	tmp2 := entry.NewMul(tmp1, a)
	//
	//	tmp3 := entry.NewAdd(tmp2, c)
	//	_ = tmp3
	//	hi := global.NewGlobalString(entry,
	//		`Hello world!
	//
	//Testing print function...
	//
	//Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin vehicula nec sem a vulputate. Fusce eget mollis arcu, eu semper dui. Cras facilisis, leo ac dignissim faucibus, sem nisl commodo nisl, efficitur volutpat eros felis non neque. Interdum et malesuada fames ac ante ipsum primis in faucibus. Aenean justo felis, vulputate viverra dictum id, aliquet eget urna. Vestibulum pellentesque felis quis purus facilisis, vel vestibulum eros porta. In vestibulum massa tellus, ac tincidunt ex ullamcorper a. Quisque ipsum lacus, commodo eu elit nec, aliquet elementum lacus. Nam facilisis luctus interdum. Integer nec augue neque.
	//
	//Test Complete
	//
	//`)
	//	bye := global.NewGlobalString(entry, "bye world!")
	//	_ = entry.NewCall(Print, hi)
	//	_ = entry.NewCall(Println, bye)
	//	entry.NewRet(constant.NewInt(types.I32, 0))
}
