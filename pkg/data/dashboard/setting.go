package dashboard

import (
	"github.com/pborman/uuid"
	"github.com/ranger/ranger/pkg/settings"
)

func addSetting() error {
	return settings.InstallUUID.SetIfUnset(uuid.NewRandom().String())
}
