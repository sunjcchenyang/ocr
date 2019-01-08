package service

import "github.com/sbinet/go-python"

type PythonLib struct {
	Str2Py func(string)*python.PyObject
	Py2Str func(*python.PyObject)string
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
func (PythonLib *PythonLib)ImportModule(dir, name string) *python.PyObject {
	module := python.PyImport_ImportModule("sys")
	path := module.GetAttrString("path")
	python.PyList_Insert(path, 0, PythonLib.Str2Py(dir))
	return python.PyImport_ImportModule(name)
}
