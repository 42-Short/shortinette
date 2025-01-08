#[cfg(test)]
mod tests {
    use rand::{distributions::Alphanumeric, random, thread_rng, Rng};

    use super::*;

    fn rand_string() -> String {
        thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect()
    }

    #[test]
    fn test_independent_cloning() {
        let mut tableau = Tableau::new();

        let s0 = rand_string();
        let s1 = rand_string();
        let s2 = rand_string();

        tableau.push(s0);
        tableau.push(s1.clone());

        let mut tableau_clone = tableau.clone();
        tableau_clone.push(s2.clone());

        assert_eq!(tableau.len(), 2);
        assert_eq!(tableau_clone.len(), 3);

        assert_eq!(tableau.pop().unwrap(), s1);
        assert_eq!(tableau_clone.pop().unwrap(), s2);
    }

    #[test]
    fn test_iteration_with_for_loop() {
        let mut tableau = Tableau::new();

        let slice: Vec<u8> = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(random::<u16>() as usize)
            .map(u8::from)
            .collect();

        for ueight in slice.clone() {
            tableau.push(ueight);
        }

        assert_eq!(tableau.len(), slice.len());

        for idx in 0..tableau.len() {
            assert_eq!(tableau[idx], slice[idx]);
        }
    }

    #[test]
    fn test_deref_mut() {
        let mut tableau = Tableau::new();

        let n0 = random::<i32>();
        let n1 = random::<i32>();
        let n2 = random::<i32>();
        let n3 = random::<i32>();
        let n4 = random::<i32>();

        tableau.push(n0);
        tableau.push(n1);
        tableau.push(n2);

        {
            let tableau_mut: &mut [i32] = &mut *tableau;
            tableau_mut[0] = n3;
            tableau_mut[2] = n4;
        }

        assert_eq!(tableau.pop(), Some(n4));
        assert_eq!(tableau.pop(), Some(n1));
        assert_eq!(tableau.pop(), Some(n3));
    }

    #[test]
    fn test_with_tuples() {
        let mut tableau = Tableau::new();

        let s0 = rand_string();
        let s1 = rand_string();
        let s2 = rand_string();

        let n0 = random::<i32>();
        let n1 = random::<i32>();
        let n2 = random::<i32>();

        tableau.push((n0, s0));
        tableau.push((n1, s1));
        tableau.push((n2, s2.clone()));

        assert_eq!(tableau.len(), 3);

        let tuple = tableau.pop();
        assert_eq!(tuple, Some((n2, s2)));
        assert_eq!(tableau.len(), 2);
    }

    #[test]
    fn test_i32() {
        let mut tableau = Tableau::new();

        let n0 = random::<i32>();
        let n1 = random::<i32>();

        tableau.push(n0);
        tableau.push(n1);

        assert_eq!(tableau.len(), 2);

        let value = tableau.pop().unwrap();
        assert_eq!(value, n1);

        let value = tableau.pop().unwrap();
        assert_eq!(value, n0);
    }

    #[test]
    fn test_large_push() {
        let mut tableau = Tableau::new();

        let len = random::<u16>();
        let mut cpy = Vec::new();

        for _ in 0..len {
            let val = random::<u64>();
            tableau.push(val);
            cpy.push(val);
        }

        assert_eq!(tableau.len(), len as usize);

        for _ in (0..len).rev() {
            assert_eq!(tableau.pop().unwrap(), cpy.pop().unwrap());
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
    fn test_new_tableau() {
        let tableau: Tableau<i32> = Tableau::new();
        assert_eq!(tableau.len(), 0);
        assert!(tableau.is_empty());
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

        let s0 = rand_string();
        let s1 = rand_string();
        let s2 = rand_string();

        tableau.push(s0.clone());
        tableau.push(s1.clone());
        tableau.push(s2.clone());

        let mut iter = tableau.into_iter();
        assert_eq!(iter.next(), Some(s0));
        assert_eq!(iter.next(), Some(s1));
        assert_eq!(iter.next(), Some(s2));
        assert_eq!(iter.next(), None);
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
