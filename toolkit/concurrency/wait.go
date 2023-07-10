// Copyright (c) 2020 StackRox Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package concurrency

import "time"

// Wait waits indefinitely until the condition represented by the given Waitable is fulfilled.
func Wait(w Waitable) {
	<-w.Done()
}

// IsDone checks if the given waitable's condition is fulfilled.
func IsDone(w Waitable) bool {
	select {
	case <-w.Done():
		return true
	default:
		return false
	}
}

// WaitWithTimeout waits for the given Waitable with a specified timeout. It returns false if the timeout expired
// before the condition was fulfilled, true otherwise.
func WaitWithTimeout(w Waitable, timeout time.Duration) bool {
	if timeout <= 0 {
		return IsDone(w)
	}

	t := time.NewTimer(timeout)
	select {
	case <-w.Done():
		if !t.Stop() {
			<-t.C
		}
		return true
	case <-t.C:
		return false
	}
}
