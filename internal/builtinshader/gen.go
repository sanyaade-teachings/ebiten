// Copyright 2024 The Ebitengine Authors
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
// limitations under the License.

//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2/internal/builtinshader"
)

func main() {
	if err := xmain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const license = `// Copyright 2024 The Ebitengine Authors
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
// limitations under the License.
`

func xmain() error {
	f, err := os.Create("defs.go")
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	if _, err := w.WriteString("// Code generated by gen.go using 'go generate'. DO NOT EDIT.\n\n"); err != nil {
		return err
	}
	if _, err := w.WriteString(license); err != nil {
		return err
	}
	if _, err := w.WriteString("\npackage builtinshader\n"); err != nil {
		return err
	}

	for _, s := range builtinshader.AppendShaderSources(nil) {
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
		if _, err := w.WriteString("//ebitengine:shader\n"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "const _ = %q\n", s); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
