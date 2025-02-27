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

func TestCmdBlueprintsSave(t *testing.T) {
	// Test the "blueprints save " command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
        {
            "description": "simple blueprint",
            "groups": [],
            "modules": [],
            "name": "simple",
            "packages": [
                {
                    "name": "bash",
                    "version": "*"
                }
            ],
            "version": "0.1.0"
        }
    ],
    "changes": [
        {
            "changed": false,
            "name": "simple"
        }
    ],
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

	cmd, out, err := root.ExecuteTest("blueprints", "save", "simple")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/simple", mc.Req.URL.Path)

	_, err = os.Stat("simple.toml")
	assert.Nil(t, err)
}

func TestCmdBlueprintsSaveUnknown(t *testing.T) {
	// Test the "blueprints save " command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
    ],
    "changes": [
    ],
    "errors": [
		{
            "id": "UnknownBlueprint",
            "msg": "test-no-bp: "
        }
	]
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

	cmd, out, err := root.ExecuteTest("blueprints", "save", "test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint: test-no-bp")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/test-no-bp", mc.Req.URL.Path)

	_, err = os.Stat("test-no-bp.toml")
	assert.NotNil(t, err)
}

func TestCmdBlueprintsSaveJSON(t *testing.T) {
	// Test the "blueprints save " command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
        {
            "description": "simple blueprint",
            "groups": [],
            "modules": [],
            "name": "simple",
            "packages": [
                {
                    "name": "bash",
                    "version": "*"
                }
            ],
            "version": "0.1.0"
        }
    ],
    "changes": [
        {
            "changed": false,
            "name": "simple"
        }
    ],
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

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "save", "simple")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"name\": \"simple\"")
	assert.Contains(t, string(stdout), "\"changed\": false")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/simple", mc.Req.URL.Path)

	_, err = os.Stat("simple.toml")
	assert.Nil(t, err)
}

func TestCmdBlueprintsSaveUnknownJSON(t *testing.T) {
	// Test the "blueprints save " command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
    ],
    "changes": [
    ],
    "errors": [
		{
            "id": "UnknownBlueprint",
            "msg": "test-no-bp: "
        }
	]
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

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "save", "test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"test-no-bp: \"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/info/test-no-bp\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/test-no-bp", mc.Req.URL.Path)

	_, err = os.Stat("test-no-bp.toml")
	assert.NotNil(t, err)
}
