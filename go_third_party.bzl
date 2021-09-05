load("@bazel_gazelle//:deps.bzl", "go_repository")

# Generate Go dependencies macro with:
# gazelle update-repos --from_file=go.mod -index=false -to_macro=go_third_party.bzl%go_dependencies
def go_dependencies():
    pass
