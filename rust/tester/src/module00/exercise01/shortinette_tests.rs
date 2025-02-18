mod min;

#[cfg(test)]
mod shortinette_tests_0001 {
    use crate::min::min;
    use rand::{random, Rng};

    #[test]
    fn test_equal() {
        let number = random::<i32>();

        assert_eq!(
            min(number, number),
            number,
            "Failed with ({}, {})",
            number,
            number
        );
    }

    #[test]
    fn test_a_lower() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(i32::MIN..i32::MAX);
        let b = rng.gen_range(a + 1..=i32::MAX);

        assert_eq!(min(a, b), a, "Failed with ({}, {})", a, b);
    }

    #[test]
    fn test_b_lower() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(i32::MIN + 1..=i32::MAX);
        let b = rng.gen_range(i32::MIN..a);

        assert_eq!(min(a, b), b, "Failed with ({}, {})", a, b);
    }

    #[test]
    fn test_negative_and_zero() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(i32::MIN..0);
        let b = 0;

        assert_eq!(min(a, b), a, "Failed with ({}, {})", a, b);
    }
}
