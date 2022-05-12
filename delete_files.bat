@REM 刪除資料
echo "delete files"
pause

rd /Q /S ".\mongodb\data"
rd /Q /S ".\postgres\data\"
rd /Q /S ".\rabbitmq\data\"
rd /Q /S ".\redis\data\"

mkdir ".\mongodb\data"
mkdir ".\postgres\data"
mkdir ".\rabbitmq\data"
mkdir ".\redis\data"