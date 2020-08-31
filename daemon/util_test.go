package daemon

import "testing"

func TestGetLocalIpAddresses(t *testing.T) {
	t.Logf("start testing.")
	res := GetLocalIpAddresses()
	if len(res) == 0 {
		t.FailNow()
	}
	for i := range res {
		t.Logf("address %s", res[i])
		t.Logf("address %s", res[i])
	}
}

func TestGetPublicIpAddresses(t *testing.T) {
	t.Logf("start testing.")
	res := GetPublicIpAddresses()
	if len(res) == 0 {
		t.FailNow()
	}
	t.Logf("address %s", res)
}
