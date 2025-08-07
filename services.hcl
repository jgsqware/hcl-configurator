locals {
  cloudrun_services = {
    datalog-projection = {
      features = {
        mysql_access = true
      }
    }
  }
}
