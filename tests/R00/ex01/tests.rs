
mod min;

#[test]
fn test_0() {
    assert_eq!(min::min(1, 2), 1);
}

#[test]
fn test_1() {
    assert_eq!(min::min(2, 1), 1);
}

#[test]
fn test_2() {
    assert_eq!(min::min(1, 1), 1);
}

#[test]
fn test_3() {
    assert_eq!(min::min(-1, 0), -1);
}