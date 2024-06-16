#!bin/bash
go build  -o /app/dtogj8t72lnt585p /app/dtogj8t72lnt585p.go && chmod +x /app/dtogj8t72lnt585p
echo '===0 START==='
echo -e "3
 1 2
100 200
-10 20"| /app/dtogj8t72lnt585p
echo '===0 END==='
