#[cfg(test)]
mod shortinette_tests_rust_0602 {
    use super::*;

    #[derive(Clone, Debug, PartialEq)]
    struct Point {
        x: u32,
        y: u32,
    }

    #[test]
    fn test_new_and_deref() {
        let point_in_carton = Carton::new(Point { x: 1, y: 2 });
        assert_eq!(point_in_carton.x, 1);
        assert_eq!(point_in_carton.y, 2);
    }

    #[test]
    fn test_deref_mut() {
        let mut point_in_carton = Carton::new(Point { x: 1, y: 2 });
        point_in_carton.x = 10;
        point_in_carton.y = 20;
        assert_eq!(point_in_carton.x, 10);
        assert_eq!(point_in_carton.y, 20);
    }

    #[test]
    fn test_into_inner() {
        let point_in_carton = Carton::new(Point { x: 3, y: 4 });
        let point = point_in_carton.into_inner();
        assert_eq!(point, Point { x: 3, y: 4 });
    }

    #[test]
    fn test_clone() {
        let point_in_carton = Carton::new(Point { x: 5, y: 6 });
        let mut cloned_point = point_in_carton.clone();

        assert_eq!(cloned_point.x, 5);
        assert_eq!(cloned_point.y, 6);

        cloned_point.x = 7;
        cloned_point.y = 8;

        assert_eq!(point_in_carton.x, 5);
        assert_eq!(point_in_carton.y, 6);
        assert_eq!(cloned_point.x, 7);
        assert_eq!(cloned_point.y, 8);
    }

    #[test]
    fn test_drop() {
        struct TestDrop {
            dropped: std::rc::Rc<std::cell::Cell<bool>>,
        }

        impl Drop for TestDrop {
            fn drop(&mut self) {
                self.dropped.set(true);
            }
        }

        let dropped = std::rc::Rc::new(std::cell::Cell::new(false));
        let test_drop = TestDrop {
            dropped: dropped.clone(),
        };
        {
            let _carton = Carton::new(test_drop);
        }
        assert!(dropped.get());
    }

    #[test]
    fn test_carton_of_i32() {
        let int_in_carton = Carton::new(10);
        assert_eq!(*int_in_carton, 10);
    }

    #[test]
    fn test_carton_of_str() {
        let str_in_carton = Carton::new(String::from("Hello, world!"));
        assert_eq!(str_in_carton.as_str(), "Hello, world!");
    }

    #[test]
    fn test_carton_of_vec() {
        let vec_in_carton = Carton::new(vec![1, 2, 3, 4]);
        assert_eq!(vec_in_carton.len(), 4);
        assert_eq!(vec_in_carton[0], 1);
        assert_eq!(vec_in_carton[3], 4);
    }

    #[test]
    fn test_large_carton() {
        let large_carton = Carton::new(vec![0u8; 1024 * 1024]);
        assert_eq!(large_carton.len(), 1024 * 1024);
        assert_eq!(large_carton.iter().all(|&x| x == 0), true);
    }

    #[test]
    fn test_subject() {
        #[derive(Clone)]
        struct Point {
            x: u32,
            y: u32,
        }
        let point_in_carton = Carton::new(Point { x: 1, y: 2 });
        assert_eq!(point_in_carton.x, 1);
        assert_eq!(point_in_carton.y, 2);

        let mut another_point = point_in_carton.clone();
        another_point.x = 2;
        another_point.y = 3;
        assert_eq!(another_point.x, 2);
        assert_eq!(another_point.y, 3);
        assert_eq!(point_in_carton.x, 1);
        assert_eq!(point_in_carton.y, 2);
    }
}