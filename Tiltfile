"""
Tiltfile for deploying services and tooling in development environment.
"""

load("ext://helm_resource", "helm_repo", "helm_resource")
load("ext://restart_process", "docker_build_with_restart")

def deploy_service(service_name, port_forwards = "", resource_deps = [], values_files = []):
    """Deploy a service using Tilt.

    Args:
        service_name (str): The name of the service to deploy.
        port_forwards (list, optional): List of port forwards for the service. Defaults to "".
        resource_deps (list, optional): List of dependencies for the service. Defaults to [].
        values_files (list, optional): List of Helm values files for the service. Defaults to [].
    """

    # Set default dependencies
    if not resource_deps:
        resource_deps = ["consul"]

    # Set default values files
    if not values_files:
        values_files = ["./infra/helm/values/dev/{}-values.yaml".format(service_name)]

    # Compilation
    compile_cmd = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/{} ./services/{}/cmd/main.go".format(service_name, service_name)
    if os.name == "nt":
        compile_cmd = "./infra/docker/dev/{}-build.bat".format(service_name)

    local_resource(
        "{}-compile".format(service_name),
        compile_cmd,
        deps = ["./services/{}".format(service_name), "./shared"],
        labels = "compiles",
    )

    # Docker build
    docker_build_with_restart(
        "vasapolrittideah/moneylog-api-{}".format(service_name),
        ".",
        entrypoint = ["/app/build/{}".format(service_name)],
        dockerfile = "./infra/docker/dev/{}.Dockerfile".format(service_name),
        only = ["./build/{}".format(service_name), "./shared"],
        live_update = [
            sync("./build", "/app/build"),
            sync("./shared", "/app/shared"),
        ],
    )

    # Kubernetes deployment
    k8s_yaml(helm(
        "./infra/helm/charts/{}".format(service_name),
        name = service_name,
        values = values_files,
    ))

    # Kubernetes resource configuration
    compile_dep = "{}-compile".format(service_name)
    all_deps = [compile_dep] + resource_deps

    if port_forwards:
        k8s_resource(
            service_name,
            resource_deps = all_deps,
            port_forwards = port_forwards,
            labels = "services",
        )
    else:
        k8s_resource(
            service_name,
            resource_deps = all_deps,
            labels = "services",
        )
