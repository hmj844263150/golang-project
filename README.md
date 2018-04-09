        CREATE TABLE testdata (
          id integer NOT NULL primary key,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly boolean NOT NULL,
          module_id integer NOT NULL,
          device_type TEXT NOT NULL,
          fw_ver TEXT NOT NULL,
          esp_mac TEXT NOT NULL,
          cus_mac TEXT NOT NULL,
          flash_id TEXT NOT NULL,
          test_result TEXT NOT NULL,
          test_msg TEXT NOT NULL,
          factory_sid TEXT NOT NULL,
          batch_sid TEXT NOT NULL,
          efuse TEXT NOT NULL,
          query_times integer NOT NULL,
          print_times integer NOT NULL,
          batch_index integer NOT NULL,
          latest boolean NOT NULL
        );
        CREATE TABLE testlog (
          id integer NOT NULL primary key,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly boolean NOT NULL,
          module_id integer NOT NULL,
          device_type TEXT NOT NULL,
          fw_ver TEXT NOT NULL,
          esp_mac TEXT NOT NULL,
          cus_mac TEXT NOT NULL,
          flash_id TEXT NOT NULL,
          test_result TEXT NOT NULL,
          test_msg TEXT NOT NULL,
          factory_sid TEXT NOT NULL,
          batch_sid TEXT NOT NULL,
          efuse TEXT NOT NULL,
          query_times integer NOT NULL,
          print_times integer NOT NULL
          batch_index integer NOT NULL,
          latest boolean NOT NULL
        );
        CREATE TABLE batch (
          id integer NOT NULL primary key,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          sid TEXT NOT NULL,
          name TEXT NOT NULL,
          desc TEXT NOT NULL,
          factory_sid TEXT NOT NULL,
          cnt integer NOT NULL,
          remain integer NOT NULL,
          esp_mac_from TEXT NOT NULL,
          esp_mac_to TEXT NOT NULL,
          cus_mac_from TEXT NOT NULL,
          cus_mac_to TEXT NOT NULL,
          esp_mac_num_from integer NOT NULL,
          esp_mac_num_to integer NOT NULL,
          cus_mac_num_from integer NOT NULL,
          cus_mac_num_to integer NOT NULL,
          is_cus bool NOT NULL
        );
        CREATE TABLE factory (
          id integer NOT NULL primary key,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          sid TEXT NOT NULL,
          name TEXT NOT NULL,
          location TEXT NOT NULL,
          token TEXT NOT NULL,
          is_staff bool NOT NULL
        );
        CREATE TABLE module (
          id integer NOT NULL primary key,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          esp_mac TEXT NOT NULL
        );
1 create db
echo '{"path":"/hengha/db/create"}' | nc localhost 7000 | python -mjson.tool
2 drop db
echo '{"path":"/hengha/db/drop"}' | nc localhost 7000 | python -mjson.tool
3 post testdata
echo '{"path":"/testdata","method":"POST","testdata":{"device_type":"ESP_WROOM02","fw_ver":"v1.0.0.0","esp_mac":"xx:xx:xx:xx:xx:xx","flash_id":"19191919","test_result":"success","factory_sid":"acef234562345","batch_sid":"__201609222","efuse":"xxxxxxxxxxxxxxxxxxxx"}}' | nc localhost 7000 | python -mjson.tool
4 query logs
echo '{"path":"/testdata/logs","get": {"esp_mac": "xx:xx:xx:xx:xx:xx"}, "method":"Get"}' | nc localhost 7000 | python -mjson.tool
5 print testdata
echo '{"path":"/testdata/print","get": {"esp_mac": "xx:xx:xx:xx:xx:xx", "dryrun": true}, "method":"Post"}' | nc localhost 7000 | python -mjson.tool
6 post batch
echo '{"path":"/batch", "method":"Post", "batch": {"sid": "__201609222", "factory_sid": "acef234562345", "esp_mac_from": "18:fe:34:77:b9:00", "esp_mac_to": "18:fe:34:77:b9:ef", "cus_mac_from": "b0:83:fe:b3:9a:00", "cus_mac_to": "b0:83:fe:b3:9a:ef"}}' | nc localhost 7000 | python -mjson.tool
7 list batch
echo '{"path":"/batchs","method":"Get"}' | nc localhost 7000 | python -mjson.tool
8 stats by mac
echo '{"path":"/testdata/stats","get": {"by": "mac", "esp_mac": "xx:xx:xx:xx:xx:xx"}, "method":"Get"}' | nc localhost 7000 | python -mjson.tool
9 stats by batch
echo '{"path":"/testdata/stats","get": {"by": "batch", "batch_sid": "__201609222"}, "method":"Get"}' | nc localhost 7000 | python -mjson.tool
10 stats by datatime
echo '{"path":"/testdata/stats","get": {"by": "datetime"}, "method":"Get"}' | nc localhost 7000 | python -mjson.tool
11 stats by datatime batch
echo '{"path":"/testdata/stats","get": {"by": "datetime_batch", "batch_sid": "__201609222"}, "method":"Get"}' | nc localhost 7000 | python -mjson.tool

protoc --run_out=. *.proto ; sed -i 's:espressif.com/cloud/poster/db:espressif.com/chip/factory/db:g' *.go ; sed -i 's:var _ Doer:var _ db.Doer:g' *.go
sudo apt-get install gcc-mingw-w64
env CGO_ENABLED=1 GOOS=windows GOARCH=386 CC="i686-w64-mingw32-gcc" go build -o ha.exe apps/ha/ha.go
