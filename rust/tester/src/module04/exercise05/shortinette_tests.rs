#[cfg(test)]
mod shortinette_tests {
    use rand::seq::SliceRandom;

    use ex05::*;

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

    // Will not add the 42 intra here for obvious reasons.
    // Feel free to add stuff here, but please test it well, some websites use a lot of cookies which lead to a _lot_ of unique
    // values and can make 2 valid answers' similarity scores drop to < 0.5.
    fn get_random_website() -> &'static str {
        let websites = vec!["google.com", "github.com", "stackoverflow.com"];

        websites.choose(&mut rand::thread_rng()).unwrap()
    }

    fn sample_implementation<W: std::io::Write>(
        writer: &mut W,
        address: &str,
    ) -> Result<(), String> {
        let (addr, loc) = address.split_once("/").unwrap_or((address, ""));

        let mut stream = match TcpStream::connect((addr, 80)) {
            Ok(str) => str,
            Err(e) => {
                return Err(format!("TcpStream connection failed with error: {}.", e));
            }
        };

        match writeln!(
            stream,
            "\
            GET /{loc} HTTP/1.1\r\n\
            Host: {addr}\r\n\
            Connection: close\r\n\
            \r\n\
            "
        ) {
            Ok(_) => {}
            Err(e) => return Err(format!("writeln!() to TCP stream failed: {}.", e)),
        };

        let mut buf = [0u8; 4096];

        loop {
            let count = match stream.read(&mut buf) {
                Ok(c) => c,
                Err(e) => {
                    return Err(format!("Could not read from stream: {}.", e));
                }
            };

            if count == 0 {
                break;
            }

            match writer.write_all(&buf[..count]) {
                Ok(_) => {}
                Err(e) => {
                    return Err(format!("Could not write to writer: {}.", e));
                }
            }
        }

        Ok(())
    }

    #[test]
    fn basic() {
        let mut writer = Vec::new();

        if let Err(err) = get(&mut writer, &String::from("google.com")) {
            panic!("Call to get() failed with error: {}.", err);
        }

        let output_str = String::from_utf8(writer).expect("Could not parse UTF-8 from output.");

        assert!(output_str.contains("301") || output_str.contains("200"));
    }

    #[test]
    fn arthur_bied_charreton() {
        let mut writer = Vec::new();

        // This is the default location of my (winstonallo) server, which always redirects to some other service.
        // This output will stay the same, if not my poor soul is to blame for the inconvenience.
        if let Err(err) = get(&mut writer, &String::from("arthurbiedcharreton.com")) {
            panic!("Call to get() failed with error: {}.", err);
        }

        let output_str = String::from_utf8(writer).expect("Could not parse UTF-8 from output.");

        assert!(output_str.contains("HTTP/1.1"));
        assert!(output_str.contains("301 Moved Permanently"));
        assert!(output_str.contains("Server: nginx/1.24.0 (Ubuntu)"));
        assert!(output_str.contains("Content-Type: text/html"));
        assert!(output_str.contains("Connection: close"));
    }

    #[test]
    fn donotpanic() {
        let mut writer = Vec::new();

        if let Ok(()) = get(
            &mut writer,
            &String::from("ifthiswebsiteexistsiwilljuststopcoding.xcv"),
        ) {
            panic!("Call to get() with unexisting address return Ok(()).");
        }
    }

    #[test]
    fn hardcodingwillnotworksorry() {
        let website = get_random_website();

        let mut student_writer = Vec::new();
        let mut master_writer = Vec::new();

        if let Err(err) = get(&mut student_writer, &website) {
            panic!("Student implementation errored: {}.", err);
        }
        if let Err(err) = sample_implementation(&mut master_writer, website) {
            panic!("Master implementation errored: {}.", err);
        }

        let student_output_str =
            String::from_utf8(student_writer).expect("Could not parse UTF-8 from student output.");
        let master_output_str =
            String::from_utf8(master_writer).expect("Could not parse UTF-8 from master output.");

        let similarity = levenshtein_distance(&student_output_str, &master_output_str);

        assert!(similarity > 0.95, "Similarity with sample implementation too low for address `{}`, expected >= 0.95, got: {}.", website, similarity);
    }
}
