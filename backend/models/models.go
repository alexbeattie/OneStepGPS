package models
import (
	"time"
	"gorm.io/gorm"
)

// UserPreferences stores user-specific settings
type UserPreferences struct {
    gorm.Model
    UserID          string    `json:"user_id" gorm:"uniqueIndex"`
    SortOrder       string    `json:"sort_order"`
    HiddenDevices   []string  `json:"hidden_devices" gorm:"type:text[]"`
    DefaultFilters  string    `json:"default_filters"`
    MapSettings     string    `json:"map_settings"`
    ShowAddress     bool      `json:"show_address"`
    ShowEngineHours bool      `json:"show_engine_hours"`
    ShowOdometer    bool      `json:"show_odometer"`
    ShowVin         bool      `json:"show_vin"`
    ShowSpeed       bool      `json:"show_speed"`
    ShowHeading     bool      `json:"show_heading"`
    ShowBattery     bool      `json:"show_battery"`
    ShowSatellites  bool      `json:"show_satellites"`
    ShowLastUpdate  bool      `json:"show_last_update"`
    LastUpdated     time.Time `json:"last_updated"`
}



// APIResponse wraps the device list response from OneStepGPS API
type APIResponse struct {
	ResultList []Device `json:"result_list"`
}

// Device represents a GPS tracking device
type Device struct {
	DeviceID                  string                 `json:"device_id"`
	CreatedAt                 string                 `json:"created_at"`
	UpdatedAt                 string                 `json:"updated_at"`
	ActivatedAt               string                 `json:"activated_at"`
	DeliveredAt               interface{}            `json:"delivered_at"`
	FactoryID                 string                 `json:"factory_id"`
	ActiveState               string                 `json:"active_state"`
	DisplayName               string                 `json:"display_name"`
	BccID                     string                 `json:"bcc_id"`
	Make                      string                 `json:"make"`
	Model                     string                 `json:"model"`
	ConnType                  string                 `json:"conn_type"`
	ConnData                  map[string]interface{} `json:"conn_data"`
	DataNode                  string                 `json:"data_node"`
	Settings                  map[string]interface{} `json:"settings"`
	SecondaryID              string                 `json:"secondary_id"`
	UserIDList               []string               `json:"user_id_list"`
	Online                   bool                   `json:"online"`
	LatestDevicePoint        DevicePoint            `json:"latest_device_point"`
	LatestAccurateDevicePoint DevicePoint           `json:"latest_accurate_device_point"`
	DeviceGroupsIDList       interface{}            `json:"device_groups_id_list"`
	DeviceFieldList          interface{}            `json:"device_field_list"`
	DeviceUISettings         map[string]interface{} `json:"device_ui_settings"`
}

// DevicePoint represents a GPS location point with additional metadata
type DevicePoint struct {
	DevicePointID       string                 `json:"device_point_id"`
	DtServer           string                 `json:"dt_server"`
	DtTracker          string                 `json:"dt_tracker"`
	Lat                float64                `json:"lat"`
	Lng                float64                `json:"lng"`
	Altitude           interface{}            `json:"altitude"`
	Angle              int                    `json:"angle"`
	Speed              float64                `json:"speed"`
	Params             map[string]interface{} `json:"params"`
	DevicePointExternal map[string]interface{} `json:"device_point_external"`
	DevicePointDetail   DevicePointDetail      `json:"device_point_detail"`
	DeviceState         DeviceState            `json:"device_state"`
	DeviceStateStale    bool                   `json:"device_state_stale"`
	Sequence           string                 `json:"sequence"`
}

// LatLng represents a geographical coordinate
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// DevicePointDetail contains detailed information about a device's location point
type DevicePointDetail struct {
	FactoryID               string                 `json:"factory_id"`
	TransmitTime            string                 `json:"transmit_time"`
	GpsTime                 string                 `json:"gps_time"`
	Acc                     bool                   `json:"acc"`
	LatLng                  LatLng                 `json:"lat_lng"`
	Altitude                map[string]interface{} `json:"altitude"`
	Speed                   map[string]interface{} `json:"speed"`
	Heading                 int                    `json:"heading"`
	Hdop                    float64                `json:"hdop"`
	NumSatellites          int                    `json:"num_satellites"`
	RemoteAddr             string                 `json:"remote_addr"`
	HeventList             interface{}            `json:"hevent_list"`
	DtcList                interface{}            `json:"dtc_list"`
	MotionLog              map[string]interface{} `json:"motion_log"`
	PacketSequenceID       string                 `json:"packet_sequence_id"`
	Rssi                   float64                `json:"rssi"`
	TripDistance           map[string]interface{} `json:"trip_distance"`
	TravelDistance         map[string]interface{} `json:"travel_distance"`
	ExternalVolt           float64                `json:"external_volt"`
	BackupBatteryVolt      interface{}            `json:"backup_battery_volt"`
	VideoOriginalEvent     map[string]interface{} `json:"video_original_event"`
	RemainingBatteryPercent interface{}            `json:"remaining_battery_percent"`
	Historic               interface{}            `json:"historic"`
	SessionConnected       interface{}            `json:"session_connected"`
	DiffCorrected          interface{}            `json:"diff_corrected"`
	Predicted              interface{}            `json:"predicted"`
	Input1High             interface{}            `json:"input_1_high"`
	Input2High             interface{}            `json:"input_2_high"`
	VIN             string          `json:"vin" gorm:"column:detail_vin"`

}

