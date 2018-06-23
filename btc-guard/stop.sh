pid=`pidof bitcoind`
echo $pid
kill -9 $pid
