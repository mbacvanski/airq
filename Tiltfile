# -*- mode: Go -*-

docker_build('airquality', '.', dockerfile='deployments/Dockerfile')
k8s_yaml('deployments/kubernetes.yaml')
k8s_resource('airquality', port_forwards=8080)
