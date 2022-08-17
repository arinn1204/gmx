package gmx

import "github.com/google/uuid"

func (mbean *Client) Get(domain string, beanName string, attributeName string) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return mbean.internalExecuteAgainstAll(func(u uuid.UUID) (string, error) {
		return mbean.GetById(u, domain, beanName, attributeName)
	})
}

func (mbean *Client) Put(domain string, beanName string, attributeName string, value any) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return mbean.internalExecuteAgainstAll(func(u uuid.UUID) (string, error) {
		return mbean.PutById(u, domain, beanName, attributeName, value)
	})
}

func (mbean *Client) GetById(id uuid.UUID, domain string, beanName string, attributeName string) (string, error) {
	return "", nil
}

func (mbean *Client) PutById(id uuid.UUID, domain string, beanName string, attributeName string, value any) (string, error) {
	return "", nil
}
