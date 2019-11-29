package configer

import (
	"testing"
)

func TestNew(t *testing.T) {
	settings := New("config.ini")
	suppliers := settings.GetSection("suppliers")
	if len(suppliers) != 1 {
		t.Logf("The number of sections is incorrect")
		t.Fail()
	}

	if value, ok := suppliers["macys"]; ok {
		if value != "15" {
			t.Logf("The value of option macys is incorrect")
			t.Fail()
		}
	} else  {
		t.Logf("The key name of option macys is incorrect")
		t.Fail()
	}

}
