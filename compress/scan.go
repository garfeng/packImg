package compress

import (
	"os"
	"path/filepath"
	"strings"
)

func scanSubRoots(name string, exts ...string) []string {
	extList := ".bmp"

	if len(exts) > 0 {
		extList = strings.Join(exts, "|")
	}
	extList = strings.ToLower(extList)

	w, _ := os.Open(name)
	list, _ := w.Readdir(-1)
	res := []string{}
	for _, v := range list {
		if v.IsDir() {
			r := filepath.Join(name, v.Name())
			res = append(res, scanSubRoots(r, exts...)...)
		} else {
			myExt := strings.ToLower(filepath.Ext(v.Name()))
			if strings.Contains(extList, myExt) {
				res = append(res, filepath.Join(name, v.Name()))
			}
		}
	}

	return res
}

func ReplaceExtTo(name string, dstExt string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return name + dstExt
	}

	idx := strings.LastIndex(name, ext)
	return name[:idx] + dstExt
}
