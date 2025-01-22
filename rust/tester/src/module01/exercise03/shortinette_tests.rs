#[cfg(test)]
mod shortinette_tests_rust_0103 {
    use ex03::first_group;
    use rand::Rng;

    #[test]
    fn test_lifetimes() {
        let haystack = [1, 2, 3, 2, 1];
        let result;

        {
            let needle = [2, 3];
            result = first_group(&haystack, &needle);
        }

        assert_eq!(result, &[2, 3]);
    }

    #[test]
    fn test_both_empty() {
        assert_eq!(
            first_group(&[], &[]),
            &[],
            "Empty slice expected when both haystack and needle are empty."
        );
    }

    #[test]
    fn test_needle_empty() {
        let mut rng = rand::thread_rng();

        let haystack: [u32; 1] = [rng.gen_range(0..=u32::MAX)];
        assert_eq!(
            first_group(&haystack, &[]),
            &[],
            "Empty slice expected when needle is empty."
        );
    }

    #[test]
    fn test_haystack_empty() {
        let mut rng = rand::thread_rng();

        let needle: [u32; 1] = [rng.gen_range(0..=u32::MAX)];
        assert_eq!(
            first_group(&[], &needle),
            &[],
            "Empty slice expected when haystack is empty."
        );
    }

    #[test]
    fn test_match_at_start() {
        let mut rng = rand::thread_rng();

        let needle_length = rng.gen_range(5..15);
        let needle: Vec<u32> = (0..needle_length)
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();
        let mut haystack = needle.clone();

        for _ in 0..rng.gen_range(10..100) {
            haystack.push(rng.gen_range(0..=u32::MAX));
        }

        assert_eq!(
            first_group(&haystack, &needle),
            &haystack[0..needle_length],
            "Invalid result with match at the start of the haystack."
        );
    }

    #[test]
    fn test_match_at_end() {
        let mut rng = rand::thread_rng();

        let needle_length = rng.gen_range(5..15);
        let needle: Vec<u32> = (0..needle_length)
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let mut haystack: Vec<u32> = (0..rng.gen_range(10..100))
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();
        haystack.extend(&needle);

        assert_eq!(
            first_group(&haystack, &needle),
            &haystack[haystack.len() - needle_length..],
            "Invalid result with match at the end of the haystack."
        );
    }

    #[test]
    fn test_match_in_middle() {
        let mut rng = rand::thread_rng();

        let needle_length = rng.gen_range(5..15);
        let needle: Vec<u32> = (0..needle_length)
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let mut haystack: Vec<u32> = (0..rng.gen_range(5..50))
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();
        let startpos = haystack.len();
        haystack.extend(&needle);

        for _ in 0..rng.gen_range(5..50) {
            haystack.push(rng.gen_range(0..=u32::MAX));
        }

        assert_eq!(
            first_group(&haystack, &needle),
            &haystack[startpos..startpos + needle_length],
            "Invalid result with match in the middle of the haystack."
        );
    }

    #[test]
    fn test_match_short_needle() {
        let mut rng = rand::thread_rng();

        let haystack: Vec<u32> = (0..rng.gen_range(10..100))
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let mut needle_pos;
        loop {
            needle_pos = rng.gen_range(0..haystack.len());
            let needle = haystack[needle_pos];

            // Ensure that the item we are searching for doesn't exist multiple times
            // There is already another test that checks for multiple matches
            if haystack.iter().filter(|&&item| item == needle).count() == 1 {
                break;
            }
        }
        let needle = [haystack[needle_pos]];

        assert_eq!(
            first_group(&haystack, &needle),
            &[haystack[needle_pos]],
            "Invalid result if the needle consists of only one element."
        );
    }

    #[test]
    fn test_multiple_matches() {
        let mut rng = rand::thread_rng();

        let needle_length = rng.gen_range(5..15);
        let needle: Vec<u32> = (0..needle_length)
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let mut haystack: Vec<u32> = (0..rng.gen_range(5..50))
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let startpos = haystack.len();
        haystack.extend(&needle);

        for _ in 0..rng.gen_range(5..50) {
            haystack.push(rng.gen_range(0..=u32::MAX));
        }

        haystack.extend(&needle);
        for _ in 0..rng.gen_range(5..50) {
            haystack.push(rng.gen_range(0..=u32::MAX));
        }

        assert!(
            std::ptr::eq(
                first_group(&haystack, &needle),
                &haystack[startpos..startpos + needle_length]
            ),
            "Function didn't return the first match in haystack when there were multiple matches"
        );
    }

    #[test]
    fn test_needle_longer_than_haystack() {
        let mut rng = rand::thread_rng();

        let haystack: Vec<u32> = (0..rng.gen_range(10..100))
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let needle: Vec<u32> = haystack
            .iter()
            .copied()
            .chain(std::iter::once(rng.gen_range(0..=u32::MAX)))
            .collect();

        assert_eq!(
            first_group(&haystack, &needle),
            &[],
            "Expected empty slice when needle is longer than haystack."
        );
    }

    #[test]
    fn test_no_match() {
        let mut rng = rand::thread_rng();

        let haystack: Vec<u32> = (0..rng.gen_range(10..100))
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let mut needle_val;
        loop {
            needle_val = rng.gen_range(0..=u32::MAX);

            // Ensure item doesn't exist
            if haystack.iter().all(|&item| item != needle_val) {
                break;
            }
        }
        let needle = [needle_val];

        assert_eq!(
            first_group(&haystack, &needle),
            &[],
            "Expected empty slice when item from needle is not present in haystack."
        );
    }

    #[test]
    fn test_needle_wrong_order() {
        let mut rng = rand::thread_rng();

        let haystack: Vec<u32> = (0..rng.gen_range(10..100))
            .map(|_| rng.gen_range(0..=u32::MAX))
            .collect();

        let mut needle: Vec<u32>;
        // What are the chances that the haystack is the same if you reverse it
        loop {
            needle = haystack.iter().rev().copied().collect();
            if haystack != needle {
                break;
            }
        }

        assert_eq!(
            first_group(&haystack, &needle),
            &[],
            "Expected empty slice when all items from the needle are present in haystack, but the order is wrong."
        );
    }
}
