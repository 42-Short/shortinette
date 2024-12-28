#[cfg(test)]
mod shortinette_tests {
    use std::cell::Cell;

    use super::*;

    #[test]
    fn new_type() {
        #[derive(PartialEq, Eq, Debug)]
        struct NewType;

        impl FortyTwo for NewType {
            fn forty_two() -> Self {
                Self
            }
        }

        assert_eq!(
            <NewType as FortyTwo>::forty_two(),
            NewType,
            "FortyTwo::forty_two() did not return correct result for a custom type"
        );
    }

    #[test]
    fn obvious() {
        assert_eq!(u32::forty_two(), 42);

        String::forty_two();
    }

    // Probably the most hacky test ever.
    // It does not check whether the println! output is correct.
    // But come on, if FortyTwo::forty_two() was called AND the debug formatting
    // was called then it makes no sense to not call println! lol.
    #[test]
    fn forty_two() {
        thread_local! {
            static FT_CALLED: Cell<bool> = const { Cell::new(false) };
            static DEBUG_CALLED: Cell<bool> = const { Cell::new(false) };
        }

        #[derive(PartialEq, Eq)]
        struct NewType;

        impl FortyTwo for NewType {
            fn forty_two() -> Self {
                FT_CALLED.set(true);

                Self
            }
        }

        impl Debug for NewType {
            fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                DEBUG_CALLED.set(true);

                f.write_str("42")
            }
        }

        print_forty_two::<NewType>();

        assert!(FT_CALLED.get());
        assert!(DEBUG_CALLED.get());
    }
}
