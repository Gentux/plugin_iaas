package libnan

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ===================================================================================================

var ()

func CopyFile(_sourcePath, _destPath string) error {

	sourceBytes, err := ioutil.ReadFile(_sourcePath)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(_destPath, sourceBytes, 0777)
}

func ReplaceInFile(_filePath, _oldString, _newString string) error {

	var perm os.FileMode

	if fi, e := os.Stat(_filePath); e != nil {
		return e
	} else {
		perm = fi.Mode()
	}

	sourceBytes, err := ioutil.ReadFile(_filePath)
	if err != nil {
		return err
	}

	sNewContents := ""
	sFileLines := strings.Split(string(sourceBytes), "\n")

	for _, sLine := range sFileLines {
		if idx := strings.Index(sLine, _oldString); idx >= 0 {
			sLine = strings.Replace(sLine, _oldString, _newString, -1)
		}
		sNewContents += fmt.Sprintf("%s\n", sLine)
	}

	return ioutil.WriteFile(_filePath, []byte(sNewContents), perm)
}
