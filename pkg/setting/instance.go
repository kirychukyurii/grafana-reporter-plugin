package setting

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store"
)

const storeBucket = "instance_settings"

type InstanceSetting struct {
	OrgID    int
	Timezone string `json:"timezone"`
}

func NewInstanceSetting(settings backend.AppInstanceSettings) *InstanceSetting {
	return nil
}

func (s *InstanceSetting) Store(db store.DatabaseManager) error {
	if err := db.CreateObjectWithId(storeBucket, s.OrgID, s); err != nil {
		return err
	}

	return nil
}

func InstanceSettingFromStore(db store.DatabaseManager) ([]InstanceSetting, error) {
	var (
		setting  InstanceSetting
		settings = make([]InstanceSetting, 0)
	)

	if err := db.GetAll(storeBucket, &setting, boltdb.AppendFn(&settings)); err != nil {
		return nil, err
	}

	return settings, nil
}
