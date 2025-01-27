package workers

type Closer struct {
	epCloser            func(int) error
	targetSocketHandler int
}

func (c *Closer) SetTarget(targetSocketHandler int) {
	c.targetSocketHandler = targetSocketHandler
}

func (c Closer) Close() error {
	if c.epCloser != nil {
		return c.epCloser(c.targetSocketHandler)
	}

	return nil
}
