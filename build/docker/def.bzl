load("@io_bazel_rules_docker//repositories:repositories.bzl", container_repositories = "repositories")
load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")
load("@io_bazel_rules_docker//go:image.bzl", _go_image_repos = "repositories")
load("@io_bazel_rules_docker//container:container.bzl", "container_pull")
load("@io_bazel_rules_docker//toolchains/docker:toolchain.bzl", "toolchain_configure")

def extra_container_repos():
    container_pull(
        name = "alpine",
        registry = "index.docker.io",
        repository = "alpine",
        tag = "3.14",
    )

def docker_deps():
    toolchain_configure(
        name = "docker_config",
        docker_path = "/usr/bin/docker",
    )
    container_repositories()
    container_deps()
    _go_image_repos()
    extra_container_repos()
