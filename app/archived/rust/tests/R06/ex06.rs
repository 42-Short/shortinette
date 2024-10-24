#[cfg(test)]
mod shortinette_tests_rust_0605 {
    use super::*;
    use std::ffi::CString;

    #[test]
    fn test_create_database() {
        let db = Database::new();
        assert!(db.is_ok());
        let db = db.unwrap();
        assert_eq!(db.next_user_id, 1);
        assert_eq!(db.count, 0);
        assert!(db.allocated > 1);
    }

    #[test]
    fn test_create_user_success() {
        let mut db = Database::new().expect("Failed to create database");
        let name = CString::new("Alice").expect("CString::new failed");
        let user_id = db.create_user(&name);
        assert!(user_id.is_ok());
        assert_eq!(user_id.unwrap(), 1);
    }

    #[test]
    fn test_create_multiple_users() {
        let mut db = Database::new().expect("Failed to create database");

        let name1 = CString::new("Alice").expect("CString::new failed");
        let name2 = CString::new("Bob").expect("CString::new failed");
        let name3 = CString::new("Charlie").expect("CString::new failed");

        let id1 = db.create_user(&name1).expect("Failed to create user Alice");
        let id2 = db.create_user(&name2).expect("Failed to create user Bob");
        let id3 = db
            .create_user(&name3)
            .expect("Failed to create user Charlie");

        assert_eq!(id1, 1);
        assert_eq!(id2, 2);
        assert_eq!(id3, 3);
    }

    #[test]
    fn test_delete_user_success() {
        let mut db = Database::new().expect("Failed to create database");
        let name = CString::new("Alice").expect("CString::new failed");
        let user_id = db.create_user(&name).expect("Failed to create user Alice");

        let result = db.delete_user(user_id);
        assert!(result.is_ok());
    }

    #[test]
    fn test_delete_nonexistent_user() {
        let mut db = Database::new().expect("Failed to create database");
        let result = db.delete_user(999);
        assert!(result.is_err());
        assert_eq!(result.err().unwrap(), Error::ErrUnknownId);
    }

    #[test]
    fn test_get_user_success() {
        let mut db = Database::new().expect("Failed to create database");
        let name = CString::new("Alice").expect("CString::new failed");
        let user_id = db.create_user(&name).expect("Failed to create user Alice");

        let user = db.get_user(user_id);
        assert!(user.is_ok());
        let user = user.unwrap();
        let user_name = unsafe { CStr::from_ptr(user.name) };
        assert_eq!(user.id, user_id);
        assert_eq!(user_name.to_str().unwrap(), "Alice");
    }

    #[test]
    fn test_get_nonexistent_user() {
        let db = Database::new().expect("Failed to create database");
        let result = db.get_user(999);
        assert!(result.is_err());
        assert_eq!(result.err().unwrap(), Error::ErrUnknownId);
    }

    #[test]
    fn test_database_drop() {
        let db = Database::new().expect("Failed to create database");
        drop(db);
    }
}