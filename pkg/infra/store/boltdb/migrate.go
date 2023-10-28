package boltdb

type bucket int

const (
	SettingsBucket bucket = iota

	lastBucket
)

func (s bucket) Int() int {
	return int(s)
}

func (s bucket) String() string {
	return [...]string{"instance_settings"}[s]
}

func (d *Database) Migrate() error {
	for i := 0; i < lastBucket.Int(); i++ {
		if err := d.SetServiceName(bucketName(i)); err != nil {
			return err
		}
	}

	return nil
}

func bucketName(b int) string {
	return bucket(b).String()
}
