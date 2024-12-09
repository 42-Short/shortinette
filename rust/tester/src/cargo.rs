use std::{
    fs,
    io::{self, BufRead, Write},
    path, process,
};

pub struct Cargo {
    dir: tempfile::TempDir,
}

impl Cargo {
    pub fn new(name: &str, is_lib: bool) -> Self {
        let dir = tempfile::tempdir().expect("Failed to create directory for cargo project");

        let mut init_command = process::Command::new("cargo");

        init_command
            .arg("init")
            .arg(dir.path())
            .args(["--name", name]);

        if is_lib {
            init_command.arg("--lib");
        }

        let init_output = init_command
            .output()
            .expect("Failed to execute cargo command");

        assert!(
            init_output.status.success(),
            "Failed to cargo init a new project"
        );

        Self { dir }
    }

    pub fn copy_from(path: &path::PathBuf) -> Self {
        let dir = tempfile::tempdir().expect("Failed to create directory for cargo project");

        fs_extra::dir::copy(
            path,
            dir.path(),
            &fs_extra::dir::CopyOptions::new().content_only(true),
        )
        .expect("Failed to copy content of cargo project");

        // Just to go sure
        _ = fs::remove_dir(dir.path().join("target"));

        Self { dir }
    }

    pub fn compile(&self) -> Result<path::PathBuf, ()> {
        let compile_output = process::Command::new("cargo")
            .current_dir(self.dir.path())
            .arg("build")
            .arg("--release")
            .arg("--message-format=json")
            .output()
            .expect("Failed to execute cargo build");

        if !compile_output.status.success() {
            return Err(());
        }

        let reader = io::BufReader::new(compile_output.stdout.as_slice());

        #[derive(Debug, serde::Deserialize)]
        struct CompilerOutput {
            reason: String,
            executable: String,
        }

        let executable_path = reader
            .lines()
            .map_while(Result::ok)
            .map(|line| serde_json::from_str::<CompilerOutput>(&line))
            .filter_map(Result::ok)
            .filter(|output| output.reason.as_str() == "compiler-artifact")
            .last()
            .map(|output| path::PathBuf::from(output.executable))
            .expect("Failed to get path to compiled exercise");

        Ok(executable_path)
    }

    pub fn path(&self) -> &path::Path {
        self.dir.path()
    }

    pub fn add_dependency(&self, name: &str, version: &str) -> Result<(), ()> {
        let add_output = process::Command::new("cargo")
            .current_dir(self.dir.path())
            .arg("add")
            .arg(name)
            .arg("--version")
            .arg(version)
            .output()
            .expect("Failed to execute cargo add");

        if add_output.status.success() {
            Ok(())
        } else {
            Err(())
        }
    }

    pub fn add_local_dependency(&self, path: &str) -> Result<(), ()> {
        let add_output = process::Command::new("cargo")
            .current_dir(self.dir.path())
            .arg("add")
            .arg("--path")
            .arg(path)
            .output()
            .expect("Failed to execute cargo add");

        if add_output.status.success() {
            Ok(())
        } else {
            Err(())
        }
    }

    pub fn run_test<'a>(
        &self,
        test_filter: impl IntoIterator<Item = &'a str>,
    ) -> Result<(), String> {
        let test_output = process::Command::new("cargo")
            .current_dir(self.dir.path())
            .arg("test")
            .args(test_filter)
            .output()
            .expect("Failed to execute cargo test");

        if test_output.status.success() {
            return Ok(());
        };

        let test_log = String::from_utf8_lossy(&test_output.stdout);

        Err(test_log.into_owned())
    }

    pub fn check_clippy(&self, config: &str) -> bool {
        let config_file_path = self.dir.path().join("clippy.toml");
        let mut config_file =
            fs::File::create(&config_file_path).expect("Failed to create clippy.toml file");

        config_file
            .write_all(config.as_bytes())
            .expect("Failed to write clippy.toml file");

        let clippy_output = process::Command::new("cargo")
            .current_dir(self.dir.path())
            .arg("clippy")
            .arg("--")
            .arg("-D")
            .arg("warnings")
            .output()
            .expect("Failed to execute cargo clippy");

        _ = fs::remove_file(config_file_path);

        clippy_output.status.success()
    }

    pub fn test_list(&self, test_filter: &[&str]) -> Vec<String> {
        let test_list_output = process::Command::new("cargo")
            .current_dir(self.dir.path())
            .arg("test")
            .args(test_filter)
            .arg("--")
            .arg("--list")
            // `--format=json` is a nightly feature...
            .arg("--format=terse")
            .output()
            .expect("Failed to execute cargo test --list");

        assert!(
            test_list_output.status.success(),
            "cargo test --list failed"
        );

        let out = String::from_utf8_lossy(&test_list_output.stdout);

        out.lines()
            .map(|line| {
                &line[..line
                    .rfind(": test")
                    .expect("Test list is not correctly formatted")]
            })
            .map(String::from)
            .collect()
    }
}
