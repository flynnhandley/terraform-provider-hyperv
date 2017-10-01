package hyperv

// Config ...
type Config struct {
	Hypervisor string
	Username   string
	Password   string
	UseSSL     bool
	Driver     Driver
}

// GetDriver returns a new client for accessing Hypev
func (c *Config) GetDriver() (Driver, error) {

	driver, err := NewPS5Driver(c.Username, c.Password, c.Hypervisor, c.UseSSL)
	if err != nil {
		return nil, err
	}

	err = driver.TestConnectivity()

	if err != nil {
		return nil, err
	}
	return driver, nil
}
