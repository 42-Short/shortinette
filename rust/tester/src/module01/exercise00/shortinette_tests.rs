#[cfg(test)]
mod shortinette_tests_0100 {
    use ex00::{add, add_assign};
    use rand::Rng;

    #[test]
    fn test_add_1() {
        let mut rng = rand::thread_rng();

        let a: i32 = rng.gen_range(1..=100);
        let b: i32 = rng.gen_range(1..=100);

        assert_eq!(add(&a, b), a + b, "Failed for ({}, {})", a, b);
    }

    #[test]
    fn test_add_2() {
        let a = i32::MAX;

        assert_eq!(add(&a, 0), i32::MAX, "Failed for ({}, {})", a, 0);
    }

    #[test]
    fn test_add_3() {
        let mut rng = rand::thread_rng();

        let a = i32::MAX;
        let b = rng.gen_range(i32::MIN..0);

        assert_eq!(add(&a, b), a + b, "Failed for ({}, {})", a, b);
    }

    #[test]
    fn test_add_4() {
        let mut rng = rand::thread_rng();

        let a = rng.gen_range(i32::MIN..0);
        let b = i32::MAX;

        assert_eq!(add(&a, b), a + b, "Failed for ({}, {})", a, b);
    }

    #[test]
    fn test_add_assign_1() {
        let mut rng = rand::thread_rng();

        let mut a: i32 = rng.gen_range(1..=100);
        let b: i32 = rng.gen_range(1..=100);

        let expected = a + b;
        add_assign(&mut a, b);

        assert_eq!(a, expected, "Failed for ({}, {})", a, b);
    }

    #[test]
    fn test_add_assign_2() {
        let mut a = i32::MAX;

        add_assign(&mut a, 0);
        assert_eq!(a, i32::MAX, "Failed for ({}, {})", a, 0);
    }

    #[test]
    fn test_add_assign_3() {
        let mut rng = rand::thread_rng();

        let mut a = i32::MAX;
        let b = rng.gen_range(i32::MIN..0);

        let expected = a + b;
        add_assign(&mut a, b);

        assert_eq!(a, expected, "Failed for ({}, {})", a, b);
    }

    #[test]
    fn test_add_assign_4() {
        let mut rng = rand::thread_rng();

        let mut a = rng.gen_range(i32::MIN..0);
        let b = i32::MAX;

        let expected = a + b;
        add_assign(&mut a, b);

        assert_eq!(a, expected, "Failed for ({}, {})", a, b);
    }
}
