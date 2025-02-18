#[cfg(test)]
mod shortinette_tests_0005 {
    use chrono::{Datelike, NaiveDate, Weekday};
    use ex05::{friday_the_13th, is_leap_year, num_days_in_month};
    use rand::Rng;
    use similar_asserts;

    fn expected_output(year: u32) -> String {
        (1..=year)
            .flat_map(move |year| {
                (1..=12).filter_map(move |month| {
                    let date =
                        NaiveDate::from_ymd_opt(year.try_into().unwrap(), month, 13).unwrap();
                    if date.weekday() != Weekday::Fri {
                        None
                    } else {
                        Some(format!("Friday, {} 13, {}\n", date.format("%B"), year))
                    }
                })
            })
            .collect()
    }

    #[test]
    #[should_panic]
    fn test_invalid_year() {
        is_leap_year(0);
    }

    // This one is just to make sure it doesn't panic with year 1
    #[test]
    fn test_year_1() {
        let result = is_leap_year(1);

        assert_eq!(result, false, "Year 1 is not a leap year");
    }

    #[test]
    #[should_panic]
    fn test_invalid_month_0() {
        let mut rng = rand::thread_rng();

        num_days_in_month(rng.gen_range(1..2025), 0);
    }

    #[test]
    #[should_panic]
    fn test_invalid_month_1() {
        let mut rng = rand::thread_rng();

        num_days_in_month(rng.gen_range(1..2025), 13);
    }

    #[test]
    #[should_panic]
    fn test_invalid_month_2() {
        let mut rng = rand::thread_rng();

        num_days_in_month(rng.gen_range(1..2025), rng.gen_range(14..100));
    }

    #[test]
    #[should_panic]
    fn test_num_days_in_month_year_0() {
        let mut rng = rand::thread_rng();

        num_days_in_month(0, rng.gen_range(1..=12));
    }

    #[test]
    fn test_subject_1600() {
        assert_eq!(
            is_leap_year(1600),
            true,
            "Incorrect return value, expected true when calling is_leap_year(1600)"
        );
    }

    #[test]
    fn test_subject_1500() {
        assert_eq!(
            is_leap_year(1500),
            false,
            "Incorrect return value, expected false when calling is_leap_year(1500)"
        );
    }

    #[test]
    fn test_subject_2004() {
        assert_eq!(
            is_leap_year(2004),
            true,
            "Incorrect return value, expected true when calling is_leap_year(2004)"
        );
    }

    #[test]
    fn test_subject_2003() {
        assert_eq!(
            is_leap_year(2003),
            false,
            "Incorrect return value, expected false when calling is_leap_year(2003)"
        );
    }

    #[test]
    fn test_leap_year_100() {
        let mut rng = rand::thread_rng();

        let year: u32 = rng.gen_range(1..10) * 400 + 100;
        assert_eq!(
            is_leap_year(year),
            false,
            "Incorrect return value, expected false when calling is_leap_year({})",
            year
        );
    }

    #[test]
    fn test_leap_year_200() {
        let mut rng = rand::thread_rng();

        let year: u32 = rng.gen_range(1..10) * 400 + 200;
        assert_eq!(
            is_leap_year(year),
            false,
            "Incorrect return value, expected false when calling is_leap_year({})",
            year
        );
    }

    #[test]
    fn test_leap_year_300() {
        let mut rng = rand::thread_rng();

        let year: u32 = rng.gen_range(1..10) * 400 + 300;
        assert_eq!(
            is_leap_year(year),
            false,
            "Incorrect return value, expected false when calling is_leap_year({})",
            year
        );
    }

    #[test]
    fn test_leap_year_mod_400() {
        let mut rng = rand::thread_rng();

        let year = rng.gen_range(1..10) * 400;
        assert_eq!(
            is_leap_year(year),
            true,
            "Incorrect return value, expected true when calling is_leap_year({})",
            year
        );
    }

