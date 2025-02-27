#[cfg(test)]
mod shortinette_tests{
    use ex00::*;

    fn outcome() -> Outcome<u32, &'static str> {
        Outcome::Good(42)
    }

    fn maybe() -> Maybe<u8> {
        Maybe::Definitely(42)
    }

    #[test]
    fn test() {
        let o = outcome();
        match o {
            Outcome::Good(n) => assert_eq!(n, 42),
            Outcome::Bad(_) => panic!("Expected: Outcome::Good, got: Outcome::Bad."),
        }

        let m = maybe();
        match m {
            Maybe::Definitely(n) => assert_eq!(n, 42),
            Maybe::No => panic!("Expected: Maybe::Definitely, got: Maybe::No."),
        }
    }
}