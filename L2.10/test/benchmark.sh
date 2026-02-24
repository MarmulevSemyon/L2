#!/usr/bin/env bash

# set -e

FILE="test/big_test.txt"
MY_OUT="my_sort.txt"
SYS_OUT="sys_sort.txt"

echo "===== размер файла ====="
du -h $FILE

# echo
# echo "Идёт сортировка с помошью sort big_test.txt -n..."
# START=$(date +%s.%N)
# sort -n $FILE > $SYS_OUT
# END=$(date +%s.%N)
# SYS_TIME=$(echo "$END - $START" | bc)
# echo "сортировка GNU sort -n заняла: $SYS_TIME секунд"

echo
echo "Идёт сортировка  с помощью ./bin/sort -n big_test.txt ..."
START=$(date +%s.%N)
./bin/sort -n $FILE > $MY_OUT
END=$(date +%s.%N)
MY_TIME=$(echo "$END - $START" | bc)
echo "моя сортировка заняла: $MY_TIME секунд"

# echo
# echo "===== сравнение отсортированных файлов ====="
# if diff --strip-trailing-cr $SYS_OUT $MY_OUT > /dev/null; then
#     echo "OK. похоже на правду"
# else
#     echo "!!!ERROR!!! отсортировалось по разному"
#     diff --strip-trailing-cr -u sys_sort.txt my_sort.txt | sed -n '1,80p'
#     cmp -l sys_sort.txt my_sort.txt | head
# fi

# echo
# echo "===== итог ====="
# echo "sort GNU: $SYS_TIME секунд"
# echo "./bin/sort:     $MY_TIME секунд"