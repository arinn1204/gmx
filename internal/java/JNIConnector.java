package internal.java;

import javax.management.*;
import javax.management.remote.JMXConnector;
import javax.management.remote.JMXConnectorFactory;
import javax.management.remote.JMXServiceURL;
import java.util.Arrays;
import java.util.Set;

public class JNIConnector {
    public static String getValue(String address, String mbean, String operation, String[] args) {
        String result;

        try {
            String jndiUri = String.format("service:jmx:rmi:///jndi/rmi://%s/jmxrmi", address);
            ObjectName objectName = new ObjectName(mbean);
            JMXServiceURL url = new JMXServiceURL(jndiUri);
            try(final JMXConnector jmxc = JMXConnectorFactory.connect(url, null)) {
                MBeanServerConnection connection = jmxc.getMBeanServerConnection();

                MBeanInfo beanInfo = connection.getMBeanInfo(objectName);
                MBeanOperationInfo op = Arrays.stream(beanInfo.getOperations())
                        .filter(f -> f.getName().equals(operation))
                        .findFirst()
                        .orElseThrow(RuntimeException::new);

                String[] paramTypes = Arrays.stream(op.getSignature())
                        .map(MBeanParameterInfo::getType)
                        .toArray(String[]::new);

                result = (String) connection.invoke(
                        objectName,
                        operation,
                        args,
                        paramTypes);
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        return result;
    }
}
