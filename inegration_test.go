package main

import (
	"testing"

	"github.com/eliona-smart-building-assistant/app-integration-tests/app"
	"github.com/eliona-smart-building-assistant/app-integration-tests/assert"
	"github.com/eliona-smart-building-assistant/app-integration-tests/test"
)

func TestApp(t *testing.T) {
	app.StartApp()
	test.AppWorks(t)
	t.Run("TestAssetTypes", assetTypes)
	t.Run("TestSchema", schema)
	app.StopApp()
}

func assetTypes(t *testing.T) {
	t.Parallel()

	assert.AssetTypeExists(t, "coffeecloud_group", []string{})
	assert.AssetTypeExists(t, "coffeecloud_machine", []string{})
	assert.AssetTypeExists(t, "coffeecloud_root", []string{})
}

func schema(t *testing.T) {
	t.Parallel()

	assert.SchemaExists(t, "coffeecloud", []string{"asset", "configuration"})
}
