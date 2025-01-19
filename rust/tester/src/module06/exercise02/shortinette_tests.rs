#[cfg(test)]
mod tests {
    use rand::{distributions::Alphanumeric, random, thread_rng, Rng};

    use super::*;

    #[derive(Clone, Debug, PartialEq)]
    struct Point {
        x: u32,
        y: u32,
    }

    #[test]
    fn new_and_deref() {
        let x = random::<u32>();
        let y = random::<u32>();
        let point_in_carton = Carton::new(Point { x, y });
        assert_eq!(point_in_carton.x, x);
        assert_eq!(point_in_carton.y, y);
    }

    #[test]
    fn deref_mut() {
        let mut point_in_carton = Carton::new(Point { x: 0, y: 0 });
        let x = random::<u32>();
        let y = random::<u32>();
        point_in_carton.x = x;
        point_in_carton.y = y;
        assert_eq!(point_in_carton.x, x);
        assert_eq!(point_in_carton.y, y);
    }

    #[test]
    fn into_inner() {
        let x = random::<u32>();
        let y = random::<u32>();
        let point_in_carton = Carton::new(Point { x, y });
        let point = point_in_carton.into_inner();
        assert_eq!(point, Point { x, y });
    }

    #[test]
    fn clone() {
        let x = random::<u32>();
        let y = random::<u32>();
        let point_in_carton = Carton::new(Point { x, y });

        let mut cloned_point = point_in_carton.clone();

        assert_eq!(cloned_point.x, x);
        assert_eq!(cloned_point.y, y);

        let x1 = random::<u32>();
        let y1 = random::<u32>();

        cloned_point.x = x1;
        cloned_point.y = y1;

        assert_eq!(point_in_carton.x, x);
        assert_eq!(point_in_carton.y, y);
        assert_eq!(cloned_point.x, x1, "Does your clone() make deep copies?");
        assert_eq!(cloned_point.y, y1, "Does your clone() make deep copies?");
    }

    #[test]
    fn drop() {
        struct TestDrop {
            dropped: std::rc::Rc<std::cell::Cell<bool>>,
        }

        impl Drop for TestDrop {
            fn drop(&mut self) {
                self.dropped.set(true);
            }
        }

        let dropped = std::rc::Rc::new(std::cell::Cell::new(false));
        let drop = TestDrop {
            dropped: dropped.clone(),
        };
        {
            let _ = Carton::new(drop);
        }
        assert!(
            dropped.get(),
            "Does your drop() implementation also drop the Carton's content?"
        );
    }

    #[test]
    fn i32() {
        let val = random::<i32>();
        let int_in_carton = Carton::new(val);
        assert_eq!(*int_in_carton, val);
    }

    #[test]
    fn carton_of_str() {
        let val: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u16>() as usize)
            .map(char::from)
            .collect();

        let str_in_carton = Carton::new(&val);
        assert_eq!(*str_in_carton.into_inner(), val);
    }

    #[test]
    fn carton_of_vec() {
        let val: Vec<u8> = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u16>() as usize)
            .collect();

        let vec_in_carton = Carton::new(&val);

        assert_eq!(vec_in_carton.len(), val.len());
        assert_eq!(*vec_in_carton.into_inner(), val);
    }

    #[test]
    fn large_carton() {
        let val: Vec<u8> = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take((random::<u16>() as usize).max(256 * 256).min(512 * 512))
            .collect();

        let large_carton = Carton::new(&val);

        assert_eq!(large_carton.len(), val.len());
        assert_eq!(*large_carton.into_inner(), val);
    }
}
