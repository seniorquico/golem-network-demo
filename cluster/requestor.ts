import { CoreV1Api, V1Secret, V1Pod } from '@kubernetes/client-node';

export interface Node {
  name: string;
  environment: Record<string, string>;
  offerTemplate: Record<string, any>;
  presets: Record<string, any>;
}

function generateNames(name: string) {
  return {
    environmentName: `${name}-env`,
    podName: name,
  };
}

/**
 * Deprovisions a Salad Organization on the Golem Network by deleting its Pod and Secret.
 */
export async function deprovisionOrganization(
  k8sApi: CoreV1Api,
  namespace: string,
  node: Node
) {
  const names = generateNames(node.name);

  // Delete Pod
  try {
    await k8sApi.deleteNamespacedPod(names.podName, namespace);
  } catch (err) {
    // Ignore if not found
  }

  // Delete Secret for environment variables
  try {
    await k8sApi.deleteNamespacedSecret(names.environmentName, namespace);
  } catch (err) {
    // Ignore if not found
  }
}

/**
 * Provisions a Salad Organization on the Golem Network by creating its Pod and Secret.
 */
export async function provisionOrganization(
  k8sApi: CoreV1Api,
  namespace: string,
  node: Node
) {
  const names = generateNames(node.name);

  // Create or update Secret for environment variables
  const secret: V1Secret = {
    metadata: { name: names.environmentName, namespace },
    type: 'Opaque',
    stringData: node.environment,
  };
  try {
    await k8sApi.readNamespacedSecret(names.environmentName, namespace);
    await k8sApi.replaceNamespacedSecret(names.environmentName, namespace, secret);
  } catch {
    await k8sApi.createNamespacedSecret(namespace, secret);
  }

  // Create or update Pod
  const pod: V1Pod = {
    metadata: { name: names.podName, namespace },
    spec: {
      containers: [
        {
          name: 'golem-requestor',
          image: 'saladtechnologies/golem-requestor:v0.17.6-1',
          imagePullPolicy: 'Never',
          args: ['service', 'run'],
          resources: {
            limits: { cpu: '2', memory: '1Gi' },
            requests: { cpu: '10m', memory: '64Mi' },
          },
          envFrom: [
            { secretRef: { name: names.environmentName } },
          ],
        },
      ],
      restartPolicy: 'Never',
    },
  };
  try {
    await k8sApi.readNamespacedPod(names.podName, namespace);
    await k8sApi.replaceNamespacedPod(names.podName, namespace, pod);
  } catch {
    await k8sApi.createNamespacedPod(namespace, pod);
  }
}
