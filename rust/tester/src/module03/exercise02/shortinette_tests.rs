#[cfg(test)]
mod shortinette_tests {
    use super::*;

    #[test]
    fn display() {
        let john = John;

        assert_eq!(format!("{}", john), "Hey! I'm John.");
    }

    #[test]
    fn width() {
        let john = John;

        assert_eq!(
            format!("|{:<30}|", john),
            "|Hey! I'm John.                |"
        );

        assert_eq!(
            format!("|{:>30}|", john),
            "|                Hey! I'm John.|"
        );

        assert_eq!(
            format!("|{:^30}|", john),
            "|        Hey! I'm John.        |"
        );
    }

    #[test]
    fn precision() {
        let john = John;

        assert_eq!(format!("{john:.6}"), "Hey! I");
        assert_eq!(format!("{john:.100}"), "Hey! I'm John.");

        assert_eq!(format!("{john:.0}"), "Don't try to silence me!");
    }

    #[test]
    fn debug() {
        let john = John;

        assert_eq!(format!("{john:?}"), "John, the man himself.");
    }

    #[test]
    fn debug_alternate() {
        let john = John;

        assert_eq!(
            format!("{john:#?}"),
            "John, the man himself. He's handsome AND formidable."
        );
    }

    #[test]
    fn binary() {
        let john = John;

        assert_eq!(format!("{john:b}"), "Bip Boop?");
    }
}
