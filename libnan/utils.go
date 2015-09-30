/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

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
