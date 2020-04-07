// +build go1.12
// +build !go1.15

package goloader

import (
	"encoding/binary"
	"strconv"
	"unsafe"
)

func _addStackObject(code *CodeReloc, fi *funcInfoData, seg *segment, symPtr map[string]uintptr) {
	if len(fi.funcdata) > _FUNCDATA_StackObjects && fi.funcdata[_FUNCDATA_StackObjects] != 0xFFFFFFFF {
		stackObjectRecordSize := unsafe.Sizeof(stackObjectRecord{})
		b := code.Mod.stkmaps[fi.funcdata[_FUNCDATA_StackObjects]]
		n := *(*int)(unsafe.Pointer(&b[0]))
		p := unsafe.Pointer(&b[PtrSize])
		for i := 0; i < n; i++ {
			obj := *(*stackObjectRecord)(p)
			var name string
			for _, v := range fi.Var {
				if v.Offset == (int64)(obj.off) {
					name = v.Type.Name
					break
				}
			}
			if len(name) == 0 {
				name = fi.stkobjReloc[i].Sym.Name
			}
			ptr, ok := symPtr[name]
			if !ok {
				ptr, ok = seg.typeSymPtr[name]
			}
			if !ok {
				sprintf(&seg.err, "unresolve external:", strconv.Itoa(i), " ", fi.name, "\n")
			} else {
				off := PtrSize + i*(int)(stackObjectRecordSize) + PtrSize
				if PtrSize == 4 {
					binary.LittleEndian.PutUint32(b[off:], *(*uint32)(unsafe.Pointer(&ptr)))
				} else {
					binary.LittleEndian.PutUint64(b[off:], *(*uint64)(unsafe.Pointer(&ptr)))
				}
			}
			p = add(p, stackObjectRecordSize)
		}
	}
}
