#[cfg(test)]
mod shortinette_tests_rust_0107 {
    use ex07::time_manager;
    use ex07::Task;
    use rand::seq::SliceRandom;
    use rand::Rng;
    use std::fmt::{Display, Error, Formatter};

    struct Testcase {
        tasks: Vec<Task>,
        expected: u32,
    }

    // This way we can put the input values in the trace, even if the participants don't add
    // #[derive(Debug)] to the Task struct
    impl Display for Testcase {
        fn fmt(&self, f: &mut Formatter<'_>) -> Result<(), Error> {
            let strings: Vec<String> = self
                .tasks
                .iter()
                .map(|task| {
                    format!(
                        "{{{}, {}, {}}}",
                        task.start_time, task.end_time, task.cookies
                    )
                })
                .collect();
            write!(f, "{:?}", strings)
        }
    }

    #[test]
    fn test_empty() {
        let mut tasks = vec![];
        assert_eq!(time_manager(&mut tasks), 0, "Failed for empty input.");
    }

    // Just multiplying all the input values and the expected output by a random value, to have
    // kinda random inputs without having to have the solution in the public repo.
    // Also shuffling so the order changes.
    fn scale_tests(testcase: &mut Testcase) {
        let mut rng = rand::thread_rng();

        let multiplicator = rng.gen_range(1..100);
        for task in testcase.tasks.iter_mut() {
            task.start_time *= multiplicator;
            task.end_time *= multiplicator;
            task.cookies *= multiplicator;
        }
        testcase.expected *= multiplicator;
        testcase.tasks.shuffle(&mut rng);
    }

    #[test]
    fn test_0() {
        let tasks = vec![
            Task {
                start_time: 0,
                end_time: 3,
                cookies: 10,
            },
            Task {
                start_time: 4,
                end_time: 5,
                cookies: 5,
            },
            Task {
                start_time: 6,
                end_time: 10,
                cookies: 25,
            },
        ];
        let mut testcase = Testcase {
            tasks,
            expected: 40,
        };
        scale_tests(&mut testcase);
        assert_eq!(
            time_manager(&mut testcase.tasks),
            testcase.expected,
            "Failed for {}",
            testcase
        );
    }

    #[test]
    fn test_1() {
        let tasks = vec![
            Task {
                start_time: 0,
                end_time: 3,
                cookies: 10,
            },
            Task {
                start_time: 3,
                end_time: 5,
                cookies: 5,
            },
            Task {
                start_time: 5,
                end_time: 10,
                cookies: 25,
            },
        ];
        let mut testcase = Testcase {
            tasks,
            expected: 40,
        };
        scale_tests(&mut testcase);
        assert_eq!(
            time_manager(&mut testcase.tasks),
            testcase.expected,
            "Failed for {}",
            testcase
        );
    }

    #[test]
    fn test_2() {
        let tasks = vec![
            Task {
                start_time: 0,
                end_time: 5,
                cookies: 10,
            },
            Task {
                start_time: 3,
                end_time: 7,
                cookies: 5,
            },
            Task {
                start_time: 5,
                end_time: 10,
                cookies: 25,
            },
        ];
        let mut testcase = Testcase {
            tasks,
            expected: 35,
        };
        scale_tests(&mut testcase);
        assert_eq!(
            time_manager(&mut testcase.tasks),
            testcase.expected,
            "Failed for {}",
            testcase
        );
    }

    #[test]
    fn test_4() {
        let tasks = vec![
            Task {
                start_time: 0,
                end_time: 5,
                cookies: 1,
            },
            Task {
                start_time: 3,
                end_time: 7,
                cookies: 30,
            },
            Task {
                start_time: 5,
                end_time: 10,
                cookies: 25,
            },
        ];
        let mut testcase = Testcase {
            tasks,
            expected: 30,
        };
        scale_tests(&mut testcase);
        assert_eq!(
            time_manager(&mut testcase.tasks),
            testcase.expected,
            "Failed for {}",
            testcase
        );
    }

    #[test]
    fn test_5() {
        let tasks = vec![
            Task {
                start_time: 0,
                end_time: 5,
                cookies: 1,
            },
            Task {
                start_time: 3,
                end_time: 7,
                cookies: 24,
            },
            Task {
                start_time: 5,
                end_time: 10,
                cookies: 25,
            },
        ];
        let mut testcase = Testcase {
            tasks,
            expected: 26,
        };
        scale_tests(&mut testcase);
        assert_eq!(
            time_manager(&mut testcase.tasks),
            testcase.expected,
            "Failed for {}",
            testcase
        );
    }

