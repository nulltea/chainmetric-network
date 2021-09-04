load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

local_repository(
    name = "bazel_rules_go",
    path = "build/go/",
)

load("@bazel_rules_go//:repos.bzl", "add_go_repos")

add_go_repos()

load("@com_envoyproxy_protoc_gen_validate//bazel:repositories.bzl", "pgv_dependencies")

pgv_dependencies()

load("@bazel_rules_go//:def.bzl", "go_rules_deps")

go_rules_deps()

load("@com_envoyproxy_protoc_gen_validate//:dependencies.bzl", "go_third_party")

go_third_party()

load("//:go_third_party.bzl", "go_dependencies")

# gazelle:repository_macro go_third_party.bzl%go_dependencies
go_dependencies()
