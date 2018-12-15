# Stream package

[![Build Status](https://travis-ci.org/Protoc0d/stream.svg?branch=master)](https://travis-ci.org/Protoc0d/stream)

> This package is a fork of interrools package which is a translation of the python `itertools` module. It includes all the usual suspects except for the cartesian product and permutation operators. All iterators are `chan interface{}` which allows some type ambiguity for these generic functions. It would be completely ok, however, to reproduce these functions in your package for your type-specific iterators such as `chan MyStruct`. I did this mostly as a thought exercise on converting python generators to Go.


# License

Copyright (c) 2013 Jon Eisen

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
