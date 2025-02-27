// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdBlueprintsFreeze(t *testing.T) {
	// Test the "blueprints freeze" command
	json := `{
        "blueprints": [
		    {
                "blueprint": {
                    "description": "Install tmux",
                    "distro": "",
                    "groups": [],
                    "modules": [],
                    "name": "cli-test-bp-1",
                    "packages": [
                        {
                            "name": "tmux",
                            "version": "3.1c-2.fc34.x86_64"
                        }
                    ],
                    "version": "0.0.3"
                }
            }
        ],
        "errors": [
            {
                "id": "UnknownBlueprint",
                "msg": "test-no-bp: blueprint not found"
            }
        ]
}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "freeze", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "blueprint: cli-test-bp-1 v0.0.3")
	assert.Contains(t, string(stdout), "tmux-3.1c-2.fc34.x86_64")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint: test-no-bp: blueprint not found")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}

func TestCmdBlueprintsFreezeJSON(t *testing.T) {
	// Test the "blueprints freeze" command
	json := `{
        "blueprints": [
		    {
                "blueprint": {
                    "description": "Install tmux",
                    "distro": "",
                    "groups": [],
                    "modules": [],
                    "name": "cli-test-bp-1",
                    "packages": [
                        {
                            "name": "tmux",
                            "version": "3.1c-2.fc34.x86_64"
                        }
                    ],
                    "version": "0.0.3"
                }
            }
        ],
        "errors": [
            {
                "id": "UnknownBlueprint",
                "msg": "test-no-bp: blueprint not found"
            }
        ]
}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "freeze", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"version\": \"3.1c-2.fc34.x86_64\"")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"test-no-bp: blueprint not found\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}

func TestCmdBlueprintsFreezeSave(t *testing.T) {
	// Test the "blueprints freeze save" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
        {
			"blueprint": {
				"description": "Install tmux",
				"distro": "",
				"groups": [],
				"modules": [],
				"name": "cli-test-bp-1",
				"packages": [
					{
						"name": "tmux",
						"version": "3.1c-2.fc34.x86_64"
					}
				],
				"version": "0.0.3"
			}
	   }],
    "errors": []
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	dir, err := ioutil.TempDir("", "test-bp-save-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	cmd, out, err := root.ExecuteTest("blueprints", "freeze", "save", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeSaveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)

	_, err = os.Stat("cli-test-bp-1.frozen.toml")
	assert.Nil(t, err)
}

func TestCmdBlueprintsFreezeSaveJSON(t *testing.T) {
	// Test the "blueprints freeze save" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
        {
			"blueprint": {
				"description": "Install tmux",
				"distro": "",
				"groups": [],
				"modules": [],
				"name": "cli-test-bp-1",
				"packages": [
					{
						"name": "tmux",
						"version": "3.1c-2.fc34.x86_64"
					}
				],
				"version": "0.0.3"
			}
	   }],
    "errors": []
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	dir, err := ioutil.TempDir("", "test-bp-save-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "freeze", "save", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeSaveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"version\": \"3.1c-2.fc34.x86_64\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)

	_, err = os.Stat("cli-test-bp-1.frozen.toml")
	assert.Nil(t, err)
}

func TestCmdBlueprintsFreezeShow(t *testing.T) {
	// Test the "blueprints freeze show" command
	toml := `name = "cli-test-bp-1"
description = "Install tmux"
version = "0.0.3"
modules = []
groups = []
distro = ""

[[packages]]
name = "tmux"
version = "3.1c-2.fc34.x86_64"
}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(toml))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "freeze", "show", "cli-test-bp-1")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeShowCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "name = \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "[[packages]]")
	assert.Contains(t, string(stdout), "version = \"3.1c-2.fc34.x86_64\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)
}

func TestCmdBlueprintsFreezeShowJSON(t *testing.T) {
	// Test the "blueprints freeze show" command
	json := `{
        "blueprints": [
		    {
                "blueprint": {
                    "description": "Install tmux",
                    "distro": "",
                    "groups": [],
                    "modules": [],
                    "name": "cli-test-bp-1",
                    "packages": [
                        {
                            "name": "tmux",
                            "version": "3.1c-2.fc34.x86_64"
                        }
                    ],
                    "version": "0.0.3"
                }
            }
        ],
        "errors": [
        ]
}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "freeze", "show", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeShowCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"version\": \"3.1c-2.fc34.x86_64\"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/freeze/cli-test-bp-1\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)
}
