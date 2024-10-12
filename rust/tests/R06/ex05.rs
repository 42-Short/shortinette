#[cfg(test)]
mod shortinette_tests_rust_0605 {
    use super::*;

    #[test]
    fn test_independent_cloning() {
        let mut tableau = Tableau::new();
        tableau.push(String::from("Hello"));
        tableau.push(String::from("Rust"));
    
        let mut tableau_clone = tableau.clone();
        tableau_clone.push(String::from("World"));
    
        assert_eq!(tableau.len(), 2);
        assert_eq!(tableau_clone.len(), 3);
    
        assert_eq!(tableau.pop(), Some(String::from("Rust")));
        assert_eq!(tableau_clone.pop(), Some(String::from("World")));
    }

    #[test]
    fn test_iteration_with_for_loop() {
        let mut tableau = Tableau::new();
        tableau.push(10);
        tableau.push(20);
        tableau.push(30);
    
        let mut sum = 0;
        for value in &*tableau {
            sum += value;
        }
    
        assert_eq!(sum, 60);
    }
    
    #[test]
    fn test_deref_mut() {
        let mut tableau = Tableau::new();
        tableau.push(1);
        tableau.push(2);
        tableau.push(3);
    
        {
            let tableau_mut: &mut [i32] = &mut *tableau;
            tableau_mut[0] = 10;
            tableau_mut[2] = 30;
        }
    
        assert_eq!(tableau.pop(), Some(30));
        assert_eq!(tableau.pop(), Some(2));
        assert_eq!(tableau.pop(), Some(10));
    }
    
    #[test]
    fn test_with_tuples() {
        let mut tableau = Tableau::new();
        tableau.push((1, "one"));
        tableau.push((2, "two"));
        tableau.push((3, "three"));
    
        assert_eq!(tableau.len(), 3);
    
        let tuple = tableau.pop();
        assert_eq!(tuple, Some((3, "three")));
        assert_eq!(tableau.len(), 2);
    }
    
    #[test]
    fn test_non_copyable_types() {
        let mut tableau = Tableau::new();
        tableau.push(Box::new(10));
        tableau.push(Box::new(20));
    
        assert_eq!(tableau.len(), 2);
    
        let value = tableau.pop().unwrap();
        assert_eq!(*value, 20);
    
        let value = tableau.pop().unwrap();
        assert_eq!(*value, 10);
    }    

    #[test]
    fn test_large_push() {
        let mut tableau = Tableau::new();
    
        for i in 0..1000 {
            tableau.push(i);
        }
    
        assert_eq!(tableau.len(), 1000);
    
        for i in (0..1000).rev() {
            assert_eq!(tableau.pop(), Some(i));
        }
    
        assert!(tableau.is_empty());
    }

    #[test]
    fn test_empty_tableau_behavior() {
        let mut tableau: Tableau<i32> = Tableau::new();
    
        assert!(tableau.is_empty());
        assert_eq!(tableau.len(), 0);
        assert_eq!(tableau.pop(), None);
    
        tableau.clear();
        assert_eq!(tableau.len(), 0);
        assert!(tableau.is_empty());
        assert_eq!(tableau.pop(), None);
    }

    #[test]
    fn test_push_and_pop_multiple() {
        let mut tableau = Tableau::new();

        tableau.push(1);
        tableau.push(2);
        tableau.push(3);
        tableau.push(4);
        tableau.push(5);

        assert_eq!(tableau.len(), 5);

        assert_eq!(tableau.pop(), Some(5));
        assert_eq!(tableau.pop(), Some(4));
        assert_eq!(tableau.pop(), Some(3));
        assert_eq!(tableau.len(), 2);

        tableau.push(6);
        assert_eq!(tableau.len(), 3);
        assert_eq!(tableau.pop(), Some(6));
        assert_eq!(tableau.len(), 2);
    }

    #[test]
    fn test_new_tableau() {
        let tableau: Tableau<i32> = Tableau::new();
        assert_eq!(tableau.len(), 0);
        assert!(tableau.is_empty());
    }

    #[test]
    fn test_pop() {
        let mut tableau = Tableau::new();
        assert_eq!(tableau.pop(), None);

        tableau.push(5);
        tableau.push(15);
        assert_eq!(tableau.len(), 2);

        let popped = tableau.pop();
        assert_eq!(popped, Some(15));
        assert_eq!(tableau.len(), 1);

        let popped = tableau.pop();
        assert_eq!(popped, Some(5));
        assert_eq!(tableau.len(), 0);

        assert!(tableau.is_empty());
        assert_eq!(tableau.pop(), None);
    }

    #[test]
    fn test_clear() {
        let mut tableau = Tableau::new();
        tableau.push(1);
        tableau.push(2);
        tableau.push(3);
        assert_eq!(tableau.len(), 3);

        tableau.clear();
        assert_eq!(tableau.len(), 0);
        assert!(tableau.is_empty());
        assert_eq!(tableau.pop(), None);
    }

    #[test]
    fn test_into_iterator() {
        let mut tableau = Tableau::new();
        tableau.push("a");
        tableau.push("b");
        tableau.push("c");

        let mut iter = tableau.into_iter();
        assert_eq!(iter.next(), Some("a"));
        assert_eq!(iter.next(), Some("b"));
        assert_eq!(iter.next(), Some("c"));
        assert_eq!(iter.next(), None);
    }

    #[test]
    fn test_clone() {
        let mut tableau = Tableau::new();
        tableau.push(String::from("Hello"));
        tableau.push(String::from("World"));

        let tableau_clone = tableau.clone();
        assert_eq!(tableau.len(), tableau_clone.len());

        tableau.push(String::from("!"));
        assert_eq!(tableau.len(), 3);
        assert_eq!(tableau_clone.len(), 2);
    }

    #[test]
    fn test_subject() {
        let mut a = Tableau::new();
        a.push(1);
        a.push(2);
        a.push(4);
        let b = a.clone();

        for it in b {
            println!("{it}");
        }
        let c: &[i32] = &*a;
        assert_eq!(c, [1, 2, 4]);
    }
}