#!/usr/bin/env bash

DATA="test/testdata.txt"
SYS_OUT="sys_test_temp.txt"
MY_OUT="my_test_temp.txt"

run_test() {
    local name="$1"
    shift

    echo "=== $name ==="

    grep "$@" "$DATA" > "$SYS_OUT"

    bin/grep "$@" "$DATA" > "$MY_OUT"

    if diff --strip-trailing-cr -u "$SYS_OUT" "$MY_OUT" > /dev/null; then
        echo "PASS"
    else
        echo "FAIL"
        diff "$SYS_OUT" "$MY_OUT"
        exit 1
    fi

    echo
}

run_test "basic" "abc"
run_test "-F literal dot" -F "a.c"
run_test "-i ignore case" -i "hello"
run_test "-v invert" -v "abc"
run_test "-n line numbers" -n "abc"
run_test "-c count" -c "repeat"
run_test "-A 2 after" -A 2 "курс"
run_test "-B 2 before" -B 2 "курс"
run_test "-C 2 context" -C 2 "курс"

echo "All tests passed!"

rm -f "$SYS_OUT" "$MY_OUT"