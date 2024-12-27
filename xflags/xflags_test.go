// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xflags

import (
	"fmt"
	"os"
	"testing"

	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/util"
)

// funcName returns the caller's function name
func funcName() string {
	_, n := util.CallerFuncName(1)
	return n
}

type options struct {
	// these all use the default mapping
	String string
	Int    int
	Uint   uint
	Uint64 uint64
	Int64  int64
	// and these use custom mappings ...
	String0 string   `flag:"s0"`
	String1 string   `flag:"s1"`
	String2 string   `flag:"s2|string2"`
	String3 string   `flag:"s3,Help with s3,Wow"`
	String4 string   `flag:"s4,Help with s4"`
	Bool0   bool     `flag:"b0"`
	Bool1   bool     `flag:"b1"`
	Bool2   bool     `flag:"b2,Help with b2,true"`
	Uint0   uint     `flag:"u0"`
	Uint1   uint     `flag:"u1"`
	Uint2   uint     `flag:"u2,Help with u2,2112"`
	Uint3   uint     `flag:"u3,Help with u3"`
	Args    []string `flag:"*,extra args"` // remaining args
	Ignore1 string   `flag:"-"`            // ignore this one
	ignore2 string   // non-exported fields are ignored
}

func TestUsage(t *testing.T) {
	err := Usage(os.Stdout, funcName(), &(options{}))
	assert.Nil(t, err)
}

func TestFlags(t *testing.T) {
	opts := options{
		String4: "Hello, world!",
		Bool0:   true,
		Uint3:   2113,
		Ignore1: "Ignore1",
		ignore2: "ignore2",
	}
	args := []string{
		"--s0=Hello-0",
		"--s1=Hello-1",
		"--s2=World",
		"--b1=true",
		"--u1=42",
		"Hello",
		"World",
	}
	err := ParseArgs("TestInitFlags", &opts, args)
	assert.Nil(t, err)

	fmt.Fprintf(os.Stderr, "opts: %+v\n", opts)

	// fmt.Fprintf(os.Stderr, "flags: %+v\n", flags)
	assert.Equal(t, "Hello-0", opts.String0)
	assert.Equal(t, "Hello-1", opts.String1)
	assert.Equal(t, "World", opts.String2)
	assert.Equal(t, "Wow", opts.String3)
	assert.Equal(t, "Hello, world!", opts.String4)

	assert.True(t, opts.Bool1)
	assert.True(t, opts.Bool2)

	assert.Equal(t, 0, opts.Uint0)
	assert.Equal(t, 42, opts.Uint1)
	assert.Equal(t, 2112, opts.Uint2)
	assert.Equal(t, 2113, opts.Uint3)

	assert.Equal(t, 2, len(opts.Args))
	assert.Equal(t, "Ignore1", opts.Ignore1)
	assert.Equal(t, "ignore2", opts.ignore2)
	assert.DeepEqual(t, []string{"Hello", "World"}, opts.Args)
}

func TestFlagsWithBadUsage(t *testing.T) {
	err := ParseArgs(funcName(), &(options{}), []string{"--bogus"})
	assert.NotNil(t, err)
	t.Log(err)
}

type badOptions struct {
	Func func() // not supported
}

func TestUsageWithBadOptions(t *testing.T) {
	err := Usage(os.Stdout, funcName(), &(badOptions{}))
	assert.NotNil(t, err)
	t.Log(err)
}
func TestFlagsWithBadOptions(t *testing.T) {
	err := ParseArgs(funcName(), &(badOptions{}), []string{"--debug"})
	assert.NotNil(t, err)
	t.Log(err)
}

func TestFlagsWithWrongType(t *testing.T) {
	err := ParseArgs(funcName(), options{}, []string{"--debug"})
	assert.NotNil(t, err)
	t.Log(err)

	notAStruct := 2112
	err = ParseArgs(funcName(), &notAStruct, []string{"--debug"})
	assert.NotNil(t, err)
	t.Log(err)

}

type badRemainingArgsType struct {
	Args []bool `flag:"*"`
}

func TestFlagsWithBadRemainingArgsType(t *testing.T) {
	err := ParseArgs(funcName(), &badRemainingArgsType{}, []string{"--debug"})
	assert.NotNil(t, err)
	t.Log(err)
}

type MyOptions struct {
	Verbose  bool   // use defaults
	Debug    bool   `flag:"d|debug,Enable debug output"`
	Port     string `flag:"p|port,Port to listen on,8080"`
	Retries  int    `flag:",Number of times to try,10"`
	internal int    `flag:"-"` // ignore this one
}

func TestExample(t *testing.T) {
	opts := MyOptions{
		internal: 2112,
	}
	err := ParseArgs("MyProgram", &opts, []string{"--debug"})
	assert.Nil(t, err)
	assert.True(t, opts.Debug)
	assert.Equal(t, "8080", opts.Port)
	assert.Equal(t, 2112, opts.internal)
}
