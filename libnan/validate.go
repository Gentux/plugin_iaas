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
	"net/http"
	"regexp"
	"strings"
)

var ()

// Check if the name meets the following requirements:
//     at least 1 and less than 65 characters long
//     characters that can be used:
//        any alphanumeric character 0 to 9 OR A to Z or a to z
//        punctuation symbols . , " ' ? ! ; : # $ % & ( ) * + - / < > = @ [ ] \ ^ _ { } | ~
func ValidName(val string) bool {

	if len(val) < 1 || len(val) >= 65 || strings.Contains(val, " ") {
		return false
	}

	bMatch, e := regexp.MatchString(`[[:lower:][:upper:][:punct:]\d]+`, val)
	if e != nil {
		Log("Error when attempting to match a regexp", e)
	}
	if e != nil || bMatch == false {
		return false
	}

	return true
}

// Check if this is an email address
func ValidEmail(address string) bool {

	matched, err := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, address)
	//"^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,4}$", adress)

	if err != nil {
		LogError("Error when attempting to regexp.MatchString :", err)
	}

	if !matched {
		LogError("Failed to match conformation rules for email %s", address)
	}

	return matched
}

// Check if the password meets the following requirements:
//    at least 7 and less than 65 characters long
//    has at least one digit
//    has at least one Upper case Alphabet
//    has at least one Lower case Alphabet
//       characters that can be used:
//           any alphanumeric character 0 to 9 OR A to Z or a to z
//           punctuation symbols . , " ' ? ! ; : # $ % & ( ) * + - / < > = @ [ ] \ ^ _ { } | ~
func ValidPassword(val string) bool {

	// posixRegExp := `[[ ${#s} -ge 7 && ${#s} -le 64 && "$s" == *[[:upper:]]* && "$s" == *[[:lower:]]* && "$s" == *[[:digit:]]* && "$s" =~ ^[[:alnum:][:punct:]]+$ ]]`

	// pRegExp, e := regexp.CompilePOSIX(posixRegExp)
	// if e != nil {
	// 	LogError("Failed to generate regexp parser")
	// 	return false
	// }

	// if pRegExp.Match([]byte(val)) == false {
	// 	LogError("Couldn't match regexp")
	// 	return false
	// }

	// return true

	if len(val) < 7 || len(val) >= 65 || strings.Contains(val, " ") {
		return false
	}

	if !strings.ContainsAny(val, "0123456789") {
		return false
	}

	if bHasUpper, e := regexp.MatchString(`[[:upper:]]+`, val); e != nil {
		LogError("Error when attempting to match a regexp", e)
		return false
	} else if !bHasUpper {
		LogError("Couldn't match regexp for presence of at one uppercase char")
		return false
	}

	bHasLower, e := regexp.MatchString(`[[:lower:]]+`, val)
	if e != nil {
		LogError("Error when attempting to match regexp: ", e)
		return false
	} else if !bHasLower {
		LogError("Couldn't match regexp for presence of at one lowercase char")
		return false
	}

	if bHasNonGraphChars, e := regexp.MatchString(`[[:^graph:]]+`, val); e != nil {
		LogError("Error when attempting to match a regexp", e)
	} else if bHasNonGraphChars {
		LogError("Couldn't match regexp for detection of non graphical characters")
		return false
	}

	return true
}

func ValidUrl(url string) bool {

	resp, e := http.Get(url)
	defer resp.Body.Close()

	if e != nil {
		Log("Error when getting url :", url)
		return false
	}

	return resp.StatusCode == 200
}
