const std = @import("std");

const Target = struct {
    goos: []const u8,
    goarch: []const u8,
};

const targets = [_]Target{
    .{ .goos = "windows", .goarch = "amd64" },
    .{ .goos = "windows", .goarch = "arm64" },
    .{ .goos = "linux", .goarch = "amd64" },
    .{ .goos = "linux", .goarch = "arm64" },
};

pub fn build(b: *std.Build) void {
    const filter_goos = b.option([]const u8, "goos", "Target OS (windows|linux)");
    const filter_goarch = b.option([]const u8, "goarch", "Target arch (amd64|arm64)");

    for (targets) |t| {
        if (filter_goos) |fg| {
            if (!std.mem.eql(u8, fg, t.goos)) continue;
        }
        if (filter_goarch) |fa| {
            if (!std.mem.eql(u8, fa, t.goarch)) continue;
        }

        addGoTarget(b, t);
    }
}

fn addGoTarget(b: *std.Build, t: Target) void {
    const name = b.fmt("{s}-{s}", .{ t.goos, t.goarch });

    const host_goos = hostGOOS();
    const host_goarch = hostGOARCH();
    const is_host = std.mem.eql(u8, t.goos, host_goos) and
        std.mem.eql(u8, t.goarch, host_goarch);

    const ext = if (std.mem.eql(u8, t.goos, "windows")) ".exe" else "";
    const out_bin = b.fmt("dist/{s}/conduit{s}", .{ name, ext });

    const build_step = b.addSystemCommand(&.{ "go", "build", "-o", out_bin, "." });
    build_step.setEnvironmentVariable("GOOS", t.goos);
    build_step.setEnvironmentVariable("GOARCH", t.goarch);

    const build_named = b.step(
        b.fmt("build-{s}", .{name}),
        b.fmt("Build conduit for {s}/{s}", .{ t.goos, t.goarch }),
    );
    build_named.dependOn(&build_step.step);
    b.default_step.dependOn(&build_step.step);

    if (is_host) {
        const test_step = b.addSystemCommand(&.{ "go", "test", "./...", "-count=1", "-v" });
        test_step.setEnvironmentVariable("GOOS", t.goos);
        test_step.setEnvironmentVariable("GOARCH", t.goarch);

        const test_named = b.step(
            b.fmt("test-{s}", .{name}),
            b.fmt("Run tests for {s}/{s}", .{ t.goos, t.goarch }),
        );
        test_named.dependOn(&test_step.step);
        b.default_step.dependOn(&test_step.step);
    }
}

fn hostGOOS() []const u8 {
    return switch (@import("builtin").os.tag) {
        .windows => "windows",
        .linux => "linux",
        .macos => "darwin",
        else => "linux",
    };
}

fn hostGOARCH() []const u8 {
    return switch (@import("builtin").cpu.arch) {
        .x86_64 => "amd64",
        .aarch64 => "arm64",
        else => "amd64",
    };
}
