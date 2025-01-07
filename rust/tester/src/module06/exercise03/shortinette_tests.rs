
#[cfg(test)]
mod tests {
    use rand::random;

    use super::*;

    #[test]
    fn test_new() {
        let val = random::<i32>();
        let cell = Cellule::new(val);

        assert_eq!(
            cell.get(),
            val,
            "Call to cell.get() returned wrong value - expected: {}, got: {}.",
            val,
            cell.get()
        );
    }

    #[test]
    fn test_replace() {
        let val1 = random::<i32>();
        let mut cell = Cellule::new(val1);
        assert_eq!(
            cell.get(),
            val1,
            "Call to cell.get() returned wrong value - expected: {}, got: {}.",
            val1,
            cell.get()
        );

        let val2 = random::<i32>();
        let old_val = cell.replace(val2);

        assert_eq!(
            old_val, val1,
            "Call to cell.replace() returned wrong value - expected: {}, got: {}.",
            val1, old_val
        );
        assert_eq!(
            cell.get(),
            val2,
            "Call to cell.get() returned wrong value - expected: {}, got: {}.",
            val2,
            cell.get()
        );

        let val3 = random::<i32>();
        let old_val = cell.replace(val3);

        assert_eq!(
            old_val, val2,
            "Call to cell.replace() returned wrong value - expected: {}, got: {}.",
            val2, old_val
        );
        assert_eq!(
            cell.get(),
            val3,
            "Call to cell.get() returned wrong value - expected: {}, got: {}.",
            val3,
            cell.get()
        );
    }

    #[test]
    fn test_set() {
        let val = random::<i32>();
        let mut cell = Cellule::new(val);
        assert_eq!(
            cell.get(),
            val,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val,
            cell.get()
        );

        let val1 = random::<i32>();
        cell.set(val1);
        assert_eq!(
            cell.get(),
            val1,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val1,
            cell.get()
        );
    }

    #[test]
    fn test_get_mut() {
        let val = random::<i32>();

        let mut cell = Cellule::new(val);
        assert_eq!(
            cell.get(),
            val,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val,
            cell.get()
        );

        let val1 = random::<i32>();

        {
            let mut_ref = cell.get_mut();
            *mut_ref = val1;
        }

        assert_eq!(
            cell.get(),
            val1,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val1,
            cell.get()
        );
    }

    #[test]
    fn test_into_inner() {
        let val = random::<i32>();
        let cell = Cellule::new(val);
        let value = cell.into_inner();
        assert_eq!(value, val, "Call to cell.into_inner() returned wrong value - expected: {}, got: {}", val, val);
    }

    #[test]
    fn test_multiple_cells() {
        let val1 = random::<i32>();
        let val2 = random::<i32>();
        let mut cell1 = Cellule::new(val1);
        let mut cell2 = Cellule::new(val2);

        assert_eq!(
            cell1.get(),
            val1,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val1,
            cell1.get()
        );
        assert_eq!(
            cell2.get(),
            val2,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val2,
            cell2.get()
        );

        cell1.set(val2);
        cell2.set(val1);

        assert_eq!(
            cell1.get(),
            val2,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val2,
            cell2.get()
        );
        assert_eq!(
            cell2.get(),
            val1,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val1,
            cell2.get()
        );
    }

    #[test]
    fn test_copy_value() {
        let val = random::<i32>();
        let cell1 = Cellule::new(val);

        let value = cell1.get();
        let cell2 = Cellule::new(value);

        assert_eq!(
            cell1.get(),
            val,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val,
            cell1.get()
        );
        assert_eq!(
            cell2.get(),
            val,
            "Call to cell.get() returned wrong value - expected: {}, got: {}",
            val,
            cell2.get()
        );

        let value = cell2.into_inner();
        assert_eq!(value, val, "Call to cell.into_inner() returned wrong value - expected: {}, got: {}", val, value);
    }
}
