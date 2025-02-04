#[cfg(test)]
mod shortinette_tests {
    use std::process::Command;

    use rand::{distributions::Alphanumeric, Rng};

    use ex06::*;

    fn create_binary() -> String {
        let rng = rand::thread_rng();

        let randname: String = rng
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();
        let filename = format!("/tmp/{}", randname);

        if let Err(e) = std::fs::write("/tmp/file.c", b"int main(){return 0;}") {
            panic!("Could not write file.c: {}.", e);
        }

        if let Err(e) = Command::new("cc")
            .arg("/tmp/file.c")
            .arg("-o")
            .arg(&filename)
            .output()
        {
            panic!("could not compile file.c: {}", e);
        }

        filename
    }

    fn sample_strings(path: &str, min: usize) -> String {
        match Command::new("strings")
            .arg("-n")
            .arg(min.to_string())
            .arg(path)
            .output()
        {
            Ok(out) => {
                String::from_utf8(out.stdout).expect("Could not parse UTF-8 from strings output.")
            }
            Err(e) => panic!("could not execute strings: {}.", e),
        }
    }

    // Calculates the similarity between two strings (0 => no similarity, 1 => same string).
    fn levenshtein_distance(a: &str, b: &str) -> f32 {
        let len_a = a.chars().count();
        let len_b = b.chars().count();

        let mut matrix: Vec<Vec<usize>> = vec![vec![0; len_b + 1]; len_a + 1];

        for idx in 0..len_a {
            matrix[idx][0] = idx;
        }

        for idx in 0..len_b {
            matrix[0][idx] = idx;
        }

        for (row, ca) in a.chars().enumerate() {
            for (col, cb) in b.chars().enumerate() {
                let cost = if ca == cb { 0 } else { 1 };

                matrix[row + 1][col + 1] = [
                    matrix[row][col + 1] + 1,
                    matrix[row + 1][col] + 1,
                    matrix[row][col] + cost,
                ]
                .iter()
                .min()
                .unwrap()
                .clone();
            }
        }

        ((len_a.max(len_b) as f32 - matrix[len_a][len_b] as f32) / len_a.max(len_b) as f32) as f32
    }

    #[test]
    fn basic() {
        let filename = create_binary();
        let master_output_str = sample_strings(&filename, 1)
            .lines()
            .filter(|line| !line.contains('\t'))
            .collect::<Vec<_>>()
            .join("\n");

        let mut student_output = Vec::new();

        if let Err(e) = strings(&mut student_output, &filename, false, None, None) {
            panic!(
                "Call to strings() returned an error on a standard use case: {}.",
                e
            );
        }

        let student_output_str =
            String::from_utf8(student_output).expect("Could not parse UTF-8 from output.");

        let levenshtein_distance = levenshtein_distance(&master_output_str, &student_output_str);

        assert!(
            levenshtein_distance > 0.99,
            "Similarity with 'strings -n 1 {}' too low (expected >= 0.99), got: {}.",
            filename,
            levenshtein_distance
        );
    }

    #[test]
    fn min() {
        let min_val = rand::random::<usize>().max(2).min(20);

        let filename = create_binary();
        let master_output_str = sample_strings(&filename, min_val)
            .lines()
            .filter(|line| !line.contains('\t'))
            .collect::<Vec<_>>()
            .join("\n");

        let mut student_output = Vec::new();

        if let Err(e) = strings(&mut student_output, &filename, false, Some(min_val), None) {
            panic!(
                "Call to strings() returned an error on a standard use case: {}.",
                e
            );
        }

        let student_output_str =
            String::from_utf8(student_output).expect("could not parse UTF-8 from output");

        for line in student_output_str.lines() {
            assert!(
                line.len() >= min_val,
                "Found string shorter than 'min' argument passed to 'strings()'. Expected: > {}, got: {}.",
                min_val,
                line.len()
            );
        }

        let levenshtein_distance = levenshtein_distance(&master_output_str, &student_output_str);

        // In my tests, the levenstein distance never went under 0.996 - If 0.99 turns out to be too strict, feel free
        // to make the expected range bigger.
        // - winstonallo
        assert!(
            levenshtein_distance > 0.99,
            "Similarity with 'strings -n 1 {}' too low (expected >= 0.99), got: {}",
            filename,
            levenshtein_distance
        );
    }

    #[test]
    fn max() {
        let max_val = rand::random::<usize>().max(2).min(20);

        let filename = create_binary();
        let master_output_str = sample_strings(&filename, max_val)
            .lines()
            .filter(|line| !line.contains('\t'))
            .filter(|line| !(line.len() >= max_val))
            .collect::<Vec<_>>()
            .join("\n");

        let mut student_output = Vec::new();

        if let Err(e) = strings(&mut student_output, &filename, false, None, Some(max_val)) {
            panic!("strings returned an error on a standard use case: {}", e);
        }

        let student_output_str =
            String::from_utf8(student_output).expect("could not parse UTF-8 from output");

        for line in student_output_str.lines() {
            assert!(
                line.len() <= max_val,
                "Found string shorter than 'min' argument passed to 'strings()'. Expected: >= {}, got: {}.",
                max_val,
                line.len()
            );
        }

        let levenshtein_distance = levenshtein_distance(&master_output_str, &student_output_str);

        // In my tests, the levenshtein distance never went under 0.996 - If 0.99 turns out
        // to be too strict, feel free to make the expected range bigger.
        // @winstonallo
        assert!(
            levenshtein_distance > 0.99,
            "Similarity with 'strings -n 1 {}' too low (expected >= 0.99), got: {}",
            filename,
            levenshtein_distance
        );
    }

    #[test]
    fn min_bigger_than_max() {
        let max_val = rand::random::<usize>();
        let min_val = rand::random::<usize>().max(max_val + 1);
        let mut output = Vec::new();

        let filename = create_binary();

        if let Ok(()) = strings(&mut output, &filename, false, Some(min_val), Some(max_val)) {
            panic!(
                "Call to 'strings()' did not return any error with min_val={} and max_val={}.",
                min_val, max_val
            );
        }
    }

    #[test]
    fn missing_permissions() {
        let mut output = Vec::new();

        let filename = create_binary();

        if let Err(e) = Command::new("chmod").arg("000").arg(&filename).output() {
            panic!("Could not chmod file '{}': {}.", filename, e);
        }

        if let Ok(()) = strings(&mut output, &filename, false, None, None) {
            panic!("Call to 'strings()' did not return any error with invalid permissions.");
        }
    }
}
