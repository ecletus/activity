package activity

import (
	"github.com/aghape/db"
)

type Plugin struct {
	db.DisDBNames
}

func (p *Plugin) OnRegister() {
	p.DBOnMigrateGorm(func(e *db.GormDBEvent) error {
		return e.DB.AutoMigrate(&QorActivity{}).Error
	})
}
