// MIT License
//
// Copyright (c) 2017 Iman Tumorang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package faker

import (
	"io"
	mathrand "math/rand"
	"sync"
)

var (
	rand   *mathrand.Rand
	crypto io.Reader
)

type safeSource struct {
	mx sync.Mutex
	mathrand.Source
}

func (s *safeSource) Int63() int64 {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.Source.Int63()
}

// NewSafeSource wraps an unsafe rand.Source with a mutex to guard the random source
// against concurrent access.
func NewSafeSource(in mathrand.Source) mathrand.Source {
	return &safeSource{
		Source: in,
	}
}

// SetRandomSource sets a new random source at the package level.
//
// To use a concurrent-safe source, you may wrap it with NewSafeSource,
// e.g. SetRandomSource(NewSafeSource(mysource)).
//
// The default is the global, concurrent-safe source provided by math/rand.
func SetRandomSource(in mathrand.Source) {
	rand = mathrand.New(in)
}

// SetCryptoSource sets a new reader for functions using a cryptographically-safe random generator (e.g. UUID).
//
// The default is the global source provided by crypto/rand.
func SetCryptoSource(in io.Reader) {
	crypto = in
}
