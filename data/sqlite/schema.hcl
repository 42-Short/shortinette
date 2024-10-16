schema "main" {}

table "users" {
    schema = schema.main
    column "id" {
        null = false
        type = text
    }
    primary_key {
        columns = [column.id]
    }
}

table "mod_00" {
    schema = schema.main 
    column "id" {
        null = false 
        type = text
    }
    column "score" {
        null = false 
        type = integer
        default = 0
    }
    column "wait_time" {
        null = false 
        type = integer
        default = 0
    }
    primary_key {
        columns = [column.id]
    }
    foreign_key "fk1" {
        columns = [column.id]
        ref_columns = [table.users.column.id]
    }
}