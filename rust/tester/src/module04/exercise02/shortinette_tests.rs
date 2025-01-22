#[cfg(test)]
mod shortinette_tests {
    use ex02::*;
    use rand::{distributions::Alphanumeric, Rng};
    use std::{fs, path::Path};

    fn create_test_dir(
        base_path: &str,
        num_files: usize,
        max_size: usize,
    ) -> Result<Vec<usize>, String> {
        let base = Path::new(base_path);

        if base.exists() {
            fs::remove_dir_all(base).map_err(|e| format!("Failed to clean old test dir: {}", e))?;
        }
        fs::create_dir(base).map_err(|e| format!("Failed to create test dir: {}", e))?;

        let mut rng = rand::thread_rng();
        let mut file_sizes = Vec::new();

        for i in 0..num_files {
            let random_size = rng.gen_range(1..=max_size);
            file_sizes.push(random_size);

            let file_path = base.join(format!("file_{}", i));
            let content = vec![0u8; random_size];
            fs::write(&file_path, &content)
                .map_err(|e| format!("Failed to create file {}: {}", file_path.display(), e))?;
        }

        Ok(file_sizes)
    }

    fn approximately_equal(a: f64, b: f64, epsilon: f64) -> bool {
        (a - b).abs() < epsilon
    }

    #[test]
    fn basic() {
        let rng = rand::thread_rng();

        let filename: String = rng
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let full_path = format!("/tmp/{filename}");
        let num_files = 10;
        let max_size = 16384;
        let file_sizes = create_test_dir(&full_path, num_files, max_size)
            .expect("Failed to create test directory.");

        let mut output = Vec::new();
        duh(&mut output, &full_path).expect("Call to duh() failed.");

        let output_str = String::from_utf8(output).expect("Output is not valid UTF-8.");

        let expected_total_size: u64 = file_sizes.iter().map(|&s| s as u64).sum();

        let expected_size_str = if expected_total_size < 1000 {
            format!("{} bytes", expected_total_size)
        } else if expected_total_size < 1_000_000 {
            format!("{:.1} kilobytes", expected_total_size as f64 / 1000.0)
        } else if expected_total_size < 1_000_000_000 {
            format!("{:.1} megabytes", expected_total_size as f64 / 1_000_000.0)
        } else {
            format!(
                "{:.1} gigabytes",
                expected_total_size as f64 / 1_000_000_000.0
            )
        };

        assert_eq!(output_str.trim(), expected_size_str);

        fs::remove_dir_all(&full_path).expect("Failed to clean up test directory.");
    }

    #[test]
    fn recursivity() {
        let rng = rand::thread_rng();

        let dirname: String = rng
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let test_dir = format!("/tmp/{}", dirname);

        let num_files = 6;
        let max_size = 8192;
        let file_sizes = create_test_dir(&test_dir, num_files, max_size)
            .expect("Failed to create test directory");

        let mut total_expected_size: u64 = file_sizes.iter().map(|&size| size as u64).sum();

        // Stack overflow possible in some rare cases with depth >= 200
        let depth = (rand::random::<u8>().saturating_add(2)).min(150);
        let mut current_path = test_dir.clone();

        for _ in 0..depth {
            let subdir_name: String = rng
                .to_owned()
                .sample_iter(&Alphanumeric)
                .take(2)
                .map(char::from)
                .collect();

            current_path = format!("{}/{}", current_path, subdir_name);

            let num_files = (rand::random::<u8>().saturating_add(2)).min(10);
            let file_sizes = create_test_dir(&current_path, num_files as usize, max_size)
                .expect("Failed to create subdirectory.");

            let subdir_total_size: u64 = file_sizes.iter().map(|&size| size as u64).sum();

            total_expected_size += subdir_total_size;

            #[cfg(target_os = "macos")]
            // MacOS directory size depends on the number of files it contains
            let expected_dir_size = 64 + (32 * num_files as u64);

            #[cfg(target_os = "linux")]
            let expected_dir_size = 4096;

            total_expected_size += expected_dir_size;
        }

        let mut output = Vec::new();

        duh(&mut output, &test_dir).expect("Call to duh() failed.");

        let output_str = String::from_utf8(output).expect("Output is not valid UTF-8.");

        let (expected_size_float, unit) = match total_expected_size {
            0..1_000 => (total_expected_size as f64, "bytes"),
            1_000..1_000_000 => (total_expected_size as f64 / 1_000.0, "kilobytes"),
            1_000_000..1_000_000_000 => (total_expected_size as f64 / 1_000_000.0, "megabytes"),
            _ => (total_expected_size as f64 / 1_000_000.0, "gigabytes"),
        };

        let actual_size_float = output_str
            .trim()
            .split_whitespace()
            .next()
            .unwrap()
            .parse::<f64>()
            .expect("Failed to parse output size.");

        assert!(
            output_str.contains(unit),
            "Recursion test with depth {} failed - Expected unit: {}, Got: {}.",
            depth,
            expected_size_float,
            actual_size_float
        );

        // Float rounding errors make the tests not 100% deterministic, this lazy solution
        // should fix it according to my empirical study
        assert!(
            approximately_equal(expected_size_float, actual_size_float, 1.0),
            "Recursion test with depth {} failed - Expected size: {}, Got: {}.",
            depth,
            expected_size_float,
            actual_size_float
        );

        fs::remove_dir_all(&test_dir).expect("Failed to clean up test directory.");
    }

    #[test]
    fn nonexisting() {
        let mut output = Vec::new();
        if let Ok(()) = duh(&mut output, "idonotexist") {
            panic!("duh did not return Err on non-existing base directory path.")
        }
    }
}
