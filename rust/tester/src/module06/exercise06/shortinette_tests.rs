#[cfg(test)]
mod tests {
    use rand::{distributions::Alphanumeric, thread_rng, Rng};

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
        let mut db = Database::new().expect("Failed to create database.");

        let uid: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let name = CString::new(uid).expect("CString::new failed.");
        let user_id = db.create_user(&name);
        assert!(user_id.is_ok());
        assert_eq!(user_id.unwrap(), 1);
    }

    #[test]
    fn test_create_multiple_users() {
        let mut db = Database::new().expect("Failed to create database.");

        let uid1: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let uid2: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let uid3: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();

        let name1 = CString::new(uid1).expect("CString::new failed.");
        let name2 = CString::new(uid2).expect("CString::new failed.");
        let name3 = CString::new(uid3).expect("CString::new failed.");

        let id1 = db.create_user(&name1).expect("Failed to create user.");
        let id2 = db.create_user(&name2).expect("Failed to create user.");
        let id3 = db.create_user(&name3).expect("Failed to create user.");

        assert_eq!(id1, 1);
        assert_eq!(id2, 2);
        assert_eq!(id3, 3);
    }

    #[test]
    fn test_delete_user_success() {
        let mut db = Database::new().expect("Failed to create database.");
        let uid1: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();
        let name = CString::new(uid1).expect("CString::new failed.");

        let user_id = db.create_user(&name).expect("Failed to create user.");

        let result = db.delete_user(user_id);
        assert!(result.is_ok());
    }

    #[test]
    fn test_delete_nonexistent_user() {
        let mut db = Database::new().expect("Failed to create database.");
        let result = db.delete_user(999);
        assert!(result.is_err());
        assert_eq!(result.err().unwrap(), Error::ErrUnknownId);
    }

    #[test]
    fn test_get_user_success() {
        let mut db = Database::new().expect("Failed to create database.");
        let name: String = thread_rng()
            .to_owned()
            .sample_iter(&Alphanumeric)
            .take(10)
            .map(char::from)
            .collect();
        let user_id = db
            .create_user(&CString::new(&*name).expect("CString::new() failed."))
            .expect("Failed to create user.");

        let user = db.get_user(user_id);
        assert!(user.is_ok());
        let user = user.unwrap();
        let user_name = unsafe { CStr::from_ptr(user.name) };
        assert_eq!(user.id, user_id);
        assert_eq!(user_name.to_str().unwrap(), name);
    }

    #[test]
    fn test_get_nonexistent_user() {
        let db = Database::new().expect("Failed to create database.");
        let result = db.get_user(999);
        assert!(result.is_err());
        assert_eq!(result.err().unwrap(), Error::ErrUnknownId);
    }

    #[test]
    fn test_database_drop() {
        let db = Database::new().expect("Failed to create database.");
        drop(db);
    }
}
