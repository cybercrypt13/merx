/*
Purpose:	This file contains all of our shared objects and functions that are
		used across the entire merx implementation
Written:	05.28.2013
By:		Noel Jacques <njacques@nizex.com>
URL:		www.nizex.com

The MIT License (MIT)

Copyright (c) 2013 Nizex Inc.

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
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//this is our main interface that we use to reference the beginning communication
//packages coming through.  We build off of this to make each of our individual
//communication threads.
type Packagecontents interface {
	ParsePackage([]byte, *sql.DB, int) error
	VerifyPackage() error
	ProcessPackage(int, string) ([]byte, error)
}

type AdminPackage interface {
	ParsePackage([]byte, *sql.DB) (int, error)
	ProcessPackage() (int, error)
}
