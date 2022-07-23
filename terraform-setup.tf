variable GOOGLE_CLOUD_PROJECT_ID {
    type = "string"
    description = "Google Cloud Project Id."
}
variable GOOGLE_CLOUD_REGION {
    type = "string"
    description = "Google Cloud Region the Project Is Running on."
}

variable GOOGLE_CLOUD_PROJECT_ZONE {
    type = "string"
    description = "Google Cloud Project Time Zone."
}

variable MACHINE_TYPE{
    type = "string" 
    description = "Machine Type resource instance is going to be running on For Example `f1-micro` on Google Cloud"
}

variable OPERATIONAL_DISK_CLOUD_SYSTEM {
    type = "string"
    description = "Operational System the Google Cloud Boot Disk Is Running On."
}

provider "google" {

    version = "3.5.0"
    credentials = "google-cloud.json"

    project_id = "${GOOGLE_CLOUD_PROJECT_ID}" 
    region = "${GOOGLE_CLOUD_REGION}"
    zone = "${GOOGLE_CLOUD_PROJECT_ZONE}"
    request_timeout = 60
}

resource "google_compute_instance" "appserver"{

    name = "store-application-server"
    machine_type = "${MACHINE_TYPE}"

    boot_disk{  
        initialize_params{
            image = "${OPERATIONAL_CLOUD_SYSTEM}" // Operational System that is currently 
            // Used in the Cloud Project for managing... 
        }
    }
}