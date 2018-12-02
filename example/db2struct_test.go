package example_test

import (
	"testing"

	"git.code.oa.com/fip-team/fiorm/db2struct"
)

func TestBuild(t *testing.T) {
	db2struct.Build("user", "mypackage")
}
