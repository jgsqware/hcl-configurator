locals {
  cloudrun_services = {
    ffddsfdsfs = {
      features = {
        cloudrun_invoker = true
        bucket_creator = ["test"]
        firestore_access = true
        mysql_access = true
        postgres_access = true
      }
    }
  }
}
