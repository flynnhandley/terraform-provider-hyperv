package hyperv

// Config ...
type Config struct {
	Hypervisor string
	Username   string
	Password   string
	Clustered  bool
	Driver     Driver
}

// GetDriver returns a new client for accessing PowerDNS
func (c *Config) GetDriver() (Driver, error) {

	driver, err := NewPS5Driver(c.Username, c.Password, c.Hypervisor)
	if err != nil {
		return nil, err
	}

	err = driver.TestConnectivity()

	if err != nil {
		return nil, err
	}
	return driver, nil
}
