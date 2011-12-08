package main

import (
	"reflect"
	"unsafe"
	"x-go-binding.googlecode.com/hg/xgb"
)

type IdList []xgb.Id

func (l IdList) Contains(id xgb.Id) bool {
	for _, i := range l {
		if i == id {
			return true
		}
	}
	return false
}

func propReplyAtoms(prop *xgb.GetPropertyReply) IdList {
	if prop == nil || prop.ValueLen == 0 {
		return nil
	}
	atom_size := uintptr(prop.Format / 8)
	if atom_size != reflect.TypeOf(xgb.Id(0)).Size() {
		panic("Property reply has wrong format for atoms")
	}
	num_atoms := prop.ValueLen / uint32(atom_size)
	return (*[1<<24]xgb.Id)(unsafe.Pointer(&prop.Value[0]))[:num_atoms]
}
