package model

type Hsm struct {
	HsmIp      string `json:"hsm_ip"`
	HsmPort    int    `json:"hsm_port"` //密码机端口
	HsmCode    string `json:"hsm_code"` //密码机编码
	Weight     int    `json:"weight"`
	Status     int    `json:"status"` //1：启用，2：停用，3：上线中，4：故障
	Model      string `json:"model"`
	Sn         string `json:"sn"`
	HsmVersion string `json:"hsm_version"`
	ExpiryDate string `json:"expiry_date"`
	CreateTime string `json:"create_time"`
	MaxIdle    int    `json:"max_idle"`
	MinIdle    int    `json:"min_idle"`
	MaxActive  int    `json:"max_active"`
	Timeouts   int    `json:"timeouts"`
	HsmType    string `json:"hsm_type"` //密码机类型(SJJ1507.....)
}

func (Hsm) TableName() string {
	return "hsm"
}

func UpdateHsm(hsm Hsm) error {
	err := sqlDB.Model(&hsm).Where("hsm_code = ?", hsm.HsmCode).Update("status", hsm.Status).Error
	return err
}

func GetHsmList() ([]Hsm, error) {
	hsmList := []Hsm{}
	err := sqlDB.Find(&hsmList).Error
	return hsmList, err
}

func GetHsmListByNotCode(code interface{}) ([]Hsm, error) {
	hsmList := []Hsm{}
	err := sqlDB.Where("hsm_code != ?", code).Find(&hsmList).Error
	return hsmList, err
}

func GetHsmListByStatus(status int) ([]Hsm, error) {
	hsmList := []Hsm{}
	err := sqlDB.Where("status = ?", status).Find(&hsmList).Error
	return hsmList, err
}

func GetHsmCountByCode(code string) (int, error) {
	var count int
	err := sqlDB.Model(&Hsm{}).Where(&Hsm{HsmCode: code}).Count(&count).Error
	if err != nil {
		count = 0
	}
	return count, err
}

func GetHsmByCode(code string) (Hsm, error) {
	hsm := Hsm{}
	err := sqlDB.Where("hsm_code = ?", code).Find(&hsm).Error
	return hsm, err
}

func InsertHsm(hsm *Hsm) interface{} {
	err := sqlDB.Create(hsm).Error
	return err
}

func DeleteAllHsm() interface{} {
	err := sqlDB.Delete(&Hsm{}).Error
	return err
}

func DeleteHsmByCode(code string) error {
	err := sqlDB.Where("hsm_code = ?", code).Delete(&Hsm{}).Error
	return err
}
