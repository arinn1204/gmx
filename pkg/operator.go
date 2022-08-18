package gmx

import (
	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/google/uuid"
)

// ExecuteAgainstAll will execute a single command against every mbean that is currently registered.
// This will return a mapping of all results and errors, based on the UUID that the connection has been assigned.
//
// All executions will be run in separate go routines, so this needs to be planned for accordingly
func (operator *operator) ExecuteAgainstAll(domain string, name string, operation string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error) {
	return internalExecuteAgainstAll(operator.mbeans, operator.maxNumberOfGoRoutines, func(id uuid.UUID) (string, error) {
		return operator.ExecuteAgainstID(id, domain, name, operation, args...)
	})
}

// ExecuteAgainstID is a method that will take a given operation and MBean ID and make the JMX request.
// It will return whatever is returned downstream, errors and all
func (operator *operator) ExecuteAgainstID(id uuid.UUID, domain string, name string, operation string, args ...MBeanArgs) (string, error) {
	env := java.Attach()
	defer java.Detach()

	bean := (*operator.mbeans)[id].WithEnvironment(env)

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

	mbeanOp := mbean.Operation{
		Domain:    domain,
		Name:      name,
		Operation: operation,
		Args:      operationArgs,
	}

	return bean.Execute(mbeanOp)
}
