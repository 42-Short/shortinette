#[cfg(test)]
mod shortinette_tests_rust_0105 {
    use ex05::deduplicate;
    use rand::Rng;

    fn check_deduplicated(input: &[i32], deduplicated: &[i32]) -> bool {
        let mut input_vec = Vec::new();
        let mut deduplicated_vec = Vec::new();

        for nbr in input {
            if !input_vec.contains(&nbr) {
                input_vec.push(nbr);
            }
        }

        for nbr in deduplicated {
            if !deduplicated_vec.contains(&nbr) {
                deduplicated_vec.push(nbr);
            } else {
                return false;
            }
        }

        input_vec == deduplicated_vec
    }

    #[test]
    fn test_empty() {
        let mut v = vec![];
        deduplicate(&mut v);
        assert_eq!(v, [], "Error with empty vec as input");
    }

    #[test]
    fn test_subject() {
        let mut v = vec![1, 2, 2, 3, 2, 4, 3];
        deduplicate(&mut v);
        assert_eq!(v, [1, 2, 3, 4]);
    }

    #[test]
    fn test_only_duplicates() {
        let mut rng = rand::thread_rng();

        let value = rng.gen_range(i32::MIN..=i32::MAX);
        let mut v: Vec<i32> = (0..100).map(|_| value).collect();
        deduplicate(&mut v);
        assert_eq!(
            v,
            [value],
            "Error with vec which only contains the same value",
        );
    }

    #[test]
    fn test_mixed() {
        let mut rng = rand::thread_rng();

        let mut v: Vec<i32> = (0..1000).map(|_| rng.gen_range(-100..=100)).collect();
        let input = v.clone();

        deduplicate(&mut v);
        assert!(
            check_deduplicated(&input, &v),
            "Incorrect result\nInput: {:?}\nOutput: {:?}",
            input,
            v
        );
    }
}
