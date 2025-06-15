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
