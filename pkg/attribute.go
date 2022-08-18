package gmx

import "github.com/google/uuid"

func (manager *attributeManager) Get(domain string, beanName string, attributeName string) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return internalExecuteAgainstAll(manager.mbeans, manager.maxNumberOfGoRoutines, func(u uuid.UUID) (string, error) {
		return manager.GetById(u, domain, beanName, attributeName)
	})
}

func (manager *attributeManager) Put(domain string, beanName string, attributeName string, value any) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return internalExecuteAgainstAll(manager.mbeans, manager.maxNumberOfGoRoutines, func(u uuid.UUID) (string, error) {
		return manager.PutById(u, domain, beanName, attributeName, value)
	})
}

func (manager *attributeManager) GetById(id uuid.UUID, domain string, beanName string, attributeName string) (string, error) {
	return "", nil
}

func (manager *attributeManager) PutById(id uuid.UUID, domain string, beanName string, attributeName string, value any) (string, error) {
	return "", nil
}
