package logx

type LogRecord struct {
	ID      int64  `gorm:"primaryKey;autoIncrement"`
	Level   string `gorm:"type:varchar(16);index;"`
	Message string `gorm:"type:text;"`
	File    string `gorm:"type:varchar(255)"`
	Created int    `gorm:"index;"`
}
