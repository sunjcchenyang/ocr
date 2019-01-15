package service

import (
	"github.com/sbinet/go-python"
	"go.uber.org/zap"
)

type LibInterface interface {
	Init()
	LoadLib(path string,name string)
	GetResult(image string) string
}

type PythonLib struct {
	Str2Py func(string) *python.PyObject
	Py2Str func(*python.PyObject) string
	LibHabdle *python.PyObject
}

func NewLib() *PythonLib {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
	PythonLib := &PythonLib {
		Str2Py: python.PyString_FromString,
		Py2Str: python.PyString_AsString,
	}
	return PythonLib
}
//init

func (p *PythonLib)LoadLib(dir, name string) *python.PyObject {
	module := python.PyImport_ImportModule("sys")
	path := module.GetAttrString("path")
	python.PyList_Insert(path, 0, p.Str2Py(dir))
	p.LibHabdle = python.PyImport_ImportModule(name)
	return p.LibHabdle
}

func (p *PythonLib) Init()  {
	Logger.Info("[Init]", zap.String("@Init", "inter init function"))
	f := p.LibHabdle.GetAttrString("init")
	f.CallFunction()
}

func (p *PythonLib) GetResult() string{
	f := p.LibHabdle.GetAttrString("ocr_recog")
	argv := python.PyTuple_New(1)
	python.PyTuple_SetItem(argv, 0, p.Str2Py("image.png"))
	res := f.Call(argv, python.Py_None)
	Logger.Info("[GetResult]",zap.Any("res",p.Py2Str(res)))
	return p.Py2Str(res)
}













