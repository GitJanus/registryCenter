package model

type Instructions struct {
	Id   string
	Code int
}

func (Instructions) TableName() string {
	return "instructions"
}

func GetInstructionsListByCode(code int) ([]Instructions, error) {
	instructionsList := []Instructions{}
	err := sqlDB.Where("code = ?", code).Find(&instructionsList).Error
	return instructionsList, err
}
