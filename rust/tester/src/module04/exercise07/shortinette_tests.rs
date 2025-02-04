#[cfg(test)]
mod shortinette_tests {
    use std::process::Command;

    use rand::distributions::Alphanumeric;

    use ex07::*;

    fn generate_random_keys() -> (String, String) {
        let pub_key_path: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();
        let pub_key_path: String = format!("/tmp/{}.priv", pub_key_path);
        let priv_key_path: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();
        let priv_key_path: String = format!("/tmp/{}.priv", priv_key_path);

        if let Err(e) = gen_keys(&pub_key_path, &priv_key_path) {
            panic!("Call to 'gen_keys()' returned an error: {}.", e);
        }

        (pub_key_path, priv_key_path)
    }

    #[test]
    fn keygen_broken_paths() {
        let pub_key_path: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();
        let priv_key_path: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();
        let pub_key_path: String = format!("/{}/{}", pub_key_path, priv_key_path);
        let priv_key_path: String = format!(
            "/{}/{}/{}/{}",
            priv_key_path, pub_key_path, pub_key_path, priv_key_path
        );

        if let Ok(()) = gen_keys(&pub_key_path, &priv_key_path) {
            panic!("Call to 'gen_keys()' returned Ok(()) with broken key paths.");
        }
    }

    #[test]
    fn basic_keygen() {
        // This will panic if anything goes wrong.
        let _ = generate_random_keys();
    }

    #[test]
    fn you_shall_not_bullshit() {
        let randstring: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(rand::random::<u16>() as usize)
            .map(char::from)
            .collect();

        let mut input = randstring.as_bytes();
        let mut encrypt_output = Vec::new();

        let (pub_key_path, _) = generate_random_keys();

        if let Err(e) = encrypt(&mut input, &mut encrypt_output, &pub_key_path) {
            panic!("Call to 'encrypt()' failed with error: {}.", e);
        }

        let encrypted = encrypt_output.as_slice();

        // Make sure students are not just copying the input into the encryption.
        assert_ne!(encrypted, input);
    }

    #[test]
    fn encrypt_and_decrypt() {
        let randstring: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(rand::random::<u16>() as usize)
            .map(char::from)
            .collect();

        let mut input = randstring.as_bytes();
        let mut encrypt_output = Vec::new();

        let (pub_key_path, priv_key_path) = generate_random_keys();

        if let Err(e) = encrypt(&mut input, &mut encrypt_output, &pub_key_path) {
            panic!("Call to 'encrypt()' failed with error: {}.", e);
        }

        let mut input = encrypt_output.as_slice();
        let mut decrypt_output = Vec::new();

        if let Err(e) = decrypt(&mut input, &mut decrypt_output, &priv_key_path) {
            panic!("Call to 'decrypt()' failed with error: {}.", e);
        }

        let decrypted_str =
            String::from_utf8(decrypt_output).expect("Could not parse UTF-8 from output.");

        assert_eq!(decrypted_str, randstring);
    }

    #[test]
    fn decrypt_invalid_priv_key() {
        let randstring: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(8)
            .map(char::from)
            .collect();

        let mut input = randstring.as_bytes();
        let mut encrypt_output = Vec::new();

        let (pub_key_path, _) = generate_random_keys();
        let (_, priv_key_path) = generate_random_keys();

        if let Err(e) = encrypt(&mut input, &mut encrypt_output, &pub_key_path) {
            panic!("Call to 'encrypt()' failed with error: {}.", e);
        }

        let mut input = encrypt_output.as_slice();
        let mut decrypt_output = Vec::new();

        let _ = decrypt(&mut input, &mut decrypt_output, &priv_key_path);

        assert_ne!(decrypt_output, randstring.as_bytes());
    }

    #[test]
    fn encrypt_nonexisting_pub_key() {
        let randstring: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(8)
            .map(char::from)
            .collect();

        let mut input = randstring.as_bytes();
        let mut encrypt_output = Vec::new();

        if let Ok(()) = encrypt(&mut input, &mut encrypt_output, &randstring) {
            panic!("Call to 'encrypt()' did not fail with non-existing public key.");
        }
    }

    #[test]
    fn decrypt_nonexisting_pub_key() {
        let randstring: String = rand::thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(8)
            .map(char::from)
            .collect();

        let mut input = randstring.as_bytes();
        let mut encrypt_output = Vec::new();

        let (pub_key_path, _) = generate_random_keys();

        if let Err(e) = encrypt(&mut input, &mut encrypt_output, &pub_key_path) {
            panic!("Call to 'encrypt()' failed with error: {}.", e);
        }

        let mut decrypt_input = encrypt_output.as_slice();
        let mut decrypt_output = Vec::new();
        if let Ok(()) = decrypt(&mut decrypt_input, &mut decrypt_output, &randstring) {
            panic!("Call to 'decrypt()' did not fail with non-existing private key.");
        }
    }
}
