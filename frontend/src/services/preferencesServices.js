export async function loadPreferences(deviceId, apiBaseUrl) {
  if (!deviceId) throw new Error("Device ID is required.");
  const url = `${apiBaseUrl}/api/v1/preferences/${deviceId}`;
  const response = await fetch(url);
  if (!response.ok) throw new Error(`Failed to load preferences: ${response.status}`);
  return response.json();
}

export async function savePreferences(deviceId, preferences, apiBaseUrl) {
  if (!deviceId) throw new Error("Device ID is required.");
  const url = `${apiBaseUrl}/api/v1/preferences/${deviceId}`;
  const response = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify(preferences),
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to save preferences: ${response.status} - ${errorText}`);
  }
  return response.json();
}
