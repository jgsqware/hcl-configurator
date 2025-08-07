locals {
  cloudrun_services = {
    datalog-projection = {
      features = {
        firebaseauth = "view"
        mysql_access = true
      }
    }
  }
}
