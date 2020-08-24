package daemon

import (
	"testing"
)

func TestReadVersionFromFile(t *testing.T) {
	ver, rel, err := ReadVersionFromFile("../version.txt")
	if err != nil {
		t.Fatalf("Failed encounter %s", err)
	}
	expectedVersion := Version{
		Minor: 0,
		Major: 1,
		Patch: 0,
	}
	expectedRelease := "2020.08.14"
	if ver != expectedVersion {
		t.Errorf("expected %s, got %s", expectedVersion.ToString(), ver.ToString() )
	}
	if expectedRelease != rel {
		t.Errorf("expected %s, got %s", expectedRelease, rel )
	}
}
