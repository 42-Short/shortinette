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

        assert_eq!(
            name_of_the_best_color, "dark gray",
            "Failed for [42, 42, 42]"
        );
    }

    #[test]
    fn test_dark_gray_1() {
        let input = [127, 127, 127];
        let color = color_name(&input);
        assert_eq!(color, "dark gray", "Failed for {:?}", input);
    }

    #[test]
    fn test_dark_gray_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(0..128),
            rng.gen_range(0..128),
            rng.gen_range(0..128),
        ];
        let color = color_name(&input);
        assert_eq!(color, "dark gray", "Failed for {:?}", input);
    }

    #[test]
    fn test_dark_red_1() {
        let input = [128, 127, 127];
        let color = color_name(&input);
        assert_eq!(color, "dark red", "Failed for {:?}", input);
    }

    #[test]
    fn test_dark_red_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
            rng.gen_range(0..128),
        ];
        let color = color_name(&input);
        assert_eq!(color, "dark red", "Failed for {:?}", input);
    }

    #[test]
    fn test_dark_green_1() {
        let input = [127, 128, 127];
        let color = color_name(&input);
        assert_eq!(color, "dark green", "Failed for {:?}", input);
    }

    #[test]
    fn test_dark_green_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
        ];
        let color = color_name(&input);
        assert_eq!(color, "dark green", "Failed for {:?}", input);
    }

    #[test]
    fn test_olive_1() {
        let input = [128, 128, 127];
        let color = color_name(&input);
        assert_eq!(color, "olive", "Failed for {:?}", input);
    }

    #[test]
    fn test_olive_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
        ];
        let color = color_name(&input);
        assert_eq!(color, "olive", "Failed for {:?}", input);
    }

    #[test]
    fn test_dark_blue_1() {
        let input = [127, 127, 128];
        let color = color_name(&input);
        assert_eq!(color, "dark blue", "Failed for {:?}", input);
    }

    #[test]
    fn test_dark_blue_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(0..128),
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
        ];
        let color = color_name(&input);
        assert_eq!(color, "dark blue", "Failed for {:?}", input);
    }

    #[test]
    fn test_purple_1() {
        let input = [128, 127, 128];
        let color = color_name(&input);
        assert_eq!(color, "purple", "Failed for {:?}", input);
    }

    #[test]
    fn test_purple_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(128..=255),
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
        ];
        let color = color_name(&input);
        assert_eq!(color, "purple", "Failed for {:?}", input);
    }

    #[test]
    fn test_teal_1() {
        let input = [127, 128, 128];
        let color = color_name(&input);
        assert_eq!(color, "teal", "Failed for {:?}", input);
    }

    #[test]
    fn test_teal_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(0..128),
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
        ];
        let color = color_name(&input);
        assert_eq!(color, "teal", "Failed for {:?}", input);
    }

    #[test]
    fn test_light_gray_1() {
        let input = [128, 128, 128];
        let color = color_name(&input);
        assert_eq!(color, "light gray", "Failed for {:?}", input);
    }

    #[test]
    fn test_light_gray_2() {
        let mut rng = rand::thread_rng();

        let input = [
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
            rng.gen_range(128..=255),
        ];
        let color = color_name(&input);
        assert_eq!(color, "light gray", "Failed for {:?}", input);
    }
}
