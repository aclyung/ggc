package builtin

import "testing"
import "almeng.com/glang/core/builtin/operators"

func TestEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		if !operators.Equal(1, 1) {
			t.Error("1 != 1")
		}
		if operators.Equal(1, 2) {
			t.Error("1 == 2")
		}
	})
}
