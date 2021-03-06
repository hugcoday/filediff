// By jacenr
// Create: 2017-12-10
// Usage:
//       import "github.com/jacenr/filediff/diff"
//       result, _ := diff.Diff("file1", "file2")

package diff

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var srcFile []string
var dstFile []string
var srcLen int
var dstLen int
var path = map[*point][]*point{}
var newed = map[string]*point{}

// src:x, dst: y
type point struct {
	x      int
	y      int
	parent *point
	depth  int
}

func (p *point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func newpoint(x int, y int) *point {
	p := new(point)
	p.x = x
	p.y = y
	p.parent = nil
	p.depth = 0
	return p
}

// Check the point whether has been created.
func checkNew(x, y int, p *point) *point {
	xyStr := strconv.Itoa(x) + strconv.Itoa(y)
	v, ok := newed[xyStr]
	if !ok {
		pNew := newpoint(x, y)
		pNew.parent = p
		pNew.depth = p.depth + 1
		newed[xyStr] = pNew
		return pNew
	}
	if v.depth < p.depth+1 {
		v.parent = p
		v.depth = p.depth + 1
		return v
	}
	return nil
}

// Get all shortcut paths of a point.
func scanPath(p *point) []*point {
	shortPath := []*point{}
	xlimit := srcLen
	ylimit := dstLen
	for x0, y0 := p.x+1, p.y+1; x0 < xlimit && y0 < ylimit; x0, y0 = x0+1, y0+1 {
		if srcFile[x0] == dstFile[y0] {
			pn := checkNew(x0, y0, p)
			if pn != nil {
				shortPath = append(shortPath, pn)
			}
			return shortPath
		}
		for i := x0 + 1; i < xlimit; i++ {
			if srcFile[i] == dstFile[y0] {
				xlimit = i
				pi := checkNew(i, y0, p)
				if pi != nil {
					shortPath = append(shortPath, pi)
				}
				break
			}
		}
		for j := y0 + 1; j < ylimit; j++ {
			if srcFile[x0] == dstFile[j] {
				ylimit = j
				pj := checkNew(x0, j, p)
				if pj != nil {
					shortPath = append(shortPath, pj)
				}
				break
			}
		}
	}
	return shortPath
}

// Put all shortcut paths into a map.
func getPath(p *point) {
	if _, ok := path[p]; ok {
		return
	}
	ps := scanPath(p)
	if len(ps) == 0 {
		return
	}
	path[p] = ps
	for _, pn := range ps {
		getPath(pn)
	}
}

// Get the best path.
func getMostDepth() []*point {
	pList := []*point{}
	dp := 0
	var p *point
	for _, v := range newed {
		if v.depth > dp {
			dp = v.depth
			p = v
		}
	}
	// fmt.Println(p)
	if p == nil {
		return pList
	}
	var getParent func(pt *point)
	getParent = func(pt *point) {
		pList = append(pList, pt)
		if pt.depth == 0 {
			return
		}
		pt = pt.parent
		getParent(pt)
	}
	getParent(p)
	return pList
}

// Read file text.
func readFile(file string) ([]string, error) {
	fileContens := []string{}
	f, FErr := os.Open(file)
	defer f.Close()
	if FErr != nil {
		return nil, FErr
	}
	ScannerF := bufio.NewScanner(f)
	ScannerF.Split(bufio.ScanLines)
	for ScannerF.Scan() {
		fileContens = append(fileContens, ScannerF.Text())
	}
	return fileContens, nil
}

// Output difference of files.
func Diff(src string, dst string) ([]string, error) {
	var fileErr error
	srcFile, fileErr = readFile(src)
	if fileErr != nil {
		return nil, fileErr
	}
	dstFile, fileErr = readFile(dst)
	if fileErr != nil {
		return nil, fileErr
	}
	srcLen = len(srcFile)
	dstLen = len(dstFile)
	pTmp := newpoint(-1, -1)
	getPath(pTmp)
	// for k, v := range path { // ** FOR CHECK **
	// 	fmt.Printf("%v\t%v\n", k, v)
	// }
	pathPoint := getMostDepth()
	// fmt.Println(pathPoint) // ** FOR CHECK **

	result := []string{}
	var str string
	pOne := newpoint(0, 0)
	getResult := func(pOne, pPoint *point) {
		for j := pOne.x; j < pPoint.x; j++ {
			str = fmt.Sprintf("(%4d,    ) - %s", j+1, srcFile[j])
			result = append(result, str)
		}
		for j := pOne.y; j < pPoint.y; j++ {
			str = fmt.Sprintf("(    ,%4d) + %s", j+1, dstFile[j])
			result = append(result, str)
		}
	}
	for i := len(pathPoint) - 2; i >= 0; i-- {
		getResult(pOne, pathPoint[i])
		// dstFile[pathPoint[i].y] == srcFile[pathPoint[i].x]
		str = fmt.Sprintf("(%4d,%4d)   %s", pathPoint[i].x+1, pathPoint[i].y+1, srcFile[pathPoint[i].x])
		result = append(result, str)
		pOne = newpoint(pathPoint[i].x+1, pathPoint[i].y+1)
	}
	pEnd := newpoint(srcLen, dstLen)
	getResult(pOne, pEnd)
	return result, nil
}
