// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import "testing"

func TestContext(t *testing.T) {
	c := Context{}

	// Internal storage still nil.
	if c.Exists("key") {
		t.Error("The key 'key' should not exist in context")
	}

	if c.Get("key") != nil {
		t.Error("Expected nil value for 'key'")
	}

	if c.store != nil {
		t.Error("Expected the internal store to be nil")
	}

	c.Set("key", "value")

	// Internal storage created.
	if !c.Exists("key") {
		t.Error("Expected 'key' to exist")
	}

	if c.Get("key").(string) != "value" {
		t.Error("Expected 'key' to be 'value'")
	}

	if c.Get("key2") != nil {
		t.Error("Expected nil value for 'key2'")
	}
}
