#[cfg(test)]
mod shortinette_tests_rust_0104 {
    use ex04::sort_boxes;
    use rand::prelude::SliceRandom;
    use rand::Rng;
    use std::collections::HashMap;

    fn get_next_value(current: &[u32; 2]) -> [u32; 2] {
        let mut rng = rand::thread_rng();

        let width = current[0] + (rng.gen_range(0..10) > 4) as u32;
        let height = current[1] + (rng.gen_range(0..10) > 4) as u32;

        [width, height]
    }

    fn add_boxes_shuffle(boxes: &mut Vec<[u32; 2]>, amount: usize, sortable: bool) {
        let mut rng = rand::thread_rng();

        for _ in 0..amount {
            let prev = boxes.last().unwrap();
            boxes.push(get_next_value(prev));
        }

        if !sortable {
            let position = rng.gen_range(0..boxes.len() - 2);
            boxes[position + 1] = [boxes[position][0] + 1, boxes[position][1]];
            boxes[position + 2] = [boxes[position][0], boxes[position][1] + 1];
        }

        boxes.shuffle(&mut rng);
    }

    // Using a HashMap to check if all the values returned by the function are the same
    // as the original input, to avoid having the solution itself in the public repo.
    fn is_sorted(boxes: &[[u32; 2]], initial: &[[u32; 2]]) -> bool {
        let mut box_counter = HashMap::new();

        for item in initial {
            *box_counter.entry(item).or_insert(0) += 1;
        }

        for window in boxes.windows(2) {
            if let [a, b] = window {
                if b[0] > a[0] || b[1] > a[1] {
                    return false;
                }
            }
        }

        for item in boxes {
            if let Some(count) = box_counter.get_mut(item) {
                *count -= 1;
                if *count == 0 {
                    box_counter.remove(item);
                }
            } else {
                return false;
            }
        }

        box_counter.is_empty()
    }

    #[test]
    fn test_empty() {
        let mut boxes = [];
        sort_boxes(&mut boxes);
        assert!(boxes.is_empty(), "Failed for an empty list as input");
    }

    #[test]
    fn test_regular() {
        let mut boxes = vec![[1, 2]];
        add_boxes_shuffle(&mut boxes, 100, true);
        let initial = boxes.clone();

        let original_str = format!("{:?}", boxes);
        sort_boxes(&mut boxes);
        assert_eq!(
            is_sorted(&boxes, &initial),
            true,
            "Boxes haven't been sorted correctly\nInput: {}\nAfter Sorting: {:?}",
            original_str,
            boxes
        );
    }

    #[test]
    fn test_0_0_box() {
        let mut boxes = vec![[0, 0]];
        add_boxes_shuffle(&mut boxes, 100, true);
        let initial = boxes.clone();

        let original_str = format!("{:?}", boxes);
        sort_boxes(&mut boxes);
        assert_eq!(
            is_sorted(&boxes, &initial),
            true,
            "Boxes haven't been sorted correctly\nInput: {}\nAfter Sorting: {:?}",
            original_str,
            boxes
        );
    }

    #[test]
    #[should_panic]
    fn test_unsortable() {
        let mut boxes = vec![[1, 2]];
        add_boxes_shuffle(&mut boxes, 100, false);

        let original_str = format!("{:?}", boxes);
        sort_boxes(&mut boxes);
    }
}
