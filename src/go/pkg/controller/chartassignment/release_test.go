package chartassignment

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	apps "github.com/googlecloudrobotics/core/src/go/pkg/apis/apps/v1alpha1"
	"k8s.io/client-go/tools/record"
	"k8s.io/helm/pkg/chartutil"
)

func writeFile(t *testing.T, fn string, s string) {
	t.Helper()
	if err := ioutil.WriteFile(fn, []byte(strings.TrimSpace(s)), 0666); err != nil {
		t.Fatal(err)
	}
}

func buildInlineChart(t *testing.T, chart, values string) string {
	t.Helper()

	tmpdir, err := ioutil.TempDir("", "buildInlineChart")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	writeFile(t, path.Join(tmpdir, "Chart.yaml"), chart)
	writeFile(t, path.Join(tmpdir, "values.yaml"), values)

	ch, err := chartutil.LoadDir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	fn, err := chartutil.Save(ch, tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	rawChart, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(rawChart)
}

func verifyValues(t *testing.T, have string, wantValues chartutil.Values) {
	if want, err := wantValues.YAML(); err != nil {
		t.Fatal(err)
	} else if want != have {
		t.Fatalf("config values do not match: want\n%s\n\ngot\n%s\n", want, have)
	}
}

func Test_loadChart_mergesValues(t *testing.T) {
	chart := buildInlineChart(t, `
name: testchart
version: 2.1.0
	`, `
foo1:
  baz1: "hello"
bar1: 3
	`)

	var as apps.ChartAssignment
	unmarshalYAML(t, &as, `
metadata:
  name: test-assignment-1
spec:
  chart:
    values:
      bar1: 4
      bar2:
        baz2: test
	`)
	as.Spec.Chart.Inline = chart
	wantValues := chartutil.Values{
		"bar1": 4,
		"bar2": chartutil.Values{"baz2": "test"},
		"foo1": chartutil.Values{"baz1": "hello"},
	}

	_, vals, err := loadChart(&as.Spec.Chart)
	if err != nil {
		t.Fatal(err)
	}
	verifyValues(t, vals, wantValues)
}

func Test_updateSynk_callsApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chart := buildInlineChart(t, `
name: testchart
version: 2.1.0
	`, `
foo1:
  baz1: "hello"
bar1: 3
	`)

	var as apps.ChartAssignment
	unmarshalYAML(t, &as, `
metadata:
  name: test-assignment-1
spec:
  chart:
    values:
      bar1: 4
      bar2:
        baz2: test
	`)
	as.Spec.Chart.Inline = chart

	mockSynk := NewMockInterface(ctrl)
	r := &release{
		synk:     mockSynk,
		recorder: &record.FakeRecorder{},
	}

	rs := &apps.ResourceSet{}
	mockSynk.EXPECT().Apply(gomock.Any(), "test-assignment-1", gomock.Any(), gomock.Any()).Return(rs, nil).Times(1)

	// First apply, the chart should be installed.
	r.updateSynk(&as)
}

func Test_deleteSynk_callsDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chart := buildInlineChart(t, `
name: testchart
version: 2.1.0
	`, `
foo1:
  baz1: "hello"
bar1: 3
	`)

	var as apps.ChartAssignment
	unmarshalYAML(t, &as, `
metadata:
  name: test-assignment-1
spec:
  chart:
    values:
      bar1: 4
      bar2:
        baz2: test
	`)
	as.Spec.Chart.Inline = chart

	mockSynk := NewMockInterface(ctrl)
	r := &release{
		synk:     mockSynk,
		recorder: &record.FakeRecorder{},
	}

	mockSynk.EXPECT().Delete(gomock.Any(), "test-assignment-1").Return(nil).Times(1)

	// First apply, the chart should be installed.
	r.deleteSynk(&as)
}