import javax.management.MBeanParameterInfo;
import javax.management.MBeanServerConnection;
import javax.management.MalformedObjectNameException;
import javax.management.ObjectName;
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
            final JMXServiceURL url = new JMXServiceURL(jndiUri);
            try(final JMXConnector jmxc = JMXConnectorFactory.connect(url, null)) {
                final MBeanServerConnection connection = jmxc.getMBeanServerConnection();

                var beanInfo = connection.getMBeanInfo(objectName);
                var op = Arrays.stream(beanInfo.getOperations())
                        .filter(f -> f.getName().equals(operation))
                        .findFirst()
                        .orElseThrow();

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
