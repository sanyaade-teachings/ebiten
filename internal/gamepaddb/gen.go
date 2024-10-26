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
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

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

// gamecontrollerdb.txt is downloaded at https://github.com/mdqinc/SDL_GameControllerDB.

// To update the database file, run:
//
//     curl --location --remote-name https://raw.githubusercontent.com/mdqinc/SDL_GameControllerDB/master/gamecontrollerdb.txt

//go:embed gamecontrollerdb.txt
var gameControllerDB []byte

const dbTemplate string = `{{.License}}

{{.DoNotEdit}}

{{.BuildConstraints}}

package gamepaddb

import (
	_ "embed"
)

//go:embed gamecontrollerdb_{{.FileNameSuffix}}.txt
var controllerBytes []byte

{{if .HasGLFWGamepads}}
var additionalGLFWGamepads = []byte(` + "`" + `
78696e70757401000000000000000000,XInput Gamepad (GLFW),platform:Windows,a:b0,b:b1,x:b2,y:b3,leftshoulder:b4,rightshoulder:b5,back:b6,start:b7,leftstick:b8,rightstick:b9,leftx:a0,lefty:a1,rightx:a2,righty:a3,lefttrigger:a4,righttrigger:a5,dpup:h0.1,dpright:h0.2,dpdown:h0.4,dpleft:h0.8,
78696e70757402000000000000000000,XInput Wheel (GLFW),platform:Windows,a:b0,b:b1,x:b2,y:b3,leftshoulder:b4,rightshoulder:b5,back:b6,start:b7,leftstick:b8,rightstick:b9,leftx:a0,lefty:a1,rightx:a2,righty:a3,lefttrigger:a4,righttrigger:a5,dpup:h0.1,dpright:h0.2,dpdown:h0.4,dpleft:h0.8,
78696e70757403000000000000000000,XInput Arcade Stick (GLFW),platform:Windows,a:b0,b:b1,x:b2,y:b3,leftshoulder:b4,rightshoulder:b5,back:b6,start:b7,leftstick:b8,rightstick:b9,leftx:a0,lefty:a1,rightx:a2,righty:a3,lefttrigger:a4,righttrigger:a5,dpup:h0.1,dpright:h0.2,dpdown:h0.4,dpleft:h0.8,
78696e70757404000000000000000000,XInput Flight Stick (GLFW),platform:Windows,a:b0,b:b1,x:b2,y:b3,leftshoulder:b4,rightshoulder:b5,back:b6,start:b7,leftstick:b8,rightstick:b9,leftx:a0,lefty:a1,rightx:a2,righty:a3,lefttrigger:a4,righttrigger:a5,dpup:h0.1,dpright:h0.2,dpdown:h0.4,dpleft:h0.8,
78696e70757405000000000000000000,XInput Dance Pad (GLFW),platform:Windows,a:b0,b:b1,x:b2,y:b3,leftshoulder:b4,rightshoulder:b5,back:b6,start:b7,leftstick:b8,rightstick:b9,leftx:a0,lefty:a1,rightx:a2,righty:a3,lefttrigger:a4,righttrigger:a5,dpup:h0.1,dpright:h0.2,dpdown:h0.4,dpleft:h0.8,
78696e70757406000000000000000000,XInput Guitar (GLFW),platform:Windows,a:b0,b:b1,x:b2,y:b3,leftshoulder:b4,rightshoulder:b5,back:b6,start:b7,leftstick:b8,rightstick:b9,leftx:a0,lefty:a1,rightx:a2,righty:a3,lefttrigger:a4,righttrigger:a5,dpup:h0.1,dpright:h0.2,dpdown:h0.4,dpleft:h0.8,
78696e70757408000000000000000000,XInput Drum Kit (GLFW),platform:Windows,a:b0,b:b1,x:b2,y:b3,leftshoulder:b4,rightshoulder:b5,back:b6,start:b7,leftstick:b8,rightstick:b9,leftx:a0,lefty:a1,rightx:a2,righty:a3,lefttrigger:a4,righttrigger:a5,dpup:h0.1,dpright:h0.2,dpdown:h0.4,dpleft:h0.8,
` + "`" + `)
{{end}}

func init() {
	if err := Update(controllerBytes); err != nil {
		panic(err)
	}{{if .HasGLFWGamepads}}
	if err := Update(additionalGLFWGamepads); err != nil {
		panic(err)
	}{{end}}
}
`

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// doNotEditor is a special comment for generated files.
	// This follows the standard comment rule (https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source).
	const doNotEdit = "// Code generated by gen.go using 'go generate'. DO NOT EDIT."

	type gamePadPlatform struct {
		filenameSuffix   string
		buildConstraints string
		hasGLFWGamepads  bool
	}

	platforms := map[string]gamePadPlatform{
		"Windows": {
			filenameSuffix:   "windows",
			buildConstraints: "//go:build !microsoftgdk",
			hasGLFWGamepads:  true,
		},
		"Mac OS X": {
			filenameSuffix:   "macos",
			buildConstraints: "//go:build darwin && !ios",
		},
		"Linux": {
			filenameSuffix:   "linbsd",
			buildConstraints: "//go:build (freebsd || (linux && !android) || netbsd || openbsd) && !nintendosdk && !playstation5",
		},
		"iOS": {
			filenameSuffix: "ios",
		},
		"Android": {
			filenameSuffix: "android",
		},
	}

	controllerDBs, err := splitDBsByPlatform(gameControllerDB)
	if err != nil {
		return err
	}

	for sdlPlatformName, platform := range platforms {
		controllerDB, ok := controllerDBs[sdlPlatformName]
		if !ok {
			return fmt.Errorf("failed to find controller db for platform %s in gamecontrollerdb.txt", sdlPlatformName)
		}

		// Write each chunk into separate text file for embedding into respective generated files.
		if err = os.WriteFile(fmt.Sprintf("gamecontrollerdb_%s.txt", platform.filenameSuffix), []byte(controllerDB), 0666); err != nil {
			return err
		}

		path := fmt.Sprintf("db_%s.go", platform.filenameSuffix)
		tmpl, err := template.New(path).Parse(dbTemplate)
		if err != nil {
			return err
		}

		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		if err := tmpl.Execute(w, struct {
			License          string
			DoNotEdit        string
			BuildConstraints string
			FileNameSuffix   string
			HasGLFWGamepads  bool
		}{
			License:          license,
			DoNotEdit:        doNotEdit,
			BuildConstraints: platform.buildConstraints,
			FileNameSuffix:   platform.filenameSuffix,
			HasGLFWGamepads:  platform.hasGLFWGamepads,
		}); err != nil {
			return err
		}
		if err := w.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func splitDBsByPlatform(controllerDB []byte) (map[string]string, error) {
	s := bufio.NewScanner(bytes.NewReader(controllerDB))
	dbs := map[string]string{}

	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		for _, part := range strings.Split(line, ",") {
			if platform, found := strings.CutPrefix(part, "platform:"); found {
				dbs[platform] += line + "\n"
				break
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return dbs, nil
}
