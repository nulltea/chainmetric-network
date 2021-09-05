load("//build/go:repos.bzl", "go_repos")

go_repos()

load("@com_envoyproxy_protoc_gen_validate//bazel:repositories.bzl", "pgv_dependencies")

pgv_dependencies()

load("//build/go:def.bzl", "go_rules")

go_rules()

load("//build/docker:repos.bzl", "docker_repos")

docker_repos()

load("//build/docker:def.bzl", "docker_deps")

docker_deps()

load("//:go_third_party.bzl", "go_dependencies")

# gazelle:repository_macro go_third_party.bzl%go_dependencies
go_dependencies()
