#!/bin/bash

set -e

FILE="test/big_test.txt"
N=1000000000 #~7/9GB
touch $FILE
echo "=====генерация $N строк на питоне (числа) ====="
START=$(date +%s.%N)
python3 - <<PY
import random
random.seed(42)
with open("$FILE","w") as f:
    for i in range($N):
        if i!= $N-1:
            f.write(str(random.randint(-10**7,10**7)) + "\n")
        else:
            f.write(str(random.randint(-10**7,10**7)))
PY
END=$(date +%s.%N)
GEN_TIME=$(echo "$END - $START" | bc)
echo "===== размер файла ====="
du -h $FILE
echo "Генерация файла заняла: $GEN_TIME секунд"