#[cfg(test)]
mod shortinette_rust_test_module05_ex00_0001 {
    use super::*;

    #[test]
    fn u32() {
        let a = Cell::new(1);
        let b = Cell::new(3);

        swap_u32(&a, &b);

        assert_eq!(a.get(), 3);
        assert_eq!(b.get(), 1);
    }

    #[test]
    fn string() {
        let a = Cell::new("ABC".into());
        let b = Cell::new("DEF".into());

        swap_string(&a, &b);

        assert_eq!(a.into_inner(), "DEF");
        assert_eq!(b.into_inner(), "ABC");
    }
}