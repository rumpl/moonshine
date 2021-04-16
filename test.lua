#syntax=rumpl/moonshine

---@diagnostic disable: undefined-global

from {
    base = "busybox",
    run {
        "echo hello",
        "echo world"
    },
    run {
        "ls -la"
    }
}