// DeviceState represents the current state of a device
type DeviceState struct {
	DriveStatus               string                 `json:"drive_status"`
	DriveStatusID             string                 `json:"drive_status_id"`
	DriveStatusDuration       map[string]interface{} `json:"drive_status_duration"`
	DriveStatusDistance       map[string]interface{} `json:"drive_status_distance"`
	DriveStatusLatLngDistance map[string]interface{} `json:"drive_status_lat_lng_distance"`
	DriveStatusBeginTime      string                 `json:"drive_status_begin_time"`
	BestDistanceDelta         map[string]interface{} `json:"best_distance_delta"`
	IsNewDriveStatus          bool                   `json:"is_new_drive_status"`
	AdjustedLatLng            map[string]float64     `json:"adjusted_lat_lng"`
	BeyondMaxDriftDistance    bool                   `json:"beyond_max_drift_distance"`
	PrevDriveStatusDuration   map[string]interface{} `json:"prev_drive_status_duration"`
	PrevDriveStatusDistance   map[string]interface{} `json:"prev_drive_status_distance"`
	SoftwareOdometer  OdometerReading `json:"software_odometer" gorm:"embedded;prefix:software_"`
  HardwareOdometer  OdometerReading `json:"hardware_odometer" gorm:"embedded;prefix:hardware_"`
  Odometer         OdometerReading `json:"odometer" gorm:"embedded"`
	VIN             string          `json:"vin" gorm:"column:vin"`


}
type DevicePointExternal struct {
    SoftwareOdometerReading OdometerReading `json:"software_odometer_reading" gorm:"embedded;prefix:external_"`
}
type OdometerReading struct {
    Value    float64 `json:"value"`
    Unit     string  `json:"unit"`
    Display  string  `json:"display"`
}
// type OdometerResponse struct {
//     Data struct {
//         OdometerValue struct {
//             Value   float64 `json:"value"`
//             Unit    string  `json:"unit"`
//             Display string  `json:"display"`
//         } `json:"odometer_value"`
//     } `json:"data"`
// }
type DurationData struct {
    Value   float64 `json:"value"`
    Unit    string  `json:"unit"`
    Display string  `json:"display"`
}

// type LatLng struct {
//     Lat float64 `json:"lat"`
//     Lng float64 `json:"lng"`
// }

type DriveStopPoint struct {
    Type         string       `json:"type"`
    Duration     DurationData `json:"duration"`
    FirstLatLng  LatLng      `json:"first_valid_lat_lng"`
    LastLatLng   LatLng      `json:"last_valid_lat_lng"`
    TimeFrom     string      `json:"time_from"`
    TimeTo       string      `json:"time_to"`
    OdometerFrom struct {
        Value   float64 `json:"value"`
        Unit    string  `json:"unit"`
        Display string  `json:"display"`
    } `json:"odometer_from"`
    OdometerTo struct {
        Value   float64 `json:"value"`
        Unit    string  `json:"unit"`
        Display string  `json:"display"`
    } `json:"odometer_to"`
}

type DriveStopResponse struct {
    TimeFrom        string       `json:"time_from"`
    TimeTo          string       `json:"time_to"`
    Duration        DurationData `json:"duration"`
    Distance        DurationData `json:"distance"`
    AverageSpeed    DurationData `json:"average_speed"`
    IdleDuration    DurationData `json:"idle_duration"`
    StopDuration    DurationData `json:"stop_duration"`
    TopSpeed        DurationData `json:"top_speed"`
    DriveStopList   []DriveStopPoint `json:"drive_stop_list"`
}
type Measurement struct {
    Value   float64 `json:"value"`
    Unit    string  `json:"unit"`
    Display string  `json:"display"`
}
type DeviceInfoResponse struct {
    ResultList []DeviceInfo `json:"result_list"`
}

type DeviceInfo struct {
    DeviceID    string `json:"device_id"`
    DisplayName string `json:"display_name"`
    // Add other fields that you need from the device info response
}


// type InfoWindowPreferences struct {
//     gorm.Model
//     UserID          string `json:"user_id" gorm:"uniqueIndex"`
//     ShowAddress     bool   `json:"show_address" gorm:"default:true"`
//     ShowEngineHours bool   `json:"show_engine_hours" gorm:"default:true"`
//     ShowOdometer    bool   `json:"show_odometer" gorm:"default:true"`
//     ShowVIN         bool   `json:"show_vin" gorm:"default:true"`
//     ShowSpeed       bool   `json:"show_speed" gorm:"default:true"`
//     ShowHeading     bool   `json:"show_heading" gorm:"default:true"`
//     ShowBattery     bool   `json:"show_battery" gorm:"default:true"`
//     ShowSatellites  bool   `json:"show_satellites" gorm:"default:true"`
//     ShowLastUpdate  bool   `json:"show_last_update" gorm:"default:true"`
// }

type DriveStop struct {
    Type                string      `json:"type"`
    Duration           Measurement  `json:"duration"`
    Distance           Measurement  `json:"distance"`
    AverageSpeed       Measurement  `json:"average_speed"`
    IdleDuration       Measurement  `json:"idle_duration"`
    TopSpeed           Measurement  `json:"top_speed"`
    FirstValidLatLng   LatLng      `json:"first_valid_lat_lng"`
    LastValidLatLng    LatLng      `json:"last_valid_lat_lng"`
    TimeFrom           string      `json:"time_from"`
    TimeTo             string      `json:"time_to"`
    OdometerFrom       Measurement `json:"odometer_from"`
    OdometerTo         Measurement `json:"odometer_to"`
}