/*
Purpose:	This file handles writing to our log file.
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
	"fmt"
	"log"
	"os"
)

type LogFile struct {
	FileName string
	out      *os.File
	logger   *log.Logger
}

func NewLog(fname string) (file *LogFile, err error) {
	//05.28.2013 naj - open a new log file
	out, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	//05.28.2013 naj - link the logging service to the new file.
	l := log.New(out, "", log.Ldate+log.Ltime)

	return &LogFile{FileName: fname, out: out, logger: l}, nil
}

func (lf *LogFile) LogError(errcode string, errtext string) {
	lf.logger.Print(errcode, errtext)

	if errcode == "" {
		fmt.Println(errtext)
	} else {
		fmt.Println("Error code: " + errcode + " - " + errtext)
	}
}

func (lf *LogFile) RestartLog() error {
	//05.28.2013 naj - make a log entry in the old file.
	lf.logger.Print("", "Closing Log File")

	//05.28.2013 naj - close the old log file
	err := lf.out.Close()
	if err != nil {
		return err
	}

	//05.28.2013 naj - open a new log file
	lf.out, err = os.OpenFile(lf.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	//05.28.2013 naj - link the logging service to the new file.
	lf.logger = log.New(lf.out, "", log.Ldate+log.Ltime)

	//05.28.2013 naj - make a log entry in the new file.
	lf.logger.Print("", "Starting New Log File")

	return nil
}