    #[test]
    fn test_leap_year_4() {
        let mut rng = rand::thread_rng();
        let mut year;

        loop {
            year = rng.gen_range(1..1000) * 4;
            if year % 100 != 0 {
                break;
            }
        }

        assert_eq!(
            is_leap_year(year),
            true,
            "Incorrect return value, expected true when calling is_leap_year({})",
            year
        );
    }

    #[test]
    fn test_leap_year_false() {
        let mut rng = rand::thread_rng();
        let year = rng.gen_range(1..1000) * 4 + rng.gen_range(1..=3);

        assert_eq!(
            is_leap_year(year),
            false,
            "Incorrect return value, expected true when calling is_leap_year({})",
            year
        );
    }

    macro_rules! generate_test {
        ($name:ident, $month:literal, $leap_year:literal, $expected:literal) => {
            #[test]
            fn $name() {
                let mut rng = rand::thread_rng();
                let mut year;

                loop {
                    year = rng.gen_range(1..1000) * 4;
                    if !$leap_year {
                        year += rng.gen_range(1..=3);
                    }

                    if year % 100 != 0 {
                        break;
                    }
                }

                assert_eq!(
                    num_days_in_month(year, $month),
                    $expected,
                    "Incorrect return value, expected {} when calling num_days_in_month({}, {})",
                    $expected,
                    year,
                    $month
                );
            }
        };
    }

    generate_test!(test_january_leap_year, 1, true, 31);
    generate_test!(test_january_common_year, 1, false, 31);
    generate_test!(test_february_leap_year_0, 2, true, 29);

    #[test]
    fn test_february_leap_year_1() {
        let mut rng = rand::thread_rng();
        let year = rng.gen_range(1..10) * 400;

        assert_eq!(
            num_days_in_month(year, 2),
            29,
            "Incorrect return value, expected 29 when calling num_days_in_month({}, {})",
            year,
            2
        );
    }

    generate_test!(test_february_common_year_0, 2, false, 28);

    #[test]
    fn test_february_common_year_1() {
        let mut rng = rand::thread_rng();
        let year = rng.gen_range(1..10) * 400 + rng.gen_range(1..=3) * 100;

        assert_eq!(
            num_days_in_month(year, 2),
            28,
            "Incorrect return value, expected 28 when calling num_days_in_month({}, {})",
            year,
            2
        );
    }

    generate_test!(test_march_leap_year, 3, true, 31);
    generate_test!(test_march_common_year, 3, false, 31);
    generate_test!(test_april_leap_year, 4, true, 30);
    generate_test!(test_april_common_year, 4, false, 30);
    generate_test!(test_may_leap_year, 5, true, 31);
    generate_test!(test_may_common_year, 5, false, 31);
    generate_test!(test_june_leap_year, 6, true, 30);
    generate_test!(test_june_common_year, 6, false, 30);
    generate_test!(test_july_leap_year, 7, true, 31);
    generate_test!(test_july_common_year, 7, false, 31);
    generate_test!(test_august_leap_year, 8, true, 31);
    generate_test!(test_august_common_year, 8, false, 31);
    generate_test!(test_september_leap_year, 9, true, 30);
    generate_test!(test_september_common_year, 9, false, 30);
    generate_test!(test_october_leap_year, 10, true, 31);
    generate_test!(test_october_common_year, 10, false, 31);
    generate_test!(test_november_leap_year, 11, true, 30);
    generate_test!(test_november_common_year, 11, false, 30);
    generate_test!(test_december_leap_year, 12, true, 31);
    generate_test!(test_december_common_year, 12, false, 31);

    #[test]
    fn test_output() {
        let mut rng = rand::thread_rng();
        let mut buffer: Vec<u8> = Vec::new();

        let year: u32 = rng.gen_range(1500..4000);
        friday_the_13th(&mut buffer, year);
        let output = String::from_utf8_lossy(&buffer);
        let expected = expected_output(year);
        similar_asserts::assert_eq!(
            output,
            expected,
            "Output differs with friday_the_13th() and year {}",
            year
        );
    }

    #[test]
    #[should_panic]
    fn test_friday_the_13th_year_0() {
        let mut buffer: Vec<u8> = Vec::new();

        friday_the_13th(&mut buffer, 0);
    }
}
