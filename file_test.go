package contentaddressable

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var supOid = "a2b71d6ee8997eb87b25ab42d566c44f6a32871752c7c73eb5578cb1182f7be0"

func TestFile(t *testing.T) {
	test := SetupFile(t)
	defer test.Teardown()

	filename := filepath.Join(test.Path, supOid)
	aw, err := NewFile(filename)
	assertEqual(t, nil, err)

	n, err := aw.Write([]byte("SUP"))
	assertEqual(t, nil, err)
	assertEqual(t, 3, n)

	created, err := aw.Accept()
	assertEqual(t, created, true)
	assertEqual(t, nil, err)

	by, err := ioutil.ReadFile(filename)
	assertEqual(t, nil, err)
	assertEqual(t, "SUP", string(by))

	assertEqual(t, nil, aw.Close())
}

func TestTempFileCleanup(t *testing.T) {
	test := SetupFile(t)
	defer test.Teardown()

	filename := filepath.Join(test.Path, supOid)
	tempFilename := filepath.Join(test.Path, supOid + DefaultSuffix)

	aw, err := NewFile(filename)
	assertEqual (t, nil, err)

	// Put the destination file in place, so Accept() won't rename the
	// temp file over the destination
	fp, err := os.Create(filename)
	assertEqual(t, nil, err)
	defer fp.Close()

	n, err := aw.Write([]byte("SUP"))
	assertEqual(t, nil, err)
	assertEqual(t, 3, n)

	created, err := aw.Accept()
	assertEqual(t, false, created)
	assertEqual(t, nil, err)

	assertEqual(t, nil, aw.Close())
	if _, err := os.Stat(tempFilename); err == nil {
		t.Fatalf("Temp file should not exist: %s", tempFilename)
	}
}

func TestFileMismatch(t *testing.T) {
	test := SetupFile(t)
	defer test.Teardown()

	filename := filepath.Join(test.Path, "b2b71d6ee8997eb87b25ab42d566c44f6a32871752c7c73eb5578cb1182f7be0")
	aw, err := NewFile(filename)
	assertEqual(t, nil, err)

	n, err := aw.Write([]byte("SUP"))
	assertEqual(t, nil, err)
	assertEqual(t, 3, n)

	created, err := aw.Accept()
	if err == nil || !strings.Contains(err.Error(), "Content mismatch") {
		t.Errorf("Expected mismatch error: %s", err)
	}
	assertEqual(t, created, false)
	if _, err := os.Stat(filename); err == nil {
		t.Fatalf("%s should not exist", filename)
	}

	assertEqual(t, nil, aw.Close())

	_, err = ioutil.ReadFile(filename)
	assertEqual(t, true, os.IsNotExist(err))
}

func TestFileCancel(t *testing.T) {
	test := SetupFile(t)
	defer test.Teardown()

	filename := filepath.Join(test.Path, supOid)
	aw, err := NewFile(filename)
	assertEqual(t, nil, err)

	n, err := aw.Write([]byte("SUP"))
	assertEqual(t, nil, err)
	assertEqual(t, 3, n)

	assertEqual(t, nil, aw.Close())

	for _, name := range []string{aw.filename, aw.tempFilename} {
		if _, err := os.Stat(name); err == nil {
			t.Errorf("%s exists?", name)
		}
	}
}

func TestFileDuel(t *testing.T) {
	test := SetupFile(t)
	defer test.Teardown()

	filename := filepath.Join(test.Path, supOid)
	aw, err := NewFile(filename)
	assertEqual(t, nil, err)
	defer aw.Close()

	if _, err := NewFile(filename); err == nil {
		t.Errorf("Expected a file open conflict!")
	}
}

func SetupFile(t *testing.T) *FileTest {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	path := filepath.Join(wd, "File")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err.Error())
	}

	return &FileTest{path, t}
}

type FileTest struct {
	Path string
	*testing.T
}

func (t *FileTest) Teardown() {
	if err := os.RemoveAll(t.Path); err != nil {
		t.Fatalf("Error removing %s: %s", t.Path, err)
	}
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	checkAssertion(t, expected, actual, "")
}

func assertEqualf(t *testing.T, expected, actual interface{}, format string, args ...interface{}) {
	checkAssertion(t, expected, actual, format, args...)
}

func checkAssertion(t *testing.T, expected, actual interface{}, format string, args ...interface{}) {
	if expected == nil {
		if actual == nil {
			return
		}
	} else if reflect.DeepEqual(expected, actual) {
		return
	}

	_, file, line, _ := runtime.Caller(2) // assertEqual + checkAssertion
	t.Logf("%s:%d\nExpected: %v\nActual:   %v", file, line, expected, actual)
	if len(args) > 0 {
		t.Logf("! - "+format, args...)
	}
	t.FailNow()
}
