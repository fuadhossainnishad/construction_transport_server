package postgres

func (c *Client) Close() {
	if c.Pool != nil {
		c.Pool.Close()
		c.logger.Info("[DB] pool closed gracefully")
	}
}
