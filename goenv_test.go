package goenv_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/skipor/goenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigpanic(t *testing.T) {
	t.Parallel()
	sp := goenv.Sigpanic()
	assert.Equal(t, "runtime.sigpanic", sp.Name())
}

func TestGoRoot(t *testing.T) {
	expectedGoRootSrc := strings.TrimSuffix(gorootIoGoFile, ioGoFileTrimmed)
	assert.Equal(t, expectedGoRootSrc, goenv.GoRootSrc())
	assert.Equal(t, strings.TrimSuffix(expectedGoRootSrc, "src/"), goenv.GoRoot())
	assert.True(t, goenv.InGoroot(gorootIoGoFile))
}

func TestGoPath(t *testing.T) {
	var file string
	_, file, _, _ = runtime.Caller(0)
	require.True(t, strings.HasSuffix(file, thisFileTrimmed), file)
	expectedGoPathSrc := strings.TrimSuffix(file, thisFileTrimmed)
	assert.Equal(t, expectedGoPathSrc, goenv.GoPathSrc())
	assert.Equal(t, strings.TrimSuffix(expectedGoPathSrc, "src/"), goenv.GoPath())
	assert.False(t, goenv.InGoroot(file))
}

const ioGoFileTrimmed = "io/io.go"
const thisFileTrimmed = "github.com/skipor/goenv/goenv_test.go"

// Another way to get file in $GOROOT.
var gorootIoGoFile = func() string {
	var w callerWriter
	io.WriteString(&w, "test")
	const ioFileSuffix = "io/io.go"
	if !strings.HasSuffix(w.callerFile, ioFileSuffix) {
		panic(fmt.Sprintf("Expect get $GOROOT%s, but get %s", ioFileSuffix, w.callerFile))
	}
	return w.callerFile
}()

type callerWriter struct {
	callerFile string
}

func (w *callerWriter) Write(p []byte) (int, error) {
	var ok bool
	_, w.callerFile, _, ok = runtime.Caller(1)
	if !ok {
		panic("undefined caller")
	}
	return ioutil.Discard.Write(p)
}
