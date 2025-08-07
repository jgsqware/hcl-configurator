locals {
  cloudrun_services = {
    datalog-business = {
      features = {
        firestore_access = true
        bucket_reader = ["datalog-dmz","datalogs"]
      }
    }
    mobile-interface = {
      features = {
        firebaseauth = "viewer"
        cloudrun_invoker = true
        bucket_creator = ["datalog-dmz"]
        firestore_access = true
        bucket_reader = ["datalog-dmz","datalogs"]
      }
    }
  }
}
