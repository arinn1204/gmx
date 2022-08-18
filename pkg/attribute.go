package gmx

import (
	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/google/uuid"
)

func (manager *attributeManager) Get(domain string, beanName string, attributeName string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return internalExecuteAgainstAll(manager.mbeans, manager.maxNumberOfGoRoutines, func(u uuid.UUID) (string, error) {
		return manager.GetById(u, domain, beanName, attributeName, args...)
	})
}

func (manager *attributeManager) Put(domain string, beanName string, attributeName string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return internalExecuteAgainstAll(manager.mbeans, manager.maxNumberOfGoRoutines, func(u uuid.UUID) (string, error) {
		return manager.PutById(u, domain, beanName, attributeName, args...)
	})
}

func (manager *attributeManager) GetById(id uuid.UUID, domain string, beanName string, attributeName string, args ...MBeanArgs) (string, error) {
	operationArgs := make([]mbean.OperationArgs, 0)

	for _, arg := range args {
		operationArgs = append(
			operationArgs,
			mbean.OperationArgs{
				Value:             arg.Value,
				JavaType:          arg.JavaType,
				JavaContainerType: arg.JavaContainerType,
			},
		)
	}

	mbeanClient := (*manager.mbeans)[id]

	return mbeanClient.Get(domain, beanName, attributeName, operationArgs...)
}

func (manager *attributeManager) PutById(id uuid.UUID, domain string, beanName string, attributeName string, args ...MBeanArgs) (string, error) {
	operationArgs := make([]mbean.OperationArgs, 0)
	for _, arg := range args {
		operationArgs = append(
			operationArgs,
			mbean.OperationArgs{
				Value:             arg.Value,
				JavaType:          arg.JavaType,
				JavaContainerType: arg.JavaContainerType,
			},
		)
	}
	mbeanClient := (*manager.mbeans)[id]

	return mbeanClient.Put(domain, beanName, attributeName, operationArgs...)
}
