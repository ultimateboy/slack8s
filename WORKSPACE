http_archive(
    name = "io_bazel_rules_go",
    sha256 = "90bb270d0a92ed5c83558b2797346917c46547f6f7103e648941ecdb6b9d0e72",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.8.1/rules_go-0.8.1.tar.gz",
)

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_rules_dependencies",
    "go_register_toolchains",
    "go_repository",
)

go_rules_dependencies()

go_register_toolchains()

# You *must* import the Go rules before setting up the bazel_gazelle rules.
http_archive(
    name = "bazel_gazelle",
    sha256 = "e3dadf036c769d1f40603b86ae1f0f90d11837116022d9b06e4cd88cae786676",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.8/bazel-gazelle-0.8.tar.gz",
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# You *must* import the Go rules before setting up the go_image rules.
git_repository(
    name = "io_bazel_rules_docker",
    commit = "8aeab63328a82fdb8e8eb12f677a4e5ce6b183b1",
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
    container_repositories = "repositories",
)
load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

go_repository(
    name = "com_github_nlopes_slack",
    commit = "280245264a9d8c496e95ead447421d37cd95f93a",
    importpath = "github.com/nlopes/slack",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "292fd08b2560ad524ee37396253d71570339a821",
    importpath = "github.com/gorilla/websocket",
)
