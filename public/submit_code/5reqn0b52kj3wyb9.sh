#!bin/bash
go build  -o /app/5reqn0b52kj3wyb9 /app/5reqn0b52kj3wyb9.go && chmod +x /app/5reqn0b52kj3wyb9
echo '===0 START==='
echo -e "3
 1 2
100 200
-10 20"| /app/5reqn0b52kj3wyb9
echo '===0 END==='
