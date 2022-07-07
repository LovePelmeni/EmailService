#!/bin/sh 
## there goes method that initialize database and creates customer for access.
EOF << 
    use emails_db 
    db.SiblingsDB("emails_db")
    use emails_db 
    db.CreateUser({
        "username": "mongo_user",
        "password": "mongo_password"
    })
    db.Auth(mongo_user, mongo_password)
>> EOF 