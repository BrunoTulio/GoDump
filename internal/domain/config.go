package domain

import "time"

const (
	Weak Priority = iota
	Normal
	High
)

type (
	Priority int

	Alert struct {
		SendMail []string `yaml:"send_mail"`
		Enable   bool     `yaml:"enable"`
		Mail     Mail     `yaml:"mail"`
	}
	Mail struct {
		User           string `yaml:"user"`
		From           string `yaml:"from"`
		Password       string `yaml:"password"`
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		InsureSecurity bool   `yaml:"insureSecurity"`
	}

	Dump struct {
		Key      string   `yaml:"key"`
		Host     string   `yaml:"host"`
		Port     int      `yaml:"port"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
		Database string   `yaml:"database"`
		Priority Priority `yaml:"priority"`
		Type     *Type    `yaml:"type"`
	}

	GoogleDrive struct {
		Enable bool `yaml:"enable"`
	}

	Dumps []Dump

	Config struct {
		Cron                 string        `yaml:"cron"`
		DurationFileInFolder time.Duration `yaml:"duration_file_in_folder"`
		Alert                Alert         `yaml:"alert"`
		Dumps                Dumps         `yaml:"dumps"`
		GoogleDrive          GoogleDrive   `yaml:"google_drive"`
	}
)

func (c *Dump) UpdateType(fileType *Type) {
	c.Type = fileType
}

func (c *Dump) NoHasType() bool {
	return !c.HasType()
}

func (c *Dump) HasType() bool {
	return c.Type != nil
}

func (c *Dump) Extension() string {
	if c.NoHasType() {
		return DefaultTypeExtension
	}

	return c.Type.Extension()
}

func (c *Dump) IsTypeSQL() bool {
	return c.HasType() && *c.Type == SQL
}

func (c *Dump) IsTypeDump() bool {
	return c.HasType() && *c.Type == DUMP
}

func (c Dumps) UpdateType(typeName *Type) {
	for i := range c {
		c[i].UpdateType(typeName)
	}

}

func (c Dumps) Len() int {
	return len(c)
}
func (c Dumps) Less(i, j int) bool {
	return c[i].Priority > c[j].Priority
}
func (c Dumps) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c *Config) IsValid() error {
	return nil
}
