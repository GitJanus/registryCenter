package model

import (
	"github.com/jinzhu/gorm"
	"registryCenter/db"
)

type Cluster struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	HsmType    string `json:"hsm_type"`    //集群设备型号
	HsmVersion string `json:"hsm_version"` //版本号
	Status     int    `json:"status"`      //0停用，1启用
	PollMode   int    `json:"poll_mode"`   //1随机，2轮询，3加权随机，4加权轮询
	Domain     string `json:"domain"`
	MaxIdle    int    `json:"max_idle"`
	MinIdle    int    `json:"min_idle"`
	MaxActive  int    `json:"max_active"`
	Timeouts   int    `json:"timeouts"`
	CommModel  int    `json:"comm_model"` //通信方式   0：明文/1:单向ssl/2:双向ssl
}

var sqlDB *gorm.DB

func Init() {
	sqlDB = db.DB
}

func (Cluster) TableName() string {
	return "cluster"
}

func GetCluster() (Cluster, error) {
	cluster := Cluster{}
	err := sqlDB.First(&cluster).Error
	return cluster, err
}

func GetClusterCount() (int, error) {
	var count int
	err := sqlDB.Table("cluster").Count(&count).Error
	if err != nil {
		count = 0
	}
	return count, err
}

func InsertCluster(cluster *Cluster) interface{} {
	err := sqlDB.Create(cluster).Error
	return err
}

func DeleteAllCluster() error {
	err := sqlDB.Delete(&Cluster{}).Error
	return err
}

func UpdateCluster(clusterChannel chan Cluster) error {
	cluster := <-clusterChannel
	err := sqlDB.Model(&cluster).Where("id=?", 1).Update("poll_mode", cluster.PollMode).Error
	return err
}

func DeleteClusterByName(name string) error {
	err := sqlDB.Where("name = ?", name).Delete(&Cluster{}).Error
	return err
}
