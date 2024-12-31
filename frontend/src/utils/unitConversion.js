// src/utils/unitConversion.js

export const kmToMph = (kmValue) => {
  const km = parseFloat(kmValue.toString().replace("km/h", ""));
  if (isNaN(km)) return kmValue;
  const mph = (km * 0.621371).toFixed(1);
  return `${mph} mph`;
};