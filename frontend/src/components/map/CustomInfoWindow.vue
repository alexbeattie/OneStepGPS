<template>
  <div>
    <!-- Add root div -->
    <!-- Drive Stop Modal -->
    <DriveStopModal
      v-if="showRouteModal"
      :device-id="device.device_id"
      @close="handleModalClose"
    />

    <!-- Add preferences modal -->
    <div
      v-if="showPreferencesModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full">
        <!-- Modal Header -->
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold">Info Window Preferences</h3>
          <button
            @click="togglePreferencesModal"
            class="text-gray-500 hover:text-gray-700"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div class="space-y-4">
          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_address"
              id="show_address"
              class="mr-2"
            />
            <label for="show_address">Show Address</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_engine_hours"
              id="show_engine_hours"
              class="mr-2"
            />
            <label for="show_engine_hours">Show Engine Hours</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_odometer"
              id="show_odometer"
              class="mr-2"
            />
            <label for="show_odometer">Show Odometer</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_vin"
              id="show_vin"
              class="mr-2"
            />
            <label for="show_vin">Show VIN</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_speed"
              id="show_speed"
              class="mr-2"
            />
            <label for="show_speed">Show Speed</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_heading"
              id="show_heading"
              class="mr-2"
            />
            <label for="show_heading">Show Heading</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_battery"
              id="show_battery"
              class="mr-2"
            />
            <label for="show_battery">Show Battery</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_satellites"
              id="show_satellites"
              class="mr-2"
            />
            <label for="show_satellites">Show Satellites</label>
          </div>

          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="preferences.show_last_update"
              id="show_last_update"
              class="mr-2"
            />
            <label for="show_last_update">Show Last Update</label>
          </div>
        </div>

        <div class="mt-6 flex justify-end">
          <div v-if="saveError" class="text-red-500 mt-2 text-sm">
            {{ saveError }}
          </div>
          <button
            @click="savePreferences"
            :disabled="isSaving"
            class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:bg-blue-300"
          >
            <span v-if="isSaving" class="spinner-border"></span>

            {{ isSaving ? "Saving..." : "Save Preferences" }}
          </button>
        </div>
      </div>
    </div>
    <!-- Mobile-friendly preferences modal -->
    <div
      v-if="showPreferencesModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
    >
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold">Preferences</h3>
          <button
            @click="togglePreferencesModal"
            class="text-gray-500 hover:text-gray-700"
          >
            <span class="sr-only">Close</span>
            ✕
          </button>
        </div>

        <div class="space-y-3">
          <label
            v-for="(value, key) in preferences"
            :key="key"
            class="flex items-center justify-between p-2 hover:bg-gray-50 rounded"
          >
            <span class="text-sm text-gray-600">
              {{ formatPreferenceLabel(key) }}
            </span>
            <input
              type="checkbox"
              v-model="preferences[key]"
              class="rounded text-blue-600"
            />
          </label>
        </div>

        <div class="mt-6 flex justify-end space-x-3">
          <button
            @click="togglePreferencesModal"
            class="px-4 py-2 text-sm text-gray-600 hover:text-gray-800"
          >
            Cancel
          </button>
          <button
            @click="savePreferences"
            :disabled="isSaving"
            class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:bg-blue-300"
          >
            {{ isSaving ? "Saving..." : "Save Changes" }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { getBatteryStatus } from "@/services/deviceMetrics";
import DriveStopModal from "@/components/device/DriveStopModal.vue";

export default {
  name: "CustomInfoWindow",
  emits: ["closeclick"], // Add this line

  props: {
    device: {
      type: Object,
      required: true,
    },
    position: {
      type: Object,
      required: true,
    },
  },
  inject: [
    "getDeviceStatus",
    "getAddressForDevice",
    "kmToMph",
    "getDeviceOdometer",
    "getDeviceEngineHours",
  ],
  data() {
    return {
      showRouteModal: false,
      infoWindow: null,
      showPreferencesModal: false,
      isSaving: false,
      saveError: null,

      preferences: {
        show_address: true,
        show_engine_hours: true,
        show_odometer: true,
        show_vin: true,
        show_speed: true,
        show_heading: true,
        show_battery: true,
        show_satellites: true,
        show_last_update: true,
      },
    };
  },

  components: {
    DriveStopModal,
  },
  computed: {
    // Add this computed property
    deviceEngineHours() {
      return this.getDeviceEngineHours(this.device.device_id);
    },
  },
  mounted() {
    console.log("API Base URL:", process.env.VUE_APP_API_BASE_URL);
    this.loadPreferences()
      .then(() => {
        // Only initialize after preferences are loaded
        this.initializeInfoWindow();
      })
      .catch((error) => {
        console.error("Error loading preferences:", error);
        // Initialize with defaults if loading fails
        this.initializeInfoWindow();
      });
  },
  watch: {
    // Add this watcher to handle position updates

    position: {
      handler(newPos) {
        if (this.infoWindow) {
          this.infoWindow.setPosition(newPos);
        }
      },
      deep: true,
    },

    showPreferencesModal: {
      handler(newVal) {
        console.log("Preferences modal visibility:", newVal);
        if (newVal) {
          this.infoWindow?.close();
        }
      },
    },
    "device.device_id": {
      handler(newDeviceId) {
        if (newDeviceId) {
          this.loadPreferences();
        }
      },
      immediate: true,
    },
    showRouteModal: {
      handler(newVal) {
        console.log("Route modal visibility:", newVal);
        if (newVal) {
          this.infoWindow?.close();
        }
      },
    },
  },
  methods: {
    formatPreferenceLabel(key) {
      return key
        .replace("show_", "")
        .replace(/_/g, " ")
        .replace(/\b\w/g, (l) => l.toUpperCase());
    },
    togglePreferencesModal() {
      if (this.infoWindow) {
        if (!this.showPreferencesModal) {
          // When opening preferences modal
          this.infoWindow.close();
          this.showPreferencesModal = true;
        } else {
          // When closing preferences modal
          this.showPreferencesModal = false;
          if (this.$parent?.map) {
            this.infoWindow.open(this.$parent.map);
          }
        }
      } else {
        // Toggle modal even if infoWindow doesn't exist
        this.showPreferencesModal = !this.showPreferencesModal;
      }
      console.log("Preferences modal state:", this.showPreferencesModal);
    },

    handleModalClose() {
      this.showRouteModal = false;
    },

    handleRouteButtonClick(e) {
      e.stopPropagation();
      this.showRouteModal = true;
      console.log("Route button clicked, showRouteModal:", this.showRouteModal); // Debug log
    },
    async initializeInfoWindow() {
      try {
        // await this.loadPreferences();

        if (!window.google?.maps) {
          console.warn("Google Maps not yet loaded");
          return;
        }

        const content = document.createElement("div");
        content.className = "bg-white rounded-lg shadow-lg max-w-md";
        content.style.paddingTop = "0"; // Removes top padding
        content.style.overflow = "hidden"; // Prevents scrolling
        content.style.maxHeight = "none"; // Ensures height is unrestricted
        content.className =
          "bg-white rounded-lg shadow-lg max-w-md custom-info-window";
        content.style.paddingTop = "0"; // Removes top padding
        // Add error handling for all the data access

        const deviceName = this.device?.display_name || "Unknown Device";
        const deviceStatus = this.getDeviceStatus(this.device) || {
          className: "",
          status: "Unknown",
        };
        const address =
          this.getAddressForDevice(this.device?.device_id) || "Loading...";
        const speed = this.device?.latest_device_point?.speed
          ? this.kmToMph(this.device.latest_device_point.speed)
          : "N/A";
        const angle = this.device?.latest_device_point?.angle || "N/A";
        const satellites =
          this.device?.latest_device_point?.params?.gpslev || "N/A";
        const batteryStatus = this.getBatteryStatusHTML();
        const deviceEngineHours = this.getDeviceEngineHours(
          this.device.device_id
        );
        const engineHours = deviceEngineHours
          ? `${Number(deviceEngineHours).toFixed(2)}`
          : "N/A";

        let contentHTML = `
      <div class="relative px-4 py-2 mt-5 border-b border-gray-200" style="padding-top: 0;">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-semibold text-gray-900">${deviceName}</h3>
          <span class="inline-block px-2.5 py-0.5 text-xs font-medium rounded-full ml-2 ${deviceStatus.className}">
            ${deviceStatus.status}
          </span>
        </div>
      </div>
      <div class="px-4 py-2">
        <div class="space-y-1.5">`;

        // Conditionally add elements based on preferences
        if (this.preferences.show_address) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Address:</div>
          <div class="text-gray-900 truncate">${address}</div>
        </div>`;
        }

        if (this.preferences.show_last_update) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Last Seen:</div>
          <div class="text-gray-900">${deviceStatus.time || "Unknown"}</div>
        </div>`;
        }

        if (this.preferences.show_engine_hours) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Engine Hours:</div>
          <div class="text-gray-900">${engineHours}</div>
        </div>`;
        }

        if (this.preferences.show_odometer) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Odometer:</div>
          <div class="text-gray-900">${this.getOdometerInMiles(
            this.device
          )}</div>
        </div>`;
        }

        // Add the rest of the conditional elements
        if (this.preferences.show_vin) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">VIN:</div>
          <div class="text-gray-900">${this.getVehicleVIN(this.device)}</div>
        </div>`;
        }

        if (this.preferences.show_speed) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Speed:</div>
          <div class="text-gray-900">${speed}</div>
        </div>`;
        }

        if (this.preferences.show_heading) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Heading:</div>
          <div class="text-gray-900">${angle}°</div>
        </div>`;
        }

        if (this.preferences.show_battery) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Battery:</div>
          <div class="text-gray-900">${batteryStatus}</div>
        </div>`;
        }

        if (this.preferences.show_satellites) {
          contentHTML += `
        <div class="flex text-sm">
          <div class="font-medium text-gray-700 min-w-24">Satellites:</div>
          <div class="text-gray-900">${satellites}</div>
        </div>`;
        }

        // Close the divs and add the route button
        contentHTML += `
   <div class="flex text-sm mt-2">
    <div class="px-4 py-2 mt-2 mb-2 flex space-x-2 w-full">
      <button
        id="showRouteBtn-${this.device.device_id}"
        type="button"
        onclick="void(0)"
        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors flex-1"
      >
        Show Route History
      </button>
      <button
        id="prefsBtn-${this.device.device_id}"
        type="button"
        onclick="void(0)"
        class="bg-gray-500 text-white px-3 py-2 rounded hover:bg-gray-600 transition-colors flex-shrink-0 z-50"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
        </svg>
      </button>
    </div>
  </div>`;

        content.innerHTML = contentHTML;

        // Create and configure InfoWindow
        this.infoWindow = new window.google.maps.InfoWindow({
          content,
          position: this.position,
          pixelOffset: new window.google.maps.Size(0, -10),
          disableAutoPan: true,
          clickable: true,
        });
        // Single domready listener for both buttons
        this.infoWindow.addListener("domready", () => {
          // Setup route button
          const routeButton = document.getElementById(
            `showRouteBtn-${this.device.device_id}`
          );
          if (routeButton) {
            routeButton.onclick = (e) => {
              // Use onclick instead of addEventListener
              console.log("Route button clicked");
              e.preventDefault();
              e.stopPropagation();
              this.showRouteModal = true;
            };
          } else {
            console.error("Route button not found");
          }

          // Setup preferences button
          // Setup preferences button
          // In initializeInfoWindow, update the preferences button handler
          // Setup preferences button
          const prefsButton = document.getElementById(
            `prefsBtn-${this.device.device_id}`
          );
          if (prefsButton) {
            console.log("Found preferences button");
            prefsButton.onclick = (e) => {
              // Use onclick instead of addEventListener
              console.log("Preferences button clicked");
              e.preventDefault();
              e.stopPropagation();
              this.showPreferencesModal = true;
            };
          } else {
            console.error(
              "Preferences button not found:",
              `prefsBtn-${this.device.device_id}`
            );
          }
        });
        // Add close listener
        this.infoWindow.addListener("closeclick", () => {
          this.showRouteModal = false;
          this.$emit("closeclick");
        });

        // Open the InfoWindow
        if (this.$parent?.map) {
          this.infoWindow.open(this.$parent.map);
        } else {
          console.error("Parent map not found");
        }
      } catch (error) {
        console.error("Error initializing info window:", error);
      }
      // this.infoWindow.addListener("closeclick", () => {
      //   this.showRouteModal = false;
      //   this.$emit("closeclick");
      // });

      if (this.$parent?.map) {
        this.infoWindow.open(this.$parent.map);
      } else {
        console.error("Parent map not found");
      }
    },
    async loadPreferences() {
      try {
        const deviceId = this.device?.device_id; // Use device ID instead of getUserId
        if (!deviceId) {
          console.warn("No device ID available");
          return;
        }
        const baseUrl = process.env.VUE_APP_API_BASE_URL;
        const url = `${baseUrl}/api/v1/preferences/${deviceId}`;

        console.log("Loading preferences from:", url);

        const response = await fetch(url);

        if (!response.ok) {
          throw new Error(`Failed to load preferences: ${response.status}`);
        }

        const data = await response.json();
        console.log("Successfully loaded preferences from DB:", data);

        this.preferences = {
          show_address: data.show_address ?? true,
          show_engine_hours: data.show_engine_hours ?? true,
          show_odometer: data.show_odometer ?? true,
          show_vin: data.show_vin ?? true,
          show_speed: data.show_speed ?? true,
          show_heading: data.show_heading ?? true,
          show_battery: data.show_battery ?? true,
          show_satellites: data.show_satellites ?? true,
          show_last_update: data.show_last_update ?? true,
        };
      } catch (error) {
        console.error("Error loading preferences:", error);
        // Keep current preferences instead of resetting to defaults if load fails
        if (!this.preferences) {
          this.preferences = {
            show_address: true,
            show_engine_hours: true,
            show_odometer: true,
            show_vin: true,
            show_speed: true,
            show_heading: true,
            show_battery: true,
            show_satellites: true,
            show_last_update: true,
          };
        }
      }
    },

    async savePreferences() {
      this.isSaving = true;
      this.saveError = null;

      try {
        const deviceId = this.device?.device_id;
        if (!deviceId) {
          throw new Error("No device ID available");
        }

        const baseUrl =
          process.env.VUE_APP_API_BASE_URL || "http://localhost:8080";

        const preferencesPayload = {
          user_id: deviceId,
          sort_order: "",
          hidden_devices: [],
          default_filters: "",
          map_settings: "",
          show_address: this.preferences.show_address,
          show_engine_hours: this.preferences.show_engine_hours,
          show_odometer: this.preferences.show_odometer,
          show_vin: this.preferences.show_vin,
          show_speed: this.preferences.show_speed,
          show_heading: this.preferences.show_heading,
          show_battery: this.preferences.show_battery,
          show_satellites: this.preferences.show_satellites,
          show_last_update: this.preferences.show_last_update,
          last_updated: new Date().toISOString(),
        };

        console.log("Saving preferences:", preferencesPayload);

        const response = await fetch(
          `${baseUrl}/api/v1/preferences/${deviceId}`,
          {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            body: JSON.stringify(preferencesPayload),
          }
        );

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(
            `Failed to save preferences: ${response.status} - ${errorText}`
          );
        }

        const data = await response.json();
        console.log("Save successful:", data);

        // Close modal first
        this.showPreferencesModal = false;

        // Then re-initialize info window with current preferences
        this.$nextTick(() => {
          this.initializeInfoWindow();
        });
      } catch (error) {
        console.error("Error saving preferences:", error);
        this.saveError = error.message;
      } finally {
        this.isSaving = false;
      }
    },

    getBatteryStatusHTML() {
      const voltage =
        this.device.latest_device_point?.params?.obd_battery_voltage;
      if (!voltage) return "N/A";

      const status = getBatteryStatus(voltage); // Using the imported function
      let colorClass;

      switch (status) {
        case "critical":
          colorClass = "text-red-600";
          break;
        case "low":
          colorClass = "text-orange-500";
          break;
        case "high":
          colorClass = "text-yellow-500";
          break;
        case "normal":
          colorClass = "text-green-500";
          break;
        default:
          colorClass = "text-gray-500";
      }

      return `<span class="${colorClass}">${voltage}V (${status})</span>`;
    },
    getVehicleVIN(device) {
      // Check all possible VIN locations
      let vin = null;
      // Try device_state in latest_device_point
      vin = device?.latest_device_point?.device_state?.vin;
      // If not found, try device_point_detail
      if (!vin) {
        vin = device?.latest_device_point?.device_point_detail?.vin;
      }
      // If still not found, try direct device property
      if (!vin) {
        vin = device?.vin;
      }
      return vin || "N/A";
    },

    formatOdometer(miles) {
      return miles ? `${Number(miles).toLocaleString()} mi` : "N/A";
    },

    formatEngineHours(hours) {
      return hours ? Number(hours).toFixed(2) : "N/A";
    },

    getOdometerInMiles(device) {
      // Check all possible odometer locations
      let odometerKm = null;
      // Try device_state.odometer first
      odometerKm = device?.latest_device_point?.device_state?.odometer?.value;
      // If not found, try software_odometer
      if (!odometerKm) {
        odometerKm =
          device?.latest_device_point?.device_state?.software_odometer?.value;
      }
      // Try device_point_external
      if (!odometerKm) {
        odometerKm =
          device?.latest_device_point?.device_point_external
            ?.software_odometer_reading?.value;
      }
      // Try hardware odometer if available
      if (!odometerKm) {
        odometerKm =
          device?.latest_device_point?.device_state?.hardware_odometer?.value;
      }
      if (!odometerKm) return "N/A";
      // Convert km to miles (1 km = 0.621371 miles)
      const miles = odometerKm * 0.621371;
      // Format to 1 decimal place
      return `${Number(miles).toLocaleString(undefined, {
        minimumFractionDigits: 1,
        maximumFractionDigits: 1,
      })} mi`;
    },
  },

  beforeUnmount() {
    if (this.infoWindow) {
      this.infoWindow.close();
    }
  },
};
</script>
<style>
@media (max-width: 640px) {
  .gm-style-iw {
    max-width: 90vw !important;
  }

  .gm-style-iw-d {
    max-width: 100% !important;
    overflow: auto !important;
  }
}

.gm-style-iw-chr {
  display: none !important;
}

.gm-style-iw-tc,
.gm-style-iw-t::after,
.gm-style-iw-tip {
  display: block !important;
  visibility: visible !important;
  opacity: 1 !important;
}
</style>
