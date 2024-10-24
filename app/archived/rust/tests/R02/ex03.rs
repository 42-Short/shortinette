#[cfg(test)]
mod shortinette_tests_rust_0203 {
    use super::*;

    #[test]
    fn test_clone_trait() {
        let instance = MyType::default();
        assert_eq!(instance, instance.clone());
    }

    #[test]
    fn test_partial_eq_trait() {
        assert_eq!(MyType::default(), MyType::default());
    }

    #[test]
    fn test_partial_ord_trait() {
        let instance1 = MyType::default();
        let instance2 = MyType::default();
        assert!(instance1 <= instance2 && instance1 >= instance2);
    }

    #[test]
    fn test_debug_trait() {
        assert_eq!(format!("{:?}", MyType::default()), "MyType");
    }
}
