/*
Purpose:	Readconf is a package for reading config files and returning the specified
			default value if the parameter could not be found.
Written:	06.05.2014
By:		Noel Jacques <njacques@nizex.com>
URL:		www.nizex.com

The MIT License (MIT)

Copyright (c) 2014 Nizex Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package common

import (
	"strings"

	"github.com/dlintw/goconf"
)

var ConfigFile *goconf.ConfigFile

func OpenConfigFile(filename string) (err error) {
	//06.05.2014 naj - first read the config file
	ConfigFile, err = goconf.ReadConfigFile(filename)
	return
}

//06.05.2014 naj - This will attempt to return a value from the config file.
//If the value cannot be found or is empty then the defaultval is returned
func GetConfig(section string, param string, defaultval string) (value string) {
	value, err := ConfigFile.GetString(section, param)

	if err != nil {
		value = defaultval
	}

	value = strings.TrimSpace(value)
	if value == "" {
		value = defaultval
	}

	return
}
