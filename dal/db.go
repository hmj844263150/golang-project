package dal

import (
	"espressif.com/chip/factory/db"
	"fmt"
	"log"
)

func Createdb() {
	db.Open("factory").Exec(`
        CREATE TABLE testdata (
          id int(11) NOT NULL AUTO_INCREMENT,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          module_id int(11) NOT NULL,
          device_type varchar(64) NOT NULL,
          fw_ver varchar(64) NOT NULL,
          esp_mac varchar(64) NOT NULL,
          cus_mac varchar(64) NOT NULL,
          flash_id varchar(64) NOT NULL,
          test_result varchar(64) NOT NULL,
          test_msg varchar(64) NOT NULL,
          factory_sid varchar(64) NOT NULL,
          batch_sid varchar(64) NOT NULL,
          efuse varchar(64) NOT NULL,
          query_times int(11) NOT NULL,
          print_times int(11) NOT NULL,
          batch_index int(11) NOT NULL,
          latest bool NOT NULL,
          PRIMARY KEY (id)
        );
        CREATE TABLE testlog (
          id int(11) NOT NULL AUTO_INCREMENT,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          module_id int(11) NOT NULL,
          device_type varchar(64) NOT NULL,
          fw_ver varchar(64) NOT NULL,
          esp_mac varchar(64) NOT NULL,
          cus_mac varchar(64) NOT NULL,
          flash_id varchar(64) NOT NULL,
          test_result varchar(64) NOT NULL,
          test_msg varchar(64) NOT NULL,
          factory_sid varchar(64) NOT NULL,
          batch_sid varchar(64) NOT NULL,
          efuse varchar(64) NOT NULL,
          query_times int(11) NOT NULL,
          print_times int(11) NOT NULL,
          batch_index int(11) NOT NULL,
          latest bool NOT NULL,
          PRIMARY KEY (id)
        );
        CREATE TABLE batch (
          id int(11) NOT NULL AUTO_INCREMENT,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          sid varchar(64) NOT NULL,
          factory_sid varchar(64) NOT NULL,
          name varchar(64) NOT NULL,
          desc varchar(128) NOT NULL,
          cnt int(11) NOT NULL,
          remain int(11) NOT NULL,
          esp_mac_from varchar(64) NOT NULL,
          esp_mac_to varchar(64) NOT NULL,
          cus_mac_from varchar(64) NOT NULL,
          cus_mac_to varchar(64) NOT NULL,
          esp_mac_num_from int(11) NOT NULL,
          esp_mac_num_to int(11) NOT NULL,
          cus_mac_num_from int(11) NOT NULL,
          cus_mac_num_to int(11) NOT NULL,
          is_cus bool NOT NULL,
          PRIMARY KEY (id)
        );
        CREATE TABLE factory (
          id int(11) NOT NULL AUTO_INCREMENT,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          sid varchar(64) NOT NULL,
          name varchar(64) NOT NULL,
          location varchar(64) NOT NULL,
          token varchar(64) NOT NULL,
          is_staff bool NOT NULL,
          PRIMARY KEY (id)
        );
        CREATE TABLE module (
          id int(11) NOT NULL AUTO_INCREMENT,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly bool NOT NULL,
          esp_mac varchar(64) NOT NULL,
          PRIMARY KEY (id)
        );
        CREATE TABLE user (
          id int(11) NOT NULL AUTO_INCREMENT,
          created datetime NOT NULL,
          updated datetime NOT NULL,
          visibly tinyint(1) NOT NULL,
          account varchar(64) NOT NULL,
          password varchar(64) NOT NULL,
          name varchar(64) NOT NULL,
          factory_sid varchar(64) NOT NULL,
          group_id tinyint(4) NOT NULL,
          email varchar(128) NOT NULL,
          description varchar(256) NOT NULL,
          PRIMARY KEY (id),
          UNIQUE KEY account_UNIQUE (account)
        );
    `)
}

func Dropdb() {
	db.Open("factory").Exec(`
        Drop table testdata;
        Drop table testlog;
        Drop table batch;
        Drop table factory;
        Drop table module;
    `)
}

func init() {
	factory := FindFactoryStaff(nil)
	if factory == nil {
		factory = NewFactory(nil)
		factory.Name = "Espressif"
		factory.Location = "Shanghai"
		factory.Generate()
		factory.IsStaff = true
		factory.Save()
	}
	log.Println(fmt.Sprintf("staff factory name: %s, sid: %s, token: %s", factory.Name, factory.Sid, factory.Token))
}
