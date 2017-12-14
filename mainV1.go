// diff package usage example.
package main

import (
	"fmt"
	"github.com/jacenr/filediff/diff"
)

func main() {
	var srcFile = "testFile/SrcFile"
	var dstFile = "testFile/DstFile"
	result, diffErr := diff.Diff(dstFile, srcFile)
	if diffErr != nil {
		fmt.Println(diffErr)
	}
	for _, byList := range result {
		fmt.Printf("%s\n", byList)
	}
}