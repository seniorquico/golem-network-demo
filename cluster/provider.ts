import { CoreV1Api, V1Secret, V1ConfigMap, V1Pod } from '@kubernetes/client-node';

export interface Node {
  name: string;
  environment: Record<string, string>;
  offerTemplate: Record<string, any>;
  presets: Record<string, any>;
}

function generateNames(name: string) {
  return {
    environmentName: `${name}-env`,
    offerTemplateName: `${name}-template`,
    podName: name,
    presetsName: `${name}-presets`,
  };
}

/**
 * Deprovisions a Salad Node on the Golem Network by deleting its Pod, Secret, and ConfigMaps.
 */
export async function deprovisionNode(
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

  // Delete ConfigMap for presets
  try {
    await k8sApi.deleteNamespacedConfigMap(names.presetsName, namespace);
  } catch (err) {
    // Ignore if not found
  }

  // Delete ConfigMap for offer template
  try {
    await k8sApi.deleteNamespacedConfigMap(names.offerTemplateName, namespace);
  } catch (err) {
    // Ignore if not found
  }
}

/**
 * Provisions a Salad Node on the Golem Network by creating its Pod, Secret, and ConfigMaps.
 */
export async function provisionNode(
  k8sApi: CoreV1Api,
  namespace: string,
  node: Node
) {
  const names = generateNames(node.name);

  // Create or update ConfigMap for offer template
  const templateConfigMap: V1ConfigMap = {
    metadata: { name: names.offerTemplateName, namespace },
    data: {
      'template.json': JSON.stringify(node.offerTemplate),
    },
  };
  try {
    await k8sApi.readNamespacedConfigMap(names.offerTemplateName, namespace);
    await k8sApi.replaceNamespacedConfigMap(names.offerTemplateName, namespace, templateConfigMap);
  } catch {
    await k8sApi.createNamespacedConfigMap(namespace, templateConfigMap);
  }

  // Create or update ConfigMap for presets
  const presetsConfigMap: V1ConfigMap = {
    metadata: { name: names.presetsName, namespace },
    data: {
      'presets.json': JSON.stringify(node.presets),
    },
  };
  try {
    await k8sApi.readNamespacedConfigMap(names.presetsName, namespace);
    await k8sApi.replaceNamespacedConfigMap(names.presetsName, namespace, presetsConfigMap);
  } catch {
    await k8sApi.createNamespacedConfigMap(namespace, presetsConfigMap);
  }

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
          name: 'golem-provider',
          image: 'saladtechnologies/golem-provider:v0.17.6-1',
          imagePullPolicy: 'Never',
          args: ['run', '--no-interactive'],
          resources: {
            limits: { cpu: '2', memory: '1Gi' },
            requests: { cpu: '10m', memory: '64Mi' },
          },
          envFrom: [
            { secretRef: { name: names.environmentName } },
          ],
          volumeMounts: [
            {
              name: 'template-volume',
              mountPath: '/home/ubuntu/.config/ya-runtime-salad/template.json',
              subPath: 'template.json',
              readOnly: true,
            },
            {
              name: 'preset-volume',
              mountPath: '/home/ubuntu/.local/share/ya-provider/presets.json',
              subPath: 'presets.json',
              readOnly: true,
            },
          ],
        },
      ],
      volumes: [
        {
          name: 'template-volume',
          configMap: {
            name: names.offerTemplateName,
            items: [{ key: 'template.json', path: 'template.json' }],
          },
        },
        {
          name: 'preset-volume',
          configMap: {
            name: names.presetsName,
            items: [{ key: 'presets.json', path: 'presets.json' }],
          },
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
