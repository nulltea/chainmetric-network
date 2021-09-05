load("@io_bazel_rules_docker//container:image.bzl", "container_image")
load("@io_bazel_rules_docker//container:bundle.bzl", "container_bundle")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")

def multiarch_image(
        name,
        base,
        embed = [":go_default_library"],
        goarch = ["amd64", "arm64"],
        goos = ["linux"],
        user = "1000",
        stamp = True,
        **kwargs):
    for arch in goarch:
        for os in goos:
            go_image(
                name = "%s_%s-%s" % (name, os, arch),
                base = base,
                embed = embed,
                goarch = arch,
                goos = os,
                pure = "on",
            )

            container_image(
                name = "%s.%s-%s" % (name, os, arch),
                base = "%s_%s-%s" % (name, os, arch),
                user = user,
                stamp = stamp,
                **kwargs
            )

    container_image(
        name = name,
        base = "%s.%s-%s" % (name, goos[0], goarch[0]),
        **kwargs
    )

def multiarch_bundle(
        name,
        images,
        os = ["linux"],
        arch = ["amd64", "arm64"],
        **kwargs):
    all_images = {}
    for a in arch:
        for o in os:
            oa_images = {}
            for (k, v) in images.items():
                image_name = k.replace("{arch}", a)
                image_name = image_name.replace("{os}", o)

                oa_images[image_name] = "%s.%s-%s" % (v, o, a)

            container_bundle(
                name = "%s.%s-%s" % (name, o, a),
                images = oa_images,
                **kwargs
            )

            all_images += oa_images

    container_bundle(
        name = name,
        images = all_images,
        **kwargs
    )
