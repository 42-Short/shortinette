#[cfg(test)]
mod shortinette_rust_test_module05_ex04_0001 {
    use std::{collections::HashSet, thread};

    use super::*;

    #[test]
    fn subject() {
        let a = Unique::new();
        let b = Unique::new();
        let c = Unique::new();

        assert_eq!(a.0, 0);
        assert_eq!(b.0, 1);
        assert_eq!(c.0, 2);

        let d = a.clone();
        let e = c.clone();

        assert_eq!(d.0, 3);
        assert_eq!(e.0, 4);
    }

    #[test]
    #[should_panic]
    fn too_many_ids() {
        for id in 0..=255 {
            let unique = Unique::new();
            assert_eq!(unique.0 as usize, id);
        }
    }

    #[test]
    fn threads() {
        let mut unique_set: HashSet<_> = HashSet::new();

        let a = Unique::new();
        assert_eq!(a.0, 0);
        unique_set.insert(a.0);

        // 253 because 255 should never exist and
        // at the end another unique gets created
        let threads = (1..=253).map(|_| {
            thread::spawn(move || {
                let unique = Unique::new();
                unique.0
            })
        });

        threads.map(|thread| thread.join().unwrap()).for_each(|id| {
            unique_set.insert(id);
        });

        let unique = Unique::new();
        unique_set.insert(unique.0);

        if unique_set.contains(&255) {
            panic!("Unique(255) should never exist if AtomicU8 was used correctly");
        }

        assert_eq!(unique_set.len(), 255);
    }
}
