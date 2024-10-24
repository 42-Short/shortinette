#[cfg(test)]
mod shortinette_rust_test_module05_ex07_0001 {
    use std::{sync::Arc, thread, time};

    use super::*;

    #[test]
    fn basic() {
        let rdv = Arc::new(RendezVous::new());

        let rdv2 = rdv.clone();
        let handle = thread::spawn(move || rdv2.wait(42));

        let a = rdv.wait(21);
        let b = handle.join().unwrap();

        assert_eq!(42, a);
        assert_eq!(21, b);
    }

    #[test]
    fn generic() {
        {
            let rdv = Arc::new(RendezVous::new());

            let rdv2 = rdv.clone();
            let handle = thread::spawn(move || rdv2.wait("hello"));

            let a = rdv.wait("world");
            let b = handle.join().unwrap();

            assert_eq!("hello", a);
            assert_eq!("world", b);
        }

        {
            let rdv = Arc::new(RendezVous::new());

            let rdv2 = rdv.clone();
            let handle = thread::spawn(move || rdv2.wait(true));

            let a = rdv.wait(false);
            let b = handle.join().unwrap();

            assert!(a);
            assert!(!b);
        }
    }

    #[test]
    fn reuse_same_threads() {
        let rdv = Arc::new(RendezVous::new());

        let pairs = vec![
            (21, 42),
            (700, 890),
            (222, 111),
            (4829, 1221),
            (7273, 1284),
            (1331, 4444),
            (27, 91),
            (3417, 286),
            (58, 4923),
            (1298, 75),
            (6721, 18),
            (44, 321),
            (97, 4852),
            (739, 62),
            (8234, 57),
            (19, 481),
            (613, 7342),
            (3021, 105),
            (491, 205),
            (8762, 37),
            (55, 7219),
            (308, 4095),
            (951, 72),
            (6271, 28),
            (405, 738),
            (230, 5921),
            (18, 6438),
            (4123, 95),
            (71, 503),
            (5238, 27),
        ];

        let apairs = pairs.clone();
        let rdv2 = rdv.clone();
        let a = thread::spawn(move || {
            for pair in apairs {
                let res = rdv2.wait(pair.0);
                assert_eq!(pair.1, res);
            }
        });

        let rdv2 = rdv.clone();
        let b = thread::spawn(move || {
            for pair in pairs {
                let res = rdv2.wait(pair.1);
                assert_eq!(pair.0, res);
            }
        });

        a.join().unwrap();
        b.join().unwrap();
    }

    #[test]
    fn reuse_different_threads() {
        let rdv = Arc::new(RendezVous::new());

        let rdv2 = rdv.clone();
        let a = thread::spawn(move || rdv2.wait(21));
        let rdv2 = rdv.clone();
        let b = thread::spawn(move || rdv2.wait(42));

        let a = a.join().unwrap();
        let b = b.join().unwrap();

        assert_eq!(42, a);
        assert_eq!(21, b);

        let rdv2 = rdv.clone();
        let handle = thread::spawn(move || rdv2.wait(100));

        let a = rdv.wait(200);
        let b = handle.join().unwrap();

        assert_eq!(100, a);
        assert_eq!(200, b);
    }

    #[test]
    fn large_delay() {
        let rdv = Arc::new(RendezVous::new());
        let rdv2 = rdv.clone();

        let a = thread::spawn(move || {
            thread::sleep(time::Duration::from_millis(50));

            rdv.wait(21)
        });

        let b = thread::spawn(move || {
            thread::sleep(time::Duration::from_millis(10));

            rdv2.wait(42)
        });

        let a = a.join().unwrap();
        let b = b.join().unwrap();

        assert_eq!(42, a);
        assert_eq!(21, b);
    }

    #[test]
    fn try_wait() {
        let rdv = Arc::new(RendezVous::new());

        let res = rdv.try_wait(42);
        assert_eq!(Err(42), res);

        let rdv2 = rdv.clone();
        let a = thread::spawn(move || rdv2.try_wait(21)).join().unwrap();
        assert_eq!(Err(21), a);

        let rdv2 = rdv.clone();
        let a = thread::spawn(move || rdv2.wait(50));
        let rdv2 = rdv.clone();
        let b = thread::spawn(move || rdv2.try_wait(80));

        let a = a.join().unwrap();
        let b = b.join().unwrap();

        assert_eq!(80, a);
        assert_eq!(Ok(50), b);
    }
}
