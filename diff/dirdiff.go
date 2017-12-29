package diff

// import (
// 	"fmt"
// )

// DiffOnly used for dirdiff
func DiffOnly(src []string, dst []string) ([]string, []string) {
	srcFile = src
	dstFile = dst
	srcLen = len(srcFile)
	dstLen = len(dstFile)

	pTmp := newpoint(-1, -1)
	getPath(pTmp)
	pathPoint := getMostDepth()

	rm := []string{}
	add := []string{}
	pOne := newpoint(0, 0)
	getResult := func(pOne, pPoint *point) {
		for j := pOne.x; j < pPoint.x; j++ {
			rm = append(rm, srcFile[j])
		}
		for j := pOne.y; j < pPoint.y; j++ {
			add = append(add, dstFile[j])
		}
	}
	for i := len(pathPoint) - 2; i >= 0; i-- {
		getResult(pOne, pathPoint[i])
		pOne = newpoint(pathPoint[i].x+1, pathPoint[i].y+1)
	}
	pEnd := newpoint(srcLen, dstLen)
	getResult(pOne, pEnd)
	return rm, add
}