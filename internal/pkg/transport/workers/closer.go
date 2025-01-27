package workers

type Closer struct {
	epCloser            func(int) error
	targetSocketHandler int
}

func (c *Closer) SetTarget(targetSocketHandler int) {
	c.targetSocketHandler = targetSocketHandler
}

func (c Closer) Close() error {
	return c.epCloser(c.targetSocketHandler)
}
