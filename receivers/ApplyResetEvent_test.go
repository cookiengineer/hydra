package receivers

import "testing"

func TestApplyResetEvent_NoopWithNilBridge(t *testing.T) {

	ApplyResetEvent(nil)

}
