#[cfg(test)]
mod shortinette_tests_rust_0102 {
    use ex02::color_name;
    use rand::Rng;

    #[test]
    fn test_lifetimes() {
        let name_of_the_best_color;

        {
            let the_best_color = [42, 42, 42];
            name_of_the_best_color = color_name(&the_best_color);
        }

        assert_eq!(name_of_the_best_color, "dark gray");
    }

    #[test]
    fn test_dark_gray_1() {
        let color = color_name(&[127, 127, 127]);
        assert_eq!(color, "dark gray");
    }

    #[test]
    fn test_dark_gray_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(0..128),
            rng.gen_range(0..128),
            rng.gen_range(0..128),
        ]);
        assert_eq!(color, "dark gray");
    }

    #[test]
    fn test_dark_red_1() {
        let color = color_name(&[128, 127, 127]);
        assert_eq!(color, "dark red");
    }

    #[test]
    fn test_dark_red_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
            rng.gen_range(0..128),
        ]);
        assert_eq!(color, "dark red");
    }

    #[test]
    fn test_dark_green_1() {
        let color = color_name(&[127, 128, 127]);
        assert_eq!(color, "dark green");
    }

    #[test]
    fn test_dark_green_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
        ]);
        assert_eq!(color, "dark green");
    }

    #[test]
    fn test_olive_1() {
        let color = color_name(&[128, 128, 127]);
        assert_eq!(color, "olive");
    }

    #[test]
    fn test_olive_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
        ]);
        assert_eq!(color, "olive");
    }

    #[test]
    fn test_dark_blue_1() {
        let color = color_name(&[127, 127, 128]);
        assert_eq!(color, "dark blue");
    }

    #[test]
    fn test_dark_blue_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(0..128),
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
        ]);
        assert_eq!(color, "dark blue");
    }

    #[test]
    fn test_purple_1() {
        let color = color_name(&[128, 127, 128]);
        assert_eq!(color, "purple");
    }

    #[test]
    fn test_purple_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
        ]);
        assert_eq!(color, "purple");
    }

    #[test]
    fn test_teal_1() {
        let color = color_name(&[127, 128, 128]);
        assert_eq!(color, "teal");
    }

    #[test]
    fn test_teal_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
        ]);
        assert_eq!(color, "teal");
    }

    #[test]
    fn test_light_gray_1() {
        let color = color_name(&[128, 128, 128]);
        assert_eq!(color, "light gray");
    }

    #[test]
    fn test_light_gray_2() {
        let mut rng = rand::thread_rng();

        let color = color_name(&[
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
        ]);
        assert_eq!(color, "light gray");
    }
}
