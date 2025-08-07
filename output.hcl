locals {
  cloudrun_services = {
    datalog-atcloud-projector = {
      features = {
        cloudrun_invoker = true
        firestore_access = true
        bucket_writer = ["datalogs","files"]
        mysql_access = true
        postgres_access = true
      }
      mysql_grants = {
        read = [
          "autotuner_v2.model_year_engine",
          "autotuner_v2.model_year",
          "autotuner_v2.engine",
          "autotuner_v2.model",
          "autotuner_v2.manufacturer",
        ]
      }
      postgres_grants = {
        read = [
          "_storage.bucketz",
          "_thread.participant",
          "_thread.thread",
          "_user.firebase_reference",
          "_user.user",
          "_vehicle.owner",
          "_vehicle.work_sheet",
        ]
        read_write = [
          "_storage.file",
          "_thread.message",
          "_thread.message_attachment",
          "_thread.message_attachment_properties",
        ]
        read_write_delete = [
          "public.message_offset",
          "public.message_handler_error",
        ]
      }
    }
    datalog-business = {
      features = {
        firestore_access = true
        bucket_reader = ["datalog-dmz","datalogs"]
      }
    }
    datalog-exporter-projector = {
      features = {
        firestore_access = true
        bucket_reader = ["datalog-dmz"]
        bucket_writer = ["datalogs"]
        mysql_access = true
      }
      mysql_grants = {
        read = [
          "autotuner_v2.datalog_family",
          "autotuner_v2.datalog_measure",
          "autotuner_v2.datalog_measure_supported_ecu",
          "autotuner_v2.ecu",
          "autotuner_v2.engine",
          "autotuner_v2.engine_ecu",
          "autotuner_v2.manufacturer",
          "autotuner_v2.mcu.id",
          "autotuner_v2.mcu.name",
          "autotuner_v2.mcu.public_name",
          "autotuner_v2.model",
          "autotuner_v2.model_year",
          "autotuner_v2.model_year_engine",
          "autotuner_v2.supported_ecu.id",
          "autotuner_v2.supported_ecu.name",
          "autotuner_v2.supported_ecu.mcu_id",
          "autotuner_v2.supported_ecu.method_id",
          "autotuner_v2.supported_ecu_olsx_ecu",
          "autotuner_v2.translatable",
          "autotuner_v2.datalog_generic_excluded_ecu.ecu_id",
        ]
      }
    }
    mobile-interface = {
      features = {
        firebaseauth = "viewer"
        cloudrun_invoker = true
        bucket_creator = ["datalog-dmz"]
      }
    }
  }
}
