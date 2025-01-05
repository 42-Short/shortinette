
#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn swap_basic_type() {
        let mut a = 69;
        let mut b = 420;
        ft_swap(&mut a, &mut b);
        assert_eq!(a, 420);
        assert_eq!(b, 69);
    }

    #[test]
    fn swap_complex_type() {
        let mut a = String::from("Hello, World!");
        let mut b = String::from("Goodbye, World!");
        ft_swap(&mut a, &mut b);
        assert_eq!(a, "Goodbye, World!");
        assert_eq!(b, "Hello, World!");
    }

    #[test]
    fn strlen() {
        let s = b"Hello, World!\0";
        let len = unsafe { ft_strlen(s.as_ptr()) };
        assert_eq!(len, 13);
    }

    #[test]
    fn strcpy() {
        let s = b"Hello, World!\0";
        let mut dst = [0u8; 14];
        unsafe { ft_strcpy(dst.as_mut_ptr(), s.as_ptr()) };
        assert_eq!(&dst, s);
    }
}
