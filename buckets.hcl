locals {
  buckets = {
    datalog-dmz = {}
    datalogs = {}
    attachments = {
      existing_name = "event_message_attachments"
      storage_class = "REGIONAL"
    }
    avatars = {
      existing_name = "user_avatars-e21tn"
      storage_class = "REGIONAL"
      delete_after_days = 7
    }
    device-brandings = {
      existing_name = "device_brandings-dev"
    }
    device-brandings-autotuner = {
      existing_name = "device_brandings_autotuner-dev"
    }
    device-logs = {
      existing_name = "device_logs-dev"
    }
    files = {
      existing_name = "event_file_attachments"
      storage_class = "REGIONAL"
    }
    support = {
      existing_name = "support_attachments"
      storage_class = "REGIONAL"
    }
    vehicles = {
      existing_name = "vehicles-media"
    }
    toolbox = {
      existing_name = "toolbox-files-dev"
      delete_after_days = 1
    }
    dataflow-topic-replicate-uat = {} # TESTING
    dataflow-topic-replica-prod = {}
    at-one-firmware = {}
  }
}