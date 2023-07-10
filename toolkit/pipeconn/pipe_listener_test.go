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

package pipeconn

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func u32BE(v uint32) []byte {
	var enc [4]byte
	binary.BigEndian.PutUint32(enc[:], v)
	return enc[:]
}

func TestNetwork(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Network, pipeAddr.Network())
}

func TestPipeListener_Connections(t *testing.T) {
	t.Parallel()

	lis, dialCtx := NewPipeListener()

	var clientSum uint32
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			conn, err := dialCtx(context.Background())
			require.NoError(t, err)

			sum := uint32(idx)
			_, err = conn.Write(u32BE(sum))
			assert.NoError(t, err)

			var buf [4]byte
			_, err = io.ReadFull(conn, buf[:])
			assert.NoError(t, err)
			sum += binary.BigEndian.Uint32(buf[:])

			_, err = conn.Write(u32BE(sum))
			assert.NoError(t, err)

			assert.NoError(t, conn.Close())

			atomic.AddUint32(&clientSum, sum)
		}(i)
	}

	var serverSum uint32
	for i := 0; i < 10; i++ {
		conn, err := lis.Accept()
		require.NoError(t, err)

		wg.Add(1)
		go func(idx int, conn net.Conn) {
			defer wg.Done()

			sum := uint32(idx)

			var buf [4]byte
			_, err := io.ReadFull(conn, buf[:])
			assert.NoError(t, err)
			sum += binary.BigEndian.Uint32(buf[:])

			_, err = conn.Write(u32BE(uint32(idx)))
			assert.NoError(t, err)

			_, err = io.ReadFull(conn, buf[:])
			assert.NoError(t, err)
			assert.Equal(t, sum, binary.BigEndian.Uint32(buf[:]))

			atomic.AddUint32(&serverSum, sum)

			n, err := io.ReadFull(conn, buf[:])
			assert.Zero(t, n)
			assert.Equal(t, io.EOF, err)
		}(i, conn)
	}

	wg.Wait()
	assert.Equal(t, serverSum, clientSum)
}

func TestPipeListener_Close(t *testing.T) {
	t.Parallel()

	lis, dialCtx := NewPipeListener()

	assert.NoError(t, lis.Close())

	conn, err := lis.Accept()
	assert.Nil(t, conn)
	assert.Equal(t, ErrClosed, err)

	conn, err = dialCtx(context.Background())
	assert.Nil(t, conn)
	assert.Equal(t, ErrClosed, err)

	assert.Equal(t, ErrAlreadyClosed, lis.Close())
}
