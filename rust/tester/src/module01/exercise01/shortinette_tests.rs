#[cfg(test)]
mod shortinette_tests_0101 {
    use ex01::min;
    use rand::Rng;

    #[test]
    fn test_1() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(1..=100);
        let b = a + 1;

        assert_eq!(min(&a, &b), &a);
    }

    #[test]
    fn test_2() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(1..=100);
        let b = a - 1;

        assert_eq!(min(&a, &b), &b);
    }

    #[test]
    fn test_equal() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(1..=100);
        let b = a;

        assert!(std::ptr::eq(min(&a, &b), &b));
    }

    #[test]
    fn test_negative_numbers() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(i32::MIN..0);
        let b = rng.gen_range(i32::MIN..0);

        let expected = if b <= a { &b } else { &a };

        assert_eq!(min(&a, &b), expected);
    }

    #[test]
    fn test_i32_min() {
        let mut rng = rand::thread_rng();

        let mut a = i32::MIN;
        let mut b = rng.gen_range(i32::MIN + 1..=i32::MAX);

        if rng.gen_range(0..=1) == 0 {
            std::mem::swap(&mut a, &mut b);
        }

        let expected = if b <= a { &b } else { &a };

        assert_eq!(min(&a, &b), &a);
    }

    #[test]
    fn test_i32_max() {
        let mut rng = rand::thread_rng();

        let mut a = rng.gen_range(i32::MIN..i32::MAX);
        let mut b = i32::MAX;

        if rng.gen_range(0..=1) == 0 {
            std::mem::swap(&mut a, &mut b);
        }

        assert_eq!(min(&a, &b), &a);
    }
}
