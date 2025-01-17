#[cfg(test)]
mod shortinette_tests {
    use super::*;

    fn assert_list<T: Eq + std::fmt::Debug>(list: &List<T>, vec: &[T]) {
        assert_eq!(list.count(), vec.len());

        vec.iter()
            .enumerate()
            .for_each(|(index, element)| assert_eq!(&list[index], element));
    }

    #[test]
    fn default_list_is_empty() {
        let list: List<i32> = Default::default();
        assert_eq!(list.count(), 0);
    }

    #[test]
    fn cloned_list_are_equal() {
        let mut list = List::new();
        list.push_back(String::from("Hello"));
        list.push_back(String::from("World"));

        let cloned = list.clone();
        assert_eq!(cloned.count(), list.count());
        assert_eq!(&list[0], &cloned[0]);
        assert_eq!(&list[1], &cloned[1]);
    }

    #[test]
    #[should_panic(expected = "tried to access out of bound index 10")]
    fn out_of_bound_access_panics() {
        let mut list: List<u32> = List::new();
        list.push_back(1);
        list.push_back(2);
        list.push_back(3);

        assert_eq!(list[10], 42);
    }

    #[test]
    fn in_bound_access() {
        let mut list: List<u32> = List::new();
        list.push_back(1);
        list.push_back(2);
        list.push_back(3);

        assert_eq!(list[0], 1);
        assert_eq!(list[1], 2);
        assert_eq!(list[2], 3);
    }

    #[test]
    fn in_bound_access_mut() {
        let mut list: List<u32> = List::new();
        list.push_back(1);
        list.push_back(2);
        list.push_back(3);

        list[0] += 10;
        list[1] += 10;
        list[2] += 10;

        assert_eq!(list[0], 11);
        assert_eq!(list[1], 12);
        assert_eq!(list[2], 13);
    }

    #[test]
    fn new() {
        let list: List<i32> = List::new();

        assert!(
            list.head.is_none(),
            "List::new() should return an empty list"
        );
    }

    #[test]
    fn push_back() {
        let mut list: List<i32> = List::new();

        for i in 0..5 {
            list.push_back(i);
        }

        let expect: Vec<i32> = (0..5).collect();
        assert_list(&list, &expect);
    }

    #[test]
    fn push_front() {
        let mut list: List<i32> = List::new();

        for i in 0..5 {
            list.push_front(i);
        }

        let expect: Vec<i32> = (0..5).rev().collect();
        assert_list(&list, &expect);
    }

    #[test]
    fn count() {
        let mut list: List<i32> = List::new();

        assert_eq!(list.count(), 0);

        for i in 1..100 {
            list.push_back(i);
            assert_eq!(list.count(), i as usize);
        }

        for i in (0..99).rev() {
            list.remove_front();
            assert_eq!(list.count(), i);
        }
    }

    #[test]
    fn get() {
        let mut list: List<i32> = List::new();

        assert!(list.get(0).is_none());

        for i in 0..100 {
            list.push_back(i);
            assert_eq!(list.get(i as usize), Some(&i));
        }

        for i in 0..100 {
            assert_eq!(list.get(i as usize), Some(&i));
        }

        for i in (0..100).rev() {
            assert_eq!(list.get(i as usize), Some(&i));
            list.remove_back();
        }
    }

    #[test]
    fn get_mut() {
        let mut list: List<i32> = List::new();

        assert!(list.get(0).is_none());

        for i in 0..100 {
            list.push_back(i);
            let element = list.get_mut(i as usize).unwrap();
            *element += 100;
        }

        let expect: Vec<i32> = (100..200).collect();
        assert_list(&list, &expect);
    }

    #[test]
    fn remove_back() {
        let mut list: List<i32> = List::new();

        for i in 0..5 {
            list.push_back(i);
        }

        for i in 0..5 {
            list.remove_back();

            let expect: Vec<i32> = (0..5 - i - 1).collect();
            assert_list(&list, &expect);
        }
    }

    #[test]
    fn remove_front() {
        let mut list: List<i32> = List::new();

        for i in 0..5 {
            list.push_back(i);
        }

        for i in 0..5 {
            list.remove_front();

            let expect: Vec<i32> = (i + 1..5).collect();
            assert_list(&list, &expect);
        }
    }

    #[test]
    fn clear() {
        let mut list: List<i32> = List::new();

        assert!(list.head.is_none());
        list.clear();
        assert!(list.head.is_none());

        // Test with one element
        {
            list.push_back(0);
            assert!(list.head.is_some());
            list.clear();
            assert!(list.head.is_none());
        }

        // Test with multiple elements
        {
            for i in 0..100 {
                list.push_back(i);
            }

            assert_eq!(list.count(), 100);
            list.clear();
            assert!(list.head.is_none());
        }
    }
}
