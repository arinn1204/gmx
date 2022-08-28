package gmx

import (
	"fmt"

	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/google/uuid"
)

func (manager *attributeManager) Get(domain string, beanName string, attributeName string, args MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return internalExecuteAgainstAll(manager.numberOfConnections, manager.mbeans, manager.maxNumberOfGoRoutines, func(u uuid.UUID) (string, error) {
		return manager.GetByID(u, domain, beanName, attributeName, args)
	})
}

func (manager *attributeManager) Put(domain string, beanName string, attributeName string, args MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return internalExecuteAgainstAll(manager.numberOfConnections, manager.mbeans, manager.maxNumberOfGoRoutines, func(u uuid.UUID) (string, error) {
		return manager.PutByID(u, domain, beanName, attributeName, args)
	})
}

func (manager *attributeManager) GetByID(id uuid.UUID, domain string, beanName string, attributeName string, args MBeanArgs) (string, error) {
	operationArgs := mbean.OperationArgs{
		Value:             args.Value,
		JavaType:          args.JavaType,
		JavaContainerType: args.JavaContainerType,
	}

	env := java.Attach()
	defer java.Detach()

	if mbeanClient, ok := (*manager.mbeans).Load(id); ok {
		return mbeanClient.(mbean.BeanExecutor).WithEnvironment(env).Get(domain, beanName, attributeName, operationArgs)
	}

	return "", fmt.Errorf("id of %s does not exist as an established connection to execute get on", id.String())
}

func (manager *attributeManager) PutByID(id uuid.UUID, domain string, beanName string, attributeName string, args MBeanArgs) (string, error) {
	operationArgs := mbean.OperationArgs{
		Value:             args.Value,
		JavaType:          args.JavaType,
		JavaContainerType: args.JavaContainerType,
	}

	env := java.Attach()
	defer java.Detach()

	if mbeanClient, ok := (*manager.mbeans).Load(id); ok {
		return mbeanClient.(mbean.BeanExecutor).WithEnvironment(env).Put(domain, beanName, attributeName, operationArgs)
	}

	return "", fmt.Errorf("id of %s does not exist as an established connection to execute put on", id.String())
}
