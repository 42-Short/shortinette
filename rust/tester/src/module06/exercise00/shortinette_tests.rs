#[cfg(test)]
mod tests {
    use super::*;
    use rand::{distributions::Alphanumeric, random, thread_rng, Rng};

    #[test]
    fn swap_basic_type() {
        let mut a = random::<i32>();
        let mut b = random::<i32>();

        let a0 = a;
        let b0 = b;

        ft_swap(&mut a, &mut b);
        assert_eq!(a, b0);
        assert_eq!(b, a0);
    }

    #[test]
    fn swap_string() {
        let mut a: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u16>() as usize)
            .map(char::from)
            .collect();
        let mut b: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u16>() as usize)
            .map(char::from)
            .collect();

        let a0 = a.clone();
        let b0 = b.clone();

        ft_swap(&mut a, &mut b);
        assert_eq!(a, b0);
        assert_eq!(b, a0);
    }

    #[test]
    fn strlen_basic() {
        let randstring: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u16>() as usize)
            .map(char::from)
            .collect();

        let randstring_null_terminated = format!("{}\0", randstring);
        let randbytes: &[u8] = randstring_null_terminated.as_bytes();

        let len = unsafe { ft_strlen(randbytes.as_ptr()) };
        assert_eq!(len, randstring.len());
    }

    #[test]
    fn strlen_empty_string() {
        let empty_string = b"\0";
        assert_eq!(unsafe { ft_strlen(empty_string.as_ptr()) }, 0)
    }

    #[test]
    fn strcpy() {
        let randstring: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u16>() as usize)
            .map(char::from)
            .collect();

        let randstring_null_terminated = format!("{}\0", randstring);
        let randbytes: &[u8] = randstring_null_terminated.as_bytes();

        let mut dst = [0u8; u16::MAX as usize];
        unsafe { ft_strcpy(dst.as_mut_ptr(), randbytes.as_ptr()) };
        assert_eq!(&dst[0..randbytes.len()], randbytes);
    }
}
