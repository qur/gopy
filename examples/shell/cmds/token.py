# Copyright 2011 Julian Phillips.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import sh

def token(*args):
    for arg in args:
        print "%s -> %s" % (arg, sh.tokenise(arg))
