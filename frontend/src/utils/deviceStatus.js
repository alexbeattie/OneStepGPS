// src/utils/deviceStatus.js

const formatTimeAgo = (timestamp) => {
  if (!timestamp) return 'unknown';

  try {
    const date = new Date(timestamp);
    const now = new Date();
    const seconds = Math.floor((now - date) / 1000);

    if (seconds < 60) return "just now";
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    return `${days}d ago`;
  } catch (error) {
    console.error('Error formatting time:', error);
    return 'invalid date';
  }
};

const getTimeElapsed = (timestamp) => {
  if (!timestamp) return 'unknown';

  try {
    const now = new Date();
    const startTime = new Date(timestamp);
    const diffInMs = now - startTime;

    const hours = Math.floor(diffInMs / (1000 * 60 * 60));
    const minutes = Math.floor((diffInMs % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((diffInMs % (1000 * 60)) / 1000);

    if (hours > 0) return `${hours}h ${minutes}m`;
    if (minutes > 0) return `${minutes}m ${seconds}s`;
    return `${seconds}s`;
  } catch (error) {
    console.error('Error calculating elapsed time:', error);
    return 'unknown';
  }
};

export const getDeviceStatus = (device) => {
  // Handle undefined device
  if (!device || !device.latest_device_point) {
    return {
      status: "Unknown",
      className: "bg-gray-100 text-gray-800 px-2 py-1 rounded-full text-xs font-medium",
      time: "unknown"
    };
  }

  const point = device.latest_device_point;
  const params = point.params || {};

  // Check if device is online based on last update time
  const lastUpdate = new Date(point.dt_tracker);
  const now = new Date();
  const minutesSinceUpdate = (now - lastUpdate) / (1000 * 60);
  const isOnline = minutesSinceUpdate < 5; // Consider offline if no update in 5 minutes

  if (!isOnline) {
    return {
      status: "Offline",
      className: "bg-red-100 text-red-800 px-2 py-1 rounded-full text-xs font-medium",
      time: formatTimeAgo(point.dt_tracker)
    };
  }

  // Check signal strength using gpslev parameter
  const signalStrength = parseInt(params.gpslev);
  if (!signalStrength || signalStrength < 4) { // GPS needs at least 4 satellites for accurate position
    return {
      status: "No Signal",
      className: "bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full text-xs font-medium",
      time: formatTimeAgo(point.dt_tracker)
    };
  }

  // Check if moving based on speed
  const speed = point.speed || 0;
  if (speed > 0) {
    return {
      status: "Moving",
      className: "bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs font-medium",
      time: getTimeElapsed(point.dt_tracker)
    };
  }

  // Device is stopped
  return {
    status: "Stopped",
    className: "bg-gray-100 text-gray-800 px-2 py-1 rounded-full text-xs font-medium",
    time: formatTimeAgo(point.dt_tracker)
  };
};