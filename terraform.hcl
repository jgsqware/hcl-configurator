locals {
  cloudrun_services = {
    test-service = {
      features = {
        firebaseauth = "viewer"
        bucket_creator = ["tst",""]
        firestore_access = true
      }
    }
  }
}
