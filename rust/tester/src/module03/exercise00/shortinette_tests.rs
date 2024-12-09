#[cfg(test)]
mod shortinette_tests {
    use super::*;

    #[test]
    #[should_panic]
    fn empty() {
        let empty: Vec<i32> = Vec::new();

        let value = choose(&empty);
        println!(
            "What in the world? This should not be printed! Value: {:?}",
            value
        );
    }

    #[test]
    fn single() {
        let alone = &[0];
        let value = choose(alone);

        assert_eq!(value, &0, "How is choose(&[0]) not returning 0?");
    }

    #[test]
    fn generic() {
        let numbers = [1_u8, 2, 3];
        let _: &u8 = choose(&numbers);

        let slices = ["a", "b", "c"];
        let _: &&str = choose(&slices);

        let bools = [true, false];
        let _: &bool = choose(&bools);

        struct Foo;
        let foos = [Foo, Foo];
        let _: &Foo = choose(&foos);
    }

    #[test]
    fn randomness() {
        let huge: Vec<_> = (0..100_000).collect();

        // Well this is one of the cases where a second grademe could actually
        // pass the exercise.
        // Although it is very unlikely that this will ever happen.
        let value = choose(&huge);
        let value2 = choose(&huge);

        assert_ne!(
            value, value2,
            "choose(&huge) returned the same value twice. Do you really return a random element?"
        );
    }
}
