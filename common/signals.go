/*
Purpose:	Signals is a common package used to monitor for SIGTERM, SIGINT, and SIGHUP events
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
	"os"
	"os/signal"
	"syscall"
)

//06.05.2014 naj - this function will monitor the os for signal events
func SignalCatcher(servicename string, logging *LogFile) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for sig := range ch {
		if sig == syscall.SIGHUP {
			//05.28.2013 naj - start a new log file.
			logging.RestartLog()
		} else {
			//05.28.2013 naj - terminate the service.
			logging.LogError("", servicename+" Server Shutting Down")
			fmt.Println(servicename + " Server Shutting Down")
			os.Exit(0)
		}
	}
}
