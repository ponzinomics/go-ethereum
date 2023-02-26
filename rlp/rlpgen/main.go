// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/types"
	"os"

	"golang.org/x/tools/go/packages"
)

const pathOfPackageRLP = "github.com/ponzinomics/go-ethereum/rlp"

func main() {
	var (
		pkgdir     = flag.String("dir", ".", "input package")
		output     = flag.String("out", "-", "output file (default is stdout)")
		genEncoder = flag.Bool("encoder", true, "generate EncodeRLP?")
		genDecoder = flag.Bool("decoder", false, "generate DecodeRLP?")
		typename   = flag.String("type", "", "type to generate methods for")
	)
	flag.Parse()

	cfg := Config{
		Dir:             *pkgdir,
		Type:            *typename,
		GenerateEncoder: *genEncoder,
		GenerateDecoder: *genDecoder,
	}
	code, err := cfg.process()
	if err != nil {
		fatal(err)
	}
	if *output == "-" {
		os.Stdout.Write(code)
	} else if err := os.WriteFile(*output, code, 0600); err != nil {
		fatal(err)
	}
}

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

type Config struct {
	Dir  string // input package directory
	Type string

	GenerateEncoder bool
	GenerateDecoder bool
}

// process generates the Go code.
func (cfg *Config) process() (code []byte, err error) {
	// Load packages.
	pcfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedDeps,
		Dir:        cfg.Dir,
		BuildFlags: []string{"-tags", "norlpgen"},
	}
	ps, err := packages.Load(pcfg, pathOfPackageRLP, ".")
	if err != nil {
		return nil, err
	}
	if len(ps) == 0 {
		return nil, fmt.Errorf("no Go package found in %s", cfg.Dir)
	}
	packages.PrintErrors(ps)

	// Find the packages that were loaded.
	var (
		pkg        *types.Package
		packageRLP *types.Package
	)
	for _, p := range ps {
		if len(p.Errors) > 0 {
			return nil, fmt.Errorf("package %s has errors", p.PkgPath)
		}
		if p.PkgPath == pathOfPackageRLP {
			packageRLP = p.Types
		} else {
			pkg = p.Types
		}
	}
	bctx := newBuildContext(packageRLP)

	// Find the type and generate.
	typ, err := lookupStructType(pkg.Scope(), cfg.Type)
	if err != nil {
		return nil, fmt.Errorf("can't find %s in %s: %v", cfg.Type, pkg, err)
	}
	code, err = bctx.generate(typ, cfg.GenerateEncoder, cfg.GenerateDecoder)
	if err != nil {
		return nil, err
	}

	// Add build comments.
	// This is done here to avoid processing these lines with gofmt.
	var header bytes.Buffer
	fmt.Fprint(&header, "// Code generated by rlpgen. DO NOT EDIT.\n\n")
	fmt.Fprint(&header, "//go:build !norlpgen\n")
	fmt.Fprint(&header, "// +build !norlpgen\n\n")
	return append(header.Bytes(), code...), nil
}

func lookupStructType(scope *types.Scope, name string) (*types.Named, error) {
	typ, err := lookupType(scope, name)
	if err != nil {
		return nil, err
	}
	_, ok := typ.Underlying().(*types.Struct)
	if !ok {
		return nil, errors.New("not a struct type")
	}
	return typ, nil
}

func lookupType(scope *types.Scope, name string) (*types.Named, error) {
	obj := scope.Lookup(name)
	if obj == nil {
		return nil, errors.New("no such identifier")
	}
	typ, ok := obj.(*types.TypeName)
	if !ok {
		return nil, errors.New("not a type")
	}
	return typ.Type().(*types.Named), nil
}
