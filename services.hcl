locals {
  cloudrun_services = {
    datalog-projection = {
      features = {
        firebaseauth = "admin"
        bucket_writer = ["datalog-dmz"]
      }
    }
  }
}
