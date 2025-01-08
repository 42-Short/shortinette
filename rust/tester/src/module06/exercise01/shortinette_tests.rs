#[cfg(test)]
mod tests {
    use rand::random;

    use super::*;

    #[test]
    fn transmute_iron() {
        let iron = random::<Iron>();

        let mut nuggets: [GoldNugget; 2] = PhilosopherStone.transmute_iron(iron);

        assert!(
            nuggets.len() == 2,
            "Expected 2 nuggets from Iron transmutation, got {}.",
            nuggets.len()
        );

        let mut expected_nuggets: [GoldNugget; 2] =
            [(iron >> 16) as GoldNugget, (iron & 0xFFFF) as GoldNugget];

        expected_nuggets.sort();
        nuggets.sort();

        assert_eq!(
            nuggets, expected_nuggets,
            "Nuggets do not match expected values."
        );
    }

    #[test]
    fn transmute_mercure() {
        let mercure = random::<Mercure>();

        let mut nuggets: [GoldNugget; 4] = PhilosopherStone.transmute_mercure(mercure);

        assert!(
            nuggets.len() == 4,
            "Expected 4 nuggets from Mercure transmutation, got: {}.",
            nuggets.len()
        );

        let mut expected_nuggets: [GoldNugget; 4] = [
            (mercure >> 48) as GoldNugget,
            ((mercure >> 32) & 0xFFFF) as GoldNugget,
            ((mercure >> 16) & 0xFFFF) as GoldNugget,
            (mercure & 0xFFFF) as GoldNugget,
        ];

        nuggets.sort();
        expected_nuggets.sort();

        assert_eq!(
            nuggets, expected_nuggets,
            "Nuggets do not match expected values."
        );
    }

    #[test]
    fn transmute_metal_iron() {
        let iron: Iron = random::<Iron>();

        let nuggets: &Gold = PhilosopherStone.transmute_metal(&iron);

        assert!(
            nuggets.len() == 2,
            "Expected 2 nuggets from Iron transmutation, got {}.",
            nuggets.len()
        );

        let mut expected_nuggets: [GoldNugget; 2] =
            [(iron >> 16) as GoldNugget, (iron & 0xFFFF) as GoldNugget];
        expected_nuggets.sort();

        let mut nuggets_vec = nuggets.to_vec();
        nuggets_vec.sort();

        assert_eq!(
            nuggets_vec, expected_nuggets,
            "Nuggets do not match expected values."
        );
    }

    #[test]
    fn transmute_metal_mercure() {
        let mercure: Mercure = random::<Mercure>();

        let nuggets: &Gold = PhilosopherStone.transmute_metal(&mercure);

        assert!(
            nuggets.len() == 4,
            "Expected 2 nuggets from Iron transmutation, got {}.",
            nuggets.len()
        );

        let mut expected_nuggets: [GoldNugget; 4] = [
            (mercure >> 48) as GoldNugget,
            ((mercure >> 32) & 0xFFFF) as GoldNugget,
            ((mercure >> 16) & 0xFFFF) as GoldNugget,
            (mercure & 0xFFFF) as GoldNugget,
        ];
        expected_nuggets.sort();

        let mut nuggets_vec = nuggets.to_vec();
        nuggets_vec.sort();

        assert_eq!(
            nuggets_vec, expected_nuggets,
            "Nuggets do not match expected values."
        );
    }
}
