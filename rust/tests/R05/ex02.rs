#[cfg(test)]
mod shortinette_rust_test_module05_ex02_0001 {
    use std::thread;

    use super::*;

    #[test]
    fn last() {
        match Error::last() {
            Error::Success => (),
            _ => panic!("If Error::make_last() was never called, Error::last() call should return Error::Success")
        }
    }

    #[test]
    fn make_last() {
        {
            let err = Error::FileNotFound;
            err.make_last();
            match Error::last() {
                Error::FileNotFound => (),
                _ => panic!("Error::make_last() did not set the value correctly"),
            }
        }

        {
            let err = Error::IsDirectory;
            err.make_last();
            match Error::last() {
                Error::IsDirectory => (),
                _ => panic!("Error::make_last() did not set the value correctly"),
            }
        }

        {
            let err = Error::WriteFail;
            err.make_last();
            match Error::last() {
                Error::WriteFail => (),
                _ => panic!("Error::make_last() did not set the value correctly"),
            }
        }

        {
            let err = Error::ReadFail;
            err.make_last();
            match Error::last() {
                Error::ReadFail => (),
                _ => panic!("Error::make_last() did not set the value correctly"),
            }
        }
    }

    #[test]
    fn multiple_threads() {
        let err = Error::FileNotFound;
        err.make_last();

        let handle = thread::spawn(|| {
            match Error::last() {
                Error::Success => (),
                _ => panic!("Error::make_last() should not affect values across threads."),
            }

            let err = Error::ReadFail;
            err.make_last();
            match Error::last() {
                Error::ReadFail => (),
                _ => panic!("Error::make_last() did not set the value correctly"),
            }
        });

        let handle2 = thread::spawn(|| {
            match Error::last() {
                Error::Success => (),
                _ => panic!("Error::make_last() should not affect values across threads."),
            }

            let err = Error::IsDirectory;
            err.make_last();
            match Error::last() {
                Error::IsDirectory => (),
                _ => panic!("Error::make_last() did not set the value correctly"),
            }
        });

        handle.join().unwrap();
        handle2.join().unwrap();

        match Error::last() {
            Error::FileNotFound => (),
            _ => panic!("Error::make_last() did not set the value correctly"),
        }
    }
}
