#[cfg(test)]
mod shortinette_tests_rust_0205 {
    use super::*;

    #[test]
    fn test_color_new() {
        let color = Color::new(255, 255, 0);
        assert_eq!(color.red, 255);
        assert_eq!(color.green, 255);
        assert_eq!(color.blue, 0);
    }

    #[test]
    fn test_color_constants() {
        assert_eq!(Color::WHITE, Color::new(255, 255, 255));
        assert_eq!(Color::RED, Color::new(255, 0, 0));
        assert_eq!(Color::GREEN, Color::new(0, 255, 0));
        assert_eq!(Color::BLUE, Color::new(0, 0, 255));
    }

    #[test]
    fn test_color_closest_mix_empty_palette() {
        assert_eq!(Color::RED.closest_mix(&[], 100), Color::WHITE);
        assert_eq!(
            Color::RED.closest_mix(
                &[(Color::RED, 255), (Color::GREEN, 255), (Color::BLUE, 255)],
                0
            ),
            Color::WHITE,
        );
    }

    #[test]
    fn test_closest_mix_identical_colors_different_alpha() {
        let palette = &[(Color::RED, 255), (Color::RED, 128)];
        let target = Color::RED;

        assert_eq!(target.closest_mix(palette, 2), Color::RED);
    }

    #[test]
    fn test_closest_color_exact_match() {
        let palatte = &[(Color::RED, 255), (Color::GREEN, 255), (Color::WHITE, 255)];
        let target = Color::GREEN;

        assert_eq!(target.closest_mix(palatte, 5), Color::GREEN);
    }

    #[test]
    fn test_closest_mix_blending() {
        let target = Color::new(128, 128, 128);
        let palette = &[(Color::RED, 128), (Color::GREEN, 128), (Color::BLUE, 128)];
        let expected = Color::new(159, 95, 63);
        let result = target.closest_mix(palette, 3);

        assert_eq!(result.red, expected.red);
        assert_eq!(result.green, expected.green);
        assert_eq!(result.blue, expected.blue);
    }
}
