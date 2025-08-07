locals {
  cloudrun_services = {
    datalog-projection = {
      features = {
        firebaseauth = "admin"
      }
    }
  }
}
