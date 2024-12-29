#[cfg(test)]
mod tests {
    use super::*;
    use rand::{distributions::Alphanumeric, Rng};
    use std::{fs, io::Cursor, process::Command};

    #[test]
    fn basic() {
        let rng = rand::thread_rng();

        let filenames: Vec<String> = (0..10)
            .map(|_| {
                let name: String = rng
                    .to_owned()
                    .sample_iter(&Alphanumeric)
                    .take(10)
                    .map(char::from)
                    .collect();
                format!("/tmp/{name}")
            })
            .collect();

        let content: String = rng
            .sample_iter(&Alphanumeric)
            .take(20)
            .map(char::from)
            .collect();

        let mut input = Cursor::new(content.as_bytes());
        let mut output = Vec::new();

        tee(&mut input, &mut output, filenames.clone());

        assert_eq!(output, content.as_bytes());

        for filename in &filenames {
            let file_content = fs::read_to_string(filename).unwrap();
            assert_eq!(file_content, content);

            fs::remove_file(filename).unwrap();
        }
    }

    #[test]
    fn is_dir() {
        let rng = rand::thread_rng();

        let dirname: String = rng
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let full_path = format!("/tmp/{dirname}");

        fs::create_dir(&full_path).expect("failed to create directory");

        let mut input = Cursor::new("Don't panic!".as_bytes());
        let mut output = Vec::new();

        tee(&mut input, &mut output, vec![full_path.clone()]);

        let _ = fs::remove_dir(&full_path);
    }

    #[test]
    fn no_perm() {
        let rng = rand::thread_rng();

        let filename: String = rng
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let full_path = format!("/tmp/{filename}");

        let _ = fs::File::create_new(&full_path);

        let _ = Command::new("chmod")
            .arg("000")
            .arg(&full_path)
            .output()
            .expect("failed to execute process");

        let mut input = Cursor::new("Don't panic!".as_bytes());
        let mut output = Vec::new();

        tee(&mut input, &mut output, vec![full_path.clone()]);
    }

    #[test]
    fn overwrite() {
        let rng = rand::thread_rng();

        let filename: String = rng
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let full_path = format!("/tmp/{filename}");

        let _ = fs::write(&full_path, b"Overwrite me!").unwrap();

        let content: String = rng
            .sample_iter(&Alphanumeric)
            .take(20)
            .map(char::from)
            .collect();

        let mut input = Cursor::new(content.as_bytes());
        let mut output = Vec::new();

        tee(&mut input, &mut output, vec![full_path.clone()]);

        let file_content = fs::read_to_string(&full_path).unwrap();
        assert_eq!(file_content, content);

        fs::remove_file(&full_path).unwrap();
    }
}
