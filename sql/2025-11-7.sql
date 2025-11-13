user := models.Users{
    Username: "kube",
    Password: "1qazWSX!",
}
database.DB.Create(&user)