#[cfg(test)]
mod shortinette_tests {
    use super::*;

    #[test]
    fn subject_test_a() {
        let v = Vector {
            x: String::from("Hello, World!"),
            y: String::from("Hello, Rust!"),
        };

        let w = v.clone();

        assert_eq!(&v, &w);
    }

    #[test]
    fn subject_test_b() {
        let v = Vector::new("Hello, World!", "Hello, Rust!");
        let a = v;
        let b = v;

        assert_eq!(a, b);
    }

    #[test]
    fn new() {
        let v = Vector::new(1, 2);
        let w = Vector { x: 1, y: 2 };

        assert_eq!(
            v, w,
            "Vector::new() and Vector {{}} do not result in the same"
        );
    }

    #[test]
    fn add() {
        let lhs = Vector::new(1, 2);
        let rhs = Vector::new(5, 7);

        let res = lhs + rhs;
        let expected = Vector::new(6, 9);

        assert_eq!(res, expected);
    }

    #[test]
    fn add_assign() {
        let mut lhs = Vector::new(1, 2);
        let rhs = Vector::new(5, 7);

        lhs += rhs;
        let expected = Vector::new(6, 9);

        assert_eq!(lhs, expected);
    }

    #[test]
    fn sub() {
        let lhs = Vector::new(1, 2);
        let rhs = Vector::new(5, 7);

        let res = lhs - rhs;
        let expected = Vector::new(-4, -5);

        assert_eq!(res, expected);
    }

    #[test]
    fn sub_assign() {
        let mut lhs = Vector::new(1, 2);
        let rhs = Vector::new(5, 7);

        lhs -= rhs;
        let expected = Vector::new(-4, -5);

        assert_eq!(lhs, expected);
    }

    #[test]
    fn mul() {
        let lhs = Vector::new(5, 8);
        let rhs = 2;

        let res = lhs * rhs;
        let expected = Vector::new(10, 16);

        assert_eq!(res, expected);
    }

    #[test]
    fn mul_assign() {
        let mut lhs = Vector::new(5, 8);
        let rhs = 2;

        lhs *= rhs;
        let expected = Vector::new(10, 16);

        assert_eq!(lhs, expected);
    }

    #[test]
    fn div() {
        let lhs = Vector::new(5, 8);
        let rhs = 2;

        let res = lhs / rhs;
        let expected = Vector::new(2, 4);

        assert_eq!(res, expected);
    }

    #[test]
    fn div_assign() {
        let mut lhs = Vector::new(5, 8);
        let rhs = 2;

        lhs /= rhs;
        let expected = Vector::new(2, 4);

        assert_eq!(lhs, expected);
    }

    #[test]
    fn eq() {
        let lhs = Vector::new(1, 2);
        let rhs = Vector::new(1, 2);

        assert!(lhs == rhs);

        let rhs = Vector::new(2, 1);
        assert!(!(lhs == rhs));
    }

    #[test]
    fn ne() {
        let lhs = Vector::new(1, 2);
        let rhs = Vector::new(2, 1);

        assert!(lhs != rhs);

        let rhs = Vector::new(1, 2);
        assert!(!(lhs != rhs));
    }

    #[test]
    fn length_f32() {
        let mut vector: Vector<f32>;

        vector = Vector { x: 3.0, y: 4.0 };
        assert_eq!(5.0, vector.length());

        vector = Vector { x: 0.0, y: 0.0 };
        assert_eq!(0.0, vector.length());

        vector = Vector { x: -3.0, y: -4.0 };
        assert_eq!(5.0, vector.length());

        vector = Vector { x: 1e10, y: 1e10 };
        assert!((vector.length() - 1.414_213_5e10).abs() < 1e5);

        vector = Vector { x: 1e-10, y: 1e-10 };
        assert!((vector.length() - 1.414_213_6e-10).abs() < 1e-15);
    }

    #[test]
    fn length_f64() {
        let mut vector: Vector<f64>;

        vector = Vector { x: 3.0, y: 4.0 };
        assert_eq!(5.0, vector.length());

        vector = Vector { x: 0.0, y: 0.0 };
        assert_eq!(0.0, vector.length());

        vector = Vector { x: -3.0, y: -4.0 };
        assert_eq!(5.0, vector.length());

        vector = Vector { x: 1e10, y: 1e10 };
        assert!((vector.length() - 1.414_213_5e10).abs() < 1e5);

        vector = Vector { x: 1e-10, y: 1e-10 };
        assert!((vector.length() - 1.414_213_6e-10).abs() < 1e-15);
    }
}
