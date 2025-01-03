#[cfg(test)]
mod shortinette_tests {
    use super::*;

    #[derive(Debug, PartialEq)]
    struct User {
        name: String,
        age: u32,
    }

    impl Record for User {
        fn encode(&self, target: &mut String) -> Result<(), EncodingError> {
            self.name.encode(target)?;
            target.push(',');
            self.age.encode(target)?;

            Ok(())
        }

        fn decode(line: &str) -> Result<Self, DecodingError> {
            let mut split = line.split(',');

            let name = match split.next() {
                Some(value) => Field::decode(value),
                None => Err(DecodingError),
            }?;

            let age = match split.next() {
                Some(value) => Field::decode(value),
                None => Err(DecodingError),
            }?;

            if split.next().is_some() {
                return Err(DecodingError);
            }

            Ok(Self { name, age })
        }
    }

    #[test]
    fn test_encode() {
        let database = [
            User {
                name: "aaa".into(),
                age: 23,
            },
            User {
                name: "bb".into(),
                age: 2,
            },
        ];

        let csv = encode_csv(&database).unwrap();

        assert_eq!(
            csv,
            "\
        aaa,23\n\
        bb,2\n\
        "
        );
    }

    #[test]
    fn test_decode() {
        let csv = "\
        hello,2\n\
        yes,5\n\
        no,100\n\
    ";

        let database: Vec<User> = decode_csv(csv).unwrap();

        assert_eq!(
            database,
            [
                User {
                    name: "hello".into(),
                    age: 2
                },
                User {
                    name: "yes".into(),
                    age: 5
                },
                User {
                    name: "no".into(),
                    age: 100
                },
            ]
        );
    }

    #[test]
    fn decoding_error() {
        let csv = "\
        hello,2\n\
        yes,6\n\
        no,23,hello\n\
    ";

        decode_csv::<User>(csv).unwrap_err();

        let csv = "\
        hello,2\n\
        yes,6\n\
        no\n\
    ";

        decode_csv::<User>(csv).unwrap_err();
    }

    #[test]
    fn empty_line() {
        let csv = "\
        hello,2\n\
        \n\
        bye,3\n\
        ";

        decode_csv::<User>(csv).unwrap_err();
    }

    #[test]
    fn string_encode() {
        let mut line = String::new();
        assert!(("\n".to_string().encode(&mut line).is_err()));
        assert!((",".to_string().encode(&mut line).is_err()));
    }

    #[test]
    fn numbers() {
        let mut line = String::new();

        assert!(12_u8.encode(&mut line).is_ok());
        assert!(48_u16.encode(&mut line).is_ok());
        assert!(75_u32.encode(&mut line).is_ok());
        assert!(4831_u64.encode(&mut line).is_ok());
        assert!(1919_u128.encode(&mut line).is_ok());
        assert!(57329_usize.encode(&mut line).is_ok());
        assert!(73_i8.encode(&mut line).is_ok());
        assert!(874_i16.encode(&mut line).is_ok());
        assert!(4727_i32.encode(&mut line).is_ok());
        assert!(4994_i64.encode(&mut line).is_ok());
        assert!(9448_i128.encode(&mut line).is_ok());
        assert!(9484_isize.encode(&mut line).is_ok());

        assert_eq!("1248754831191957329738744727499494489484", line);
    }

    #[test]
    fn option() {
        let mut line = String::new();

        let hello = Some(String::from("hello"));
        assert!(hello.encode(&mut line).is_ok());
        assert_eq!("hello", line);

        line.clear();

        struct Foo;
        impl Field for Foo {
            fn encode(&self, _target: &mut String) -> Result<(), EncodingError> {
                unreachable!()
            }

            fn decode(_field: &str) -> Result<Self, DecodingError> {
                unreachable!()
            }
        }

        let world: Option<Foo> = None;
        assert!(world.encode(&mut line).is_ok());
        assert_eq!("", line);
    }
}
