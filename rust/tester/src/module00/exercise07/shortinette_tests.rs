#[cfg(test)]
mod shortinette_tests_0007 {
    use ex07::strpcmp;
    use rand::{distributions::Alphanumeric, Rng};

    fn random_u8_vec(size: usize) -> Vec<u8> {
        let mut rng = rand::thread_rng();
        (0..size).map(|_| rng.sample(Alphanumeric)).collect()
    }

    #[test]
    fn equal_test() {
        let random_input = random_u8_vec(32);
        assert!(strpcmp(&random_input, &random_input));
    }

    #[test]
    fn additional_wildcard_start() {
        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input.clone();
        with_wildcard.insert(0, b'*');

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn replace_wildcard_start() {
        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input.clone();
        with_wildcard[0] = b'*';

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn additional_multiple_wildcard_start() {
        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input.clone();
        with_wildcard.splice(0..0, std::iter::repeat(b'*').take(5));

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn one_wildcard_replace_multiple_start() {
        let mut rng = rand::thread_rng();

        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input[rng.gen_range(5..20)..].to_vec();
        with_wildcard.insert(0, b'*');

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn additional_wildcard_end() {
        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input.clone();
        with_wildcard.push(b'*');

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn replace_wildcard_end() {
        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input[..random_input.len() - 1].to_vec();
        with_wildcard.push(b'*');

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn additional_multiple_wildcard_end() {
        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input.clone();
        with_wildcard.extend(std::iter::repeat(b'*').take(5));

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn one_wildcard_replace_multiple_end() {
        let mut rng = rand::thread_rng();

        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input[..rng.gen_range(10..25)].to_vec();
        with_wildcard.push(b'*');

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn wildcard_middle_replace() {
        let mut rng = rand::thread_rng();

        let random_input = random_u8_vec(32);
        let position = rng.gen_range(1..random_input.len() - 1);

        let mut with_wildcard = random_input.clone();
        with_wildcard[position] = b'*';

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn wildcard_middle_additional() {
        let mut rng = rand::thread_rng();

        let random_input = random_u8_vec(32);
        let position = rng.gen_range(1..random_input.len() - 1);

        let mut with_wildcard = random_input.clone();
        with_wildcard.insert(position, b'*');

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn additional_multiple_wildcard_middle() {
        let mut rng = rand::thread_rng();

        let random_input = random_u8_vec(32);
        let position = rng.gen_range(1..random_input.len() - 1);

        let mut with_wildcard = random_input.clone();
        with_wildcard.splice(position..position, std::iter::repeat(b'*').take(5));
    }

    #[test]
    fn one_wildcard_replace_multiple_middle() {
        let mut rng = rand::thread_rng();

        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input[0..rng.gen_range(5..10)].to_vec();
        with_wildcard.push(b'*');
        with_wildcard.extend(&random_input[rng.gen_range(20..25)..]);

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn start_end_wildcard() {
        let mut rng = rand::thread_rng();

        let random_input = random_u8_vec(32);
        let mut with_wildcard = random_input[rng.gen_range(5..10)..rng.gen_range(15..20)].to_vec();
        with_wildcard.insert(0, b'*');
        with_wildcard.push(b'*');

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn multiple_substrings() {
        let part1 = random_u8_vec(8);
        let part2 = random_u8_vec(8);

        let random_input: Vec<u8> = part1
            .iter()
            .chain(std::iter::repeat(&part2).take(5).flat_map(|v| v))
            .copied()
            .collect();

        let with_wildcard: Vec<u8> = part1.iter().chain(b"*").chain(&part2).copied().collect();

        assert!(strpcmp(&random_input, &with_wildcard));
    }

    #[test]
    fn empty_query_only_wildcards() {
        let mut rng = rand::thread_rng();

        let pattern: Vec<u8> = std::iter::repeat(b'*').take(rng.gen_range(5..10)).collect();
        assert!(strpcmp(b"", &pattern));
    }

    #[test]
    fn both_empty() {
        assert!(strpcmp(b"", b""));
    }

    #[test]
    fn wildcard_in_query_empty_pattern() {
        let mut rng = rand::thread_rng();

        let query: Vec<u8> = std::iter::repeat(b'*').take(rng.gen_range(5..10)).collect();
        assert_eq!(strpcmp(&query, b""), false);
    }

    #[test]
    fn start_mismatch() {
        let random_input = random_u8_vec(32);
        let with_wildcard: Vec<u8> = random_input[1..].iter().chain(b"*").copied().collect();

        assert_eq!(strpcmp(&random_input, &with_wildcard), false);
    }

    #[test]
    fn end_mismatch() {
        let random_input = random_u8_vec(32);
        let with_wildcard: Vec<u8> = [vec![b'*'], random_input[..31].to_vec()]
            .iter()
            .flat_map(|v| v)
            .copied()
            .collect();

        assert_eq!(strpcmp(&random_input, &with_wildcard), false);
    }

    #[test]
    fn same_amount_of_wildcards() {
        let mut rng = rand::thread_rng();

        let part1 = random_u8_vec(8);
        let part2 = random_u8_vec(8);

        let input: Vec<u8> = part1
            .iter()
            .chain(std::iter::repeat(&b'*').take(rng.gen_range(1..5)))
            .chain(&part2)
            .copied()
            .collect();

        assert!(strpcmp(&input, &input));
    }

    #[test]
    fn more_wildcards_in_query() {
        let mut rng = rand::thread_rng();

        let part1 = random_u8_vec(8);
        let part2 = random_u8_vec(8);
        let wildcard_count = rng.gen_range(1..5);

        let query: Vec<u8> = part1
            .iter()
            .chain(std::iter::repeat(&b'*').take(wildcard_count + 1))
            .chain(&part2)
            .copied()
            .collect();

        let pattern: Vec<u8> = part1
            .iter()
            .chain(std::iter::repeat(&b'*').take(wildcard_count))
            .chain(&part2)
            .copied()
            .collect();

        assert!(strpcmp(&query, &pattern));
    }

    #[test]
    fn wildcard_wrong_side() {
        let random_input = random_u8_vec(32);

        let query: Vec<u8> = random_input.iter().chain(b"*").copied().collect();
        let pattern: Vec<u8> = [b'*'].iter().chain(&random_input).copied().collect();

        assert_eq!(strpcmp(&query, &pattern), false);
    }

    #[test]
    fn wildcard_in_query_none_in_pattern() {
        let part1 = random_u8_vec(8);
        let part2 = random_u8_vec(8);

        let query: Vec<u8> = part1.iter().chain(b"*").chain(&part2).copied().collect();
        let pattern: Vec<u8> = part1.iter().chain(&part2).copied().collect();

        assert_eq!(strpcmp(&query, &pattern), false);
    }

    #[test]
    fn substring_without_wildcard_1() {
        let random_input = random_u8_vec(32);

        let pattern = random_input[..31].to_vec();
        assert_eq!(strpcmp(&random_input, &pattern), false);
    }

    #[test]
    fn substring_without_wildcard_2() {
        let random_input = random_u8_vec(32);

        let query = random_input[..31].to_vec();
        assert_eq!(strpcmp(&query, &random_input), false);
    }
}
