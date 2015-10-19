// Copyright 2015 Datajin Technologies, Inc. All rights reserved.

package commands

import "github.com/codegangsta/cli"

// This software is provided "AS IS" but the service it connects to has its own
// terms and polcies.
const notices = `
Use of the mktmpio service is subject to the the following:
 * mktmpio Privacy Policy
   https://mktmp.io/privacy-policy
 * mktmpio Terms of Service
   https://mktmp.io/terms-of-service
`

// The dependencies list below is manually populated using the output from:
// go list -f {{.Deps}} ./...
const thirdparty = `
This binary includes 3rd party software:

 * github.com/codegangsta/cli
   MIT license
   Copyright (C) 2013 Jeremy Saenz

 * github.com/mitchellh/go-homedir
   MIT license
   Copyright (c) 2013 Mitchell Hashimoto

 * github.com/mktmpio/go-mktmpio
   Artistic-2.0 license
   Copyright (c) 2015 Datajin Technologies, Inc.

 * golang.org/x/crypto/ssh/terminal
   golang.org/x/net/websocket
   BSD-2 license
   Copyright (c) 2009 The Go Authors. All rights reserved.

 * gopkg.in/yaml.v2
   LGPL3 license
   Copyright (c) 2011-2014 - Canonical Inc.
`

// This software is provided "AS IS"
const warranties = `
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

// Contact info
const contact = `
For more information contact Datajin Technologies, Inc.:
 w: https://mktmp.io
 e: mktmpio@datajin.com
`

// Definition for the 'mktmpio config' command
var LegalCommand = cli.Command{
	Name:   "legal",
	Usage:  "display licensing, copyright, and warranty notices",
	Action: licenseAction,
}

func licenseAction(c *cli.Context) {
	c.App.Writer.Write([]byte("\nNOTICES:\n" + notices))
	c.App.Writer.Write([]byte("\nWARRANTIES:\n" + warranties))
	c.App.Writer.Write([]byte("\n3RD PARTY SOFTWARE:\n" + thirdparty))
	c.App.Writer.Write([]byte("\nCONTACT:\n" + contact))
}
