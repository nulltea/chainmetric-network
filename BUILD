load("@bazel_gazelle//:def.bzl", "gazelle")
load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle_binary")
load("@io_bazel_rules_go//proto:compiler.bzl", "go_proto_compiler")

gazelle_binary(
    name = "gazelle_binary",
    languages = DEFAULT_LANGUAGES + [
        "@golink//gazelle/go_link:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

# gazelle:prefix github.com/timoth-y/chainmetric-network
# gazelle:go_grpc_compilers @io_bazel_rules_go//proto:go_grpc
# gazelle:go_proto_compilers @io_bazel_rules_go//proto:go_grpc,//:go_validate
# gazelle:resolve proto validate/validate.proto @com_envoyproxy_protoc_gen_validate//validate:validate_proto
# gazelle:resolve proto go validate/validate.proto @com_envoyproxy_protoc_gen_validate//validate:validate_go_proto
gazelle(
    name = "gazelle",
    gazelle = "//:gazelle_binary",
)

go_proto_compiler(
    name = "go_validate",
    options = [
        "lang=go",
    ],
    plugin = "@com_envoyproxy_protoc_gen_validate//:protoc-gen-validate",
    suffix = ".pb.validate.go",
    valid_archive = False,
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
