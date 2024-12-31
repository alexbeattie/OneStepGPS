// src/services/api.js

class ApiService {
  constructor() {
    this.baseUrl = process.env.VUE_APP_API_BASE_URL || 'https://onestepgpsdemo.com';
    this.headers = {
      'Content-Type': 'application/json',
    };
  }
  async handleResponse(response) {
    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || `HTTP error! status: ${response.status}`);
    }
    return response.json();
  }

  // Get all devices
  async getDevices() {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/devices`, {
        method: 'GET',
        headers: this.headers
      });
      return this.handleResponse(response);
    } catch (error) {
      console.error('Error fetching devices:', error);
      throw error;
    }
  }

  //get Odometer data
  // Add this to your ApiService class
  async getDeviceOdometer(deviceId) {
    try {
      const response = await fetch(`${this.baseUrl}/v3/api/odometer/${deviceId}`, {
        method: 'GET',
        headers: this.headers
      });
      return this.handleResponse(response);
    } catch (error) {
      console.error(`Error fetching device ${deviceId} odometer:`, error);
      throw error;
    }
  }

  // src/services/api.js
  // Add this method to your ApiService class
  async getDriveStopRoute(deviceId, fromDate, toDate, stopDuration = '5m0s') {
    try {
      const fromStr = new Date(fromDate).toISOString();
      const toStr = new Date(toDate).toISOString();

      const response = await fetch(
        `${this.baseUrl}/v3/api/route/drive-stop?` +
        `device_id=${deviceId}&` +
        `dt_tracker_from=${fromStr}&` +
        `dt_tracker_to=${toStr}&` +
        `stop_duration=${stopDuration}&` +
        `return_points=true&` +
        `max_return_points=999`,
        {
          method: 'GET',
          headers: this.headers
        }
      );
      return this.handleResponse(response);
    } catch (error) {
      console.error('Error fetching drive-stop route:', error);
      throw error;
    }
  }
  // Get single device
  async getDevice(deviceId) {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/devices/${deviceId}`, {
        method: 'GET',
        headers: this.headers
      });
      return this.handleResponse(response);
    } catch (error) {
      console.error(`Error fetching device ${deviceId}:`, error);
      throw error;
    }
  }


  // Update device settings
  async updateDeviceSettings(deviceId, settings) {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/devices/${deviceId}/settings`, {
        method: 'PUT',
        headers: this.headers,
        body: JSON.stringify(settings)
      });
      return this.handleResponse(response);
    } catch (error) {
      console.error(`Error updating device ${deviceId} settings:`, error);
      throw error;
    }
  }

  async getDeviceHistory(deviceId, params = {}) {
    const queryParams = new URLSearchParams(params).toString();
    try {
      const response = await fetch(
        `${this.baseUrl}/api/v1/devices/${deviceId}/history?${queryParams}`,
        {
          method: 'GET',
          headers: this.headers
        }
      );
      return this.handleResponse(response);
    } catch (error) {
      console.error(`Error fetching device ${deviceId} history:`, error);
      throw error;
    }
  }

  async getDeviceAlerts(deviceId, params = {}) {
    const queryParams = new URLSearchParams(params).toString();
    try {
      const response = await fetch(
        `${this.baseUrl}/api/v1/devices/${deviceId}/alerts?${queryParams}`,
        {
          method: 'GET',
          headers: this.headers
        }
      );
      return this.handleResponse(response);
    } catch (error) {
      console.error(`Error fetching device ${deviceId} alerts:`, error);
      throw error;
    }
  }

  async updateDeviceLocation(deviceId, location) {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/devices/${deviceId}/location`, {
        method: 'PUT',
        headers: this.headers,
        body: JSON.stringify(location)
      });
      return this.handleResponse(response);
    } catch (error) {
      console.error(`Error updating device ${deviceId} location:`, error);
      throw error;
    }
  }
}

export const apiService = new ApiService();
