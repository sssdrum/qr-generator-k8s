1. Create backend ✓

2. Create frontend ✓

3. Dockerize backend ✓

4. Dockerize frontend ✓

5. Build images from dockerfiles ✓

   - Why? Minikube has its own docker environment that is not shared with your local images
   - How?
     - run `minikube start` to start cluster
     - run `eval $(minikube docker-env)` to use to minikube's docker engine by default
     - build images

6. Create Kubernetes manifests for ✓

   - Deployments
   - Services

   Optional: - Port forward to hosts's ports to check if it works and access pods logs

7. Add a postgres DB to cache qr codes as a StatefulSet

8. Build an operator in go/ansible

9. Get grilled by claude about design

# Common commands

## kubectl

- Apply deployment `kubectl apply -f <path to manifest>`
- See deployments: `kubectl get deployment`
- See pods in a deployment: `kubectl get pods -l app=<deployment>`
- See logs in a pod: `kubectl logs -f <pod name>`
- Port forward a service to access on host: `kubectl port-forward svc/<your service> <service port>:<your port>`

## minikube

Get public url of service `minikube service frontend-service --url`

## Applying changes

1. Edit code
2. Rebuild image `docker build -t <image> <path to image>`
3. Apply changes `kubectl rollout restart deployment <deployment>`

# Secrets

# Debugging

App A does not reach app B. How to debug?

1. Is app B running? `kubectl get pods/services`
2. Check logs in app B `kubectl logs -l app=<app B>`
3. Try curling form an existing pod or create a temporary debug pod
