#!/bin/sh
# Exercise the real Makefile with isolated packages, including failing tests.
set -eu

repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/../.." && pwd)
fixture=$(mktemp -d "${TMPDIR:-/tmp}/gin-make-test.XXXXXX")
trap 'rm -rf "$fixture"' 0
trap 'exit 1' HUP INT TERM

cp "$repo_root/Makefile" "$fixture/Makefile"
cd "$fixture"
export GOWORK=off GOFLAGS='' GO111MODULE=on
unset TESTTAGS MAKEFLAGS MFLAGS

cat > go.mod <<'EOF'
module example.com/gin

go 1.26.0
EOF

for package in gin binding ginS render; do
	directory=$package
	if [ "$package" = gin ]; then
		directory=.
	fi
	mkdir -p "$directory"
	cat > "$directory/fixture.go" <<'EOF'
package fixture

func Value() int { return 1 }

var flagValue = "unset"
EOF
	printf 'package fixture\n\nconst packageName = "%s"\n' "$package" > "$directory/name_test.go"
	cat > "$directory/fixture_test.go" <<'EOF'
package fixture

import (
	"os"
	"path/filepath"
	"testing"
)

func init() {
	if packageName == "gin" && os.Getenv("GIN_MAKE_TEST_FAILURE") == "panic" {
		panic("deliberate initialization panic")
	}
}

func TestValue(t *testing.T) {
	marker := filepath.Join(os.Getenv("GIN_MAKE_TEST_MARKERS"), packageName)
	if err := os.WriteFile(marker, []byte("executed"), 0600); err != nil {
		t.Fatal(err)
	}
	if Value() != 1 {
		t.Fatal("unexpected value")
	}
	if packageName == "gin" && os.Getenv("GIN_MAKE_TEST_FAILURE") == "assertion" {
		t.Fatal("deliberate assertion failure")
	}
}
EOF
done

cat > tagged_test.go <<'EOF'
//go:build makefiletest

package fixture

import "testing"

func TestTagged(t *testing.T) {
	if flagValue != "two words" {
		t.Fatalf("quoted linker flag was not forwarded: %q", flagValue)
	}
}
EOF

cat > race_test.go <<'EOF'
//go:build makefilerace

package fixture

import "testing"

func TestRace(t *testing.T) {
	var value int
	done := make(chan struct{})
	go func() {
		value = 1
		close(done)
	}()
	value = 2
	<-done
	t.Log(value)
}
EOF

fail() {
	printf 'FAIL: %s: %s\n' "$case_name" "$1" >&2
	cat "$case_name.log" >&2
	exit 1
}

run_make() {
	case_name=$1
	expected=$2
	shift 2
	# A unique directory also prevents Go's test cache from hiding execution.
	export GIN_MAKE_TEST_MARKERS="$fixture/$case_name.markers"
	mkdir "$GIN_MAKE_TEST_MARKERS"
	rm -f coverage.out
	status=0
	make --no-print-directory test "$@" > "$case_name.log" 2>&1 || status=$?
	if [ "$expected" = pass ] && [ "$status" -ne 0 ]; then
		fail "make test exited with $status"
	elif [ "$expected" = fail ] && [ "$status" -eq 0 ]; then
		fail 'make test accepted a failing command'
	fi
	printf 'PASS: %s (exit %s)\n' "$case_name" "$status"
}

check_coverage() {
	grep -qx "mode: $1" coverage.out || fail 'incorrect coverage mode'
	for package in gin binding ginS render; do
		[ -f "$GIN_MAKE_TEST_MARKERS/$package" ] || fail "$package tests did not execute"
		prefix=example.com/gin
		if [ "$package" != gin ]; then
			prefix=$prefix/$package
		fi
		grep -q "^$prefix/fixture.go:" coverage.out || fail "$package coverage is missing"
	done
	go tool cover -func=coverage.out > "$case_name.coverage" || fail 'invalid coverage profile'
}

export GIN_MAKE_TEST_FAILURE=
run_make default pass
check_coverage set

run_make empty-tags pass TESTTAGS=
check_coverage set

run_make quoted-flags pass "TESTTAGS=-tags makefiletest -ldflags=\"-X 'example.com/gin.flagValue=two words'\""
check_coverage set
grep -q '^--- PASS: TestTagged' "$case_name.log" || fail 'tagged test did not execute'

run_make race pass TESTTAGS=-race
check_coverage atomic

export GOFLAGS=-race
run_make race-go-flags pass
check_coverage atomic
export GOFLAGS=

export GIN_MAKE_TEST_FAILURE=assertion
run_make assertion-failure fail
export GIN_MAKE_TEST_FAILURE=panic
run_make initialization-panic fail
export GIN_MAKE_TEST_FAILURE=

printf 'package fixture\nvar _ = undefinedIdentifier\n' > broken.go
run_make compilation-failure fail
rm broken.go

# Package discovery must not silently fall back to testing only the root.
printf '//go:build (\n\npackage fixture\n' > ginS/broken.go
run_make package-loading-failure fail
rm ginS/broken.go

run_make invalid-flag fail TESTTAGS=-invalid-test-flag

cat > failing-go <<'EOF'
#!/bin/sh
echo 'deliberate tool failure' >&2
exit 17
EOF
chmod +x failing-go
run_make tool-failure fail GO=./failing-go
grep -q 'deliberate tool failure' "$case_name.log" || fail 'tool diagnostic was lost'

run_make detected-race fail 'TESTTAGS=-race -tags makefilerace'
grep -q 'WARNING: DATA RACE' "$case_name.log" || fail 'race detector did not report the race'
