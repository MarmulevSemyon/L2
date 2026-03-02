#!/bin/bash
set -u

BIN=./bin/cut       
FILE=test/test.txt

echo "Testing $BIN against system cut on $FILE"
echo

FIELDS=("1" "2" "1,3" "2-4" "10")
DELIMS=("," ":" "|" $'\t')

pass=0
fail=0

for d in "${DELIMS[@]}"; do
  for f in "${FIELDS[@]}"; do
    for s in 0 1; do
      if [ "$s" -eq 1 ]; then
        echo "CASE: -d $(printf '%q' "$d") -f $f -s"
        cut -d "$d" -f "$f" -s "$FILE" > expected.txt
        "$BIN" -d "$d" -f "$f" -s "$FILE" > actual.txt
      else
        echo "CASE: -d $(printf '%q' "$d") -f $f"
        cut -d "$d" -f "$f" "$FILE" > expected.txt
        "$BIN" -d "$d" -f "$f" "$FILE" > actual.txt
      fi

      if diff -u expected.txt actual.txt >/dev/null; then
        echo "  OK"
        pass=$((pass+1))
      else
        echo "  FAIL"
        diff -u expected.txt actual.txt | head -n 60
        fail=$((fail+1))
      fi
      echo
    done
  done
done

rm -f expected.txt actual.txt
echo "DONE: PASS=$pass FAIL=$fail"
exit $fail