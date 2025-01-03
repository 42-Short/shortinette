#[cfg(test)]
mod test {
    use super::*;

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
            Outcome::Bad(_) => panic!("should be Good"),
        }

        let m = maybe();
        match m {
            Maybe::Definitely(n) => assert_eq!(n, 42),
            Maybe::No => panic!("should be Definitely"),
        }
    }
}