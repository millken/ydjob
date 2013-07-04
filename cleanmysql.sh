#/usr/bin
#rotate tables

time=` date +%Y%m%d%H `
database="yundun_monitor"
echo $time
for i in firewall_hz
 do
  perl ./rotate_mrg.pl $database $i
  sleep 1
 done
