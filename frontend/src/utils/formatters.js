export const formatTimeAgo = (timestamp) => {
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
};

export const getTimeElapsed = (timestamp) => {
  const now = new Date();
  const startTime = new Date(timestamp);
  const diffInMs = now - startTime;

  const hours = Math.floor(diffInMs / (1000 * 60 * 60));
  const minutes = Math.floor((diffInMs % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((diffInMs % (1000 * 60)) / 1000);

  if (hours > 0) return `${hours}h ${minutes}m`;
  if (minutes > 0) return `${minutes}m ${seconds}s`;
  return `${seconds}s`;
};

export const formatDuration = (durationDisplay) => {
  const hours = parseInt(durationDisplay);
  if (isNaN(hours)) return durationDisplay;

  const days = Math.floor(hours / 24);
  const remainingHours = hours % 24;

  if (days === 0) return `${hours} hours`;
  if (remainingHours === 0) return `${days} days`;
  return `${days} days, ${remainingHours} hours`;
};