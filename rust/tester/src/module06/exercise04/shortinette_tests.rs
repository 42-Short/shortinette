
#[cfg(test)]
mod tests {
    use super::*;
    use libc;
    use rand::{distributions::Alphanumeric, random, thread_rng, Rng};
    use std::ffi::CString;

    fn random_file_name() -> CString {
        let s: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u8>().min(10) as usize)
            .map(char::from)
            .collect();

        CString::new(format!("/tmp/{s}")).expect("CString creation failed.")
    }

    #[test]
    fn test_errno_last() {
        unsafe {
            *libc::__errno_location() = libc::EINTR;
        }
        let errno = Errno::last();
        assert_eq!(errno.0, libc::EINTR);
    }

    #[test]
    fn test_errno_make_last() {
        let errno = Errno(libc::ENOMEM);
        errno.make_last();
        let new_errno = unsafe { *libc::__errno_location() };
        assert_eq!(new_errno, libc::ENOMEM);
    }

    #[test]
    fn test_errno_description() {
        let errno = Errno(libc::EIO);
        let description = errno.description();
        assert_eq!(description, "Input/output error");

        let errno = Errno(libc::ENOSPC);
        let description = errno.description();
        assert_eq!(description, "No space left on device");

        let errno = Errno(libc::EPERM);
        let description = errno.description();
        assert_eq!(description, "Operation not permitted");
    }

    #[test]
    fn test_fd_open_success() {
        let file_name = random_file_name();
        let fd = Fd::create(&file_name).expect("Call to Fd::create() failed.");
        assert!(fd.0 >= 0);
        fd.close().expect("Call to Fd::close() failed.");

        let fd_opened = Fd::open(&file_name).expect("Call to Fd::open() failed.");
        assert!(fd_opened.0 >= 0);
        fd_opened.close().expect("Call to Fd::close() failed.");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_open_fail() {
        let file_name = random_file_name();
        let result = Fd::open(&file_name);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
            assert_eq!(e.description(), "No such file or directory");
        }
    }

    #[test]
    fn test_fd_create_success() {
        let file_name = random_file_name();
        let fd = Fd::create(&file_name).expect("Failed to create file");
        assert!(fd.0 >= 0);
        fd.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_write_success() {
        let file_name = random_file_name();
        let fd = Fd::create(&file_name).expect("Failed to create file");
        let data = b"Test data";
        let bytes_written = fd.write(data).expect("Failed to write to file");
        assert_eq!(bytes_written, data.len());

        fd.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_write_fail() {
        let fd = Fd(-1);
        let data = b"Test data";
        let result = fd.write(data);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
        }
    }

    #[test]
    fn test_fd_read_success() {
        let file_name = random_file_name();
        let write_fd = Fd::create(&file_name).expect("Failed to create file");
        let data = b"Test data";
        write_fd.write(data).expect("Failed to write to file");

        let read_fd = Fd::open(&file_name).unwrap();
        let mut buffer = [0u8; 9];
        let bytes_read = read_fd.read(&mut buffer).expect("Failed to read from file");
        assert_eq!(bytes_read, data.len());
        assert_eq!(&buffer, data);

        read_fd.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_read_fail() {
        let fd = Fd(-1);
        let mut buffer = [0u8; 9];
        let result = fd.read(&mut buffer);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
        }
    }

    #[test]
    fn test_fd_close_success() {
        let file_name = random_file_name();
        let fd = Fd::create(&file_name).expect("Failed to create file");
        assert!(fd.close().is_ok());

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_open_success() {
        let file_name = random_file_name();
        let _ = Fd::create(&file_name).expect("Failed to create file");

        let file = File::open(&file_name).expect("Failed to open file");
        file.0.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_open_fail() {
        let file_name = random_file_name();
        let result = File::open(&file_name);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
            assert_eq!(e.description(), "No such file or directory");
        }
    }

    #[test]
    fn test_file_create_success() {
        let file_name = random_file_name();
        let file = File::create(&file_name).expect("Failed to create file");
        assert!(file.0 .0 >= 0);

        file.0.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_write_and_read_success() {
        let file_name = random_file_name();
        let write_file = File::create(&file_name).expect("Failed to create file for writing");
        let data = b"Hello, File!";

        {
            let bytes_written = write_file.write(data).expect("Failed to write to file");
            assert_eq!(bytes_written, data.len());
        }
        std::thread::sleep(std::time::Duration::from_millis(10));
        let read_file = File::open(&file_name).expect("Failed to open file for reading");
        let mut buffer = [0u8; 12];
        read_file
            .read(&mut buffer)
            .expect("Failed to read from file");
        assert_eq!(&buffer, data);

        read_file.0.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_leak() {
        let file_name = random_file_name();
        let file = File::create(&file_name).expect("Failed to create file");
        let fd = file.leak();
        assert!(fd.0 >= 0);

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_drop() {
        let file_name = random_file_name();
        {
            let file = File::create(&file_name).expect("Failed to create file");
            let _ = file.write(b"Testing drop").expect("Failed to drop file");
        }

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }
}