    #[test]
    fn test_6() {
        let tasks = vec![
            Task {
                start_time: 2,
                end_time: 25,
                cookies: 10,
            },
            Task {
                start_time: 1,
                end_time: 22,
                cookies: 23,
            },
            Task {
                start_time: 6,
                end_time: 19,
                cookies: 22,
            },
            Task {
                start_time: 6,
                end_time: 16,
                cookies: 24,
            },
            Task {
                start_time: 2,
                end_time: 7,
                cookies: 12,
            },
            Task {
                start_time: 19,
                end_time: 20,
                cookies: 20,
            },
            Task {
                start_time: 16,
                end_time: 18,
                cookies: 23,
            },
            Task {
                start_time: 4,
                end_time: 6,
                cookies: 17,
            },
            Task {
                start_time: 12,
                end_time: 13,
                cookies: 14,
            },
            Task {
                start_time: 12,
                end_time: 15,
                cookies: 23,
            },
        ];
        let mut testcase = Testcase {
            tasks,
            expected: 84,
        };
        scale_tests(&mut testcase);
        assert_eq!(
            time_manager(&mut testcase.tasks),
            testcase.expected,
            "Failed for {}",
            testcase
        );
    }

    #[test]
    fn test_long() {
        let tasks = vec![
            Task {
                start_time: 65,
                end_time: 84,
                cookies: 14,
            },
            Task {
                start_time: 88,
                end_time: 108,
                cookies: 20,
            },
            Task {
                start_time: 49,
                end_time: 113,
                cookies: 16,
            },
            Task {
                start_time: 41,
                end_time: 118,
                cookies: 21,
            },
            Task {
                start_time: 88,
                end_time: 132,
                cookies: 5,
            },
            Task {
                start_time: 105,
                end_time: 136,
                cookies: 9,
            },
            Task {
                start_time: 101,
                end_time: 146,
                cookies: 23,
            },
            Task {
                start_time: 126,
                end_time: 170,
                cookies: 2,
            },
            Task {
                start_time: 166,
                end_time: 180,
                cookies: 11,
            },
            Task {
                start_time: 164,
                end_time: 187,
                cookies: 18,
            },
            Task {
                start_time: 133,
                end_time: 197,
                cookies: 15,
            },
            Task {
                start_time: 181,
                end_time: 212,
                cookies: 15,
            },
            Task {
                start_time: 142,
                end_time: 227,
                cookies: 7,
            },
            Task {
                start_time: 223,
                end_time: 240,
                cookies: 6,
            },
            Task {
                start_time: 184,
                end_time: 241,
                cookies: 24,
            },
            Task {
                start_time: 159,
                end_time: 250,
                cookies: 15,
            },
            Task {
                start_time: 205,
                end_time: 252,
                cookies: 8,
            },
            Task {
                start_time: 192,
                end_time: 256,
                cookies: 1,
            },
            Task {
                start_time: 255,
                end_time: 258,
                cookies: 21,
            },
            Task {
                start_time: 243,
                end_time: 261,
                cookies: 1,
            },
            Task {
                start_time: 264,
                end_time: 272,
                cookies: 7,
            },
            Task {
                start_time: 216,
                end_time: 275,
                cookies: 16,
            },
            Task {
                start_time: 203,
                end_time: 275,
                cookies: 1,
            },
            Task {
                start_time: 200,
                end_time: 283,
                cookies: 16,
            },
            Task {
                start_time: 230,
                end_time: 285,
                cookies: 2,
            },
            Task {
                start_time: 231,
                end_time: 315,
                cookies: 17,
            },
            Task {
                start_time: 313,
                end_time: 323,
                cookies: 5,
            },
            Task {
                start_time: 321,
                end_time: 359,
                cookies: 7,
            },
            Task {
                start_time: 306,
                end_time: 360,
                cookies: 8,
            },
            Task {
                start_time: 324,
                end_time: 363,
                cookies: 20,
            },
            Task {
                start_time: 338,
                end_time: 382,
                cookies: 9,
            },
            Task {
                start_time: 302,
                end_time: 386,
                cookies: 4,
            },
            Task {
                start_time: 368,
                end_time: 430,
                cookies: 12,
            },
            Task {
                start_time: 408,
                end_time: 436,
                cookies: 12,
            },
            Task {
                start_time: 398,
                end_time: 437,
                cookies: 9,
            },
            Task {
                start_time: 360,
                end_time: 440,
                cookies: 9,
            },
            Task {
                start_time: 429,
                end_time: 445,
                cookies: 19,
            },
            Task {
                start_time: 412,
                end_time: 446,
                cookies: 22,
            },
            Task {
                start_time: 396,
                end_time: 451,
                cookies: 8,
            },
            Task {
                start_time: 441,
                end_time: 462,
                cookies: 1,
            },
            Task {
                start_time: 402,
                end_time: 469,
                cookies: 25,
            },
            Task {
                start_time: 418,
                end_time: 482,
                cookies: 1,
            },
            Task {
                start_time: 444,
                end_time: 486,
                cookies: 25,
            },
            Task {
                start_time: 418,
                end_time: 507,
                cookies: 24,
            },
            Task {
                start_time: 470,
                end_time: 513,
                cookies: 20,
            },
            Task {
                start_time: 493,
                end_time: 514,
                cookies: 12,
            },
            Task {
                start_time: 498,
                end_time: 531,
                cookies: 8,
            },
            Task {
                start_time: 485,
                end_time: 532,
                cookies: 18,
            },
            Task {
                start_time: 436,
                end_time: 536,
                cookies: 15,
            },
            Task {
                start_time: 453,
                end_time: 548,
                cookies: 15,
            },
            Task {
                start_time: 541,
                end_time: 567,
                cookies: 14,
            },
            Task {
                start_time: 555,
                end_time: 583,
                cookies: 3,
            },
            Task {
                start_time: 503,
                end_time: 588,
                cookies: 5,
            },
            Task {
                start_time: 529,
                end_time: 595,
                cookies: 20,
            },
            Task {
                start_time: 568,
                end_time: 609,
                cookies: 5,
            },
            Task {
                start_time: 607,
                end_time: 612,
                cookies: 21,
            },
            Task {
                start_time: 524,
                end_time: 625,
                cookies: 1,
            },
            Task {
                start_time: 584,
                end_time: 630,
                cookies: 9,
            },
            Task {
                start_time: 581,
                end_time: 630,
                cookies: 22,
            },
            Task {
                start_time: 545,
                end_time: 635,
                cookies: 5,
            },
            Task {
                start_time: 553,
                end_time: 639,
                cookies: 21,
            },
            Task {
                start_time: 633,
                end_time: 668,
                cookies: 12,
            },
            Task {
                start_time: 642,
                end_time: 677,
                cookies: 20,
            },
            Task {
                start_time: 672,
                end_time: 696,
                cookies: 19,
            },
            Task {
                start_time: 666,
                end_time: 698,
                cookies: 20,
            },
            Task {
                start_time: 665,
                end_time: 729,
                cookies: 24,
            },
            Task {
                start_time: 708,
                end_time: 743,
                cookies: 15,
            },
            Task {
                start_time: 651,
                end_time: 752,
                cookies: 6,
            },
            Task {
                start_time: 745,
                end_time: 762,
                cookies: 21,
            },
            Task {
                start_time: 676,
                end_time: 773,
                cookies: 24,
            },
            Task {
                start_time: 723,
                end_time: 777,
                cookies: 9,
            },
            Task {
                start_time: 754,
                end_time: 779,
                cookies: 23,
            },
            Task {
                start_time: 712,
                end_time: 782,
                cookies: 17,
            },
            Task {
                start_time: 755,
                end_time: 785,
                cookies: 23,
            },
            Task {
                start_time: 791,
                end_time: 800,
                cookies: 5,
            },
            Task {
                start_time: 703,
                end_time: 802,
                cookies: 5,
            },
            Task {
                start_time: 798,
                end_time: 820,
                cookies: 21,
            },
            Task {
                start_time: 822,
                end_time: 834,
                cookies: 7,
            },
            Task {
                start_time: 802,
                end_time: 836,
                cookies: 23,
            },
            Task {
                start_time: 792,
                end_time: 842,
                cookies: 11,
            },
            Task {
                start_time: 830,
                end_time: 844,
                cookies: 24,
            },
            Task {
                start_time: 794,
                end_time: 852,
                cookies: 6,
            },
            Task {
                start_time: 797,
                end_time: 857,
                cookies: 22,
            },
            Task {
                start_time: 762,
                end_time: 863,
                cookies: 19,
            },
            Task {
                start_time: 833,
                end_time: 869,
                cookies: 16,
            },
            Task {
                start_time: 887,
                end_time: 890,
                cookies: 21,
            },
            Task {
                start_time: 849,
                end_time: 898,
                cookies: 1,
            },
            Task {
                start_time: 824,
                end_time: 899,
                cookies: 8,
            },
            Task {
                start_time: 874,
                end_time: 932,
                cookies: 9,
            },
            Task {
                start_time: 885,
                end_time: 933,
                cookies: 16,
            },
            Task {
                start_time: 909,
                end_time: 959,
                cookies: 1,
            },
            Task {
                start_time: 933,
                end_time: 970,
                cookies: 14,
            },
            Task {
                start_time: 936,
                end_time: 980,
                cookies: 24,
            },
            Task {
                start_time: 975,
                end_time: 984,
                cookies: 9,
            },
            Task {
                start_time: 956,
                end_time: 993,
                cookies: 1,
            },
            Task {
                start_time: 949,
                end_time: 994,
                cookies: 19,
            },
            Task {
                start_time: 921,
                end_time: 1000,
                cookies: 1,
            },
            Task {
                start_time: 984,
                end_time: 1030,
                cookies: 21,
            },
            Task {
                start_time: 955,
                end_time: 1042,
                cookies: 13,
            },
            Task {
                start_time: 994,
                end_time: 1043,
                cookies: 1,
            },
        ];

        let mut testcase = Testcase {
            tasks,
            expected: 395,
        };

        scale_tests(&mut testcase);

        assert_eq!(
            time_manager(&mut testcase.tasks),
            testcase.expected,
            "Failed for {}",
            testcase
        );
    }
}
