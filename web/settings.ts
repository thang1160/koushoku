export enum Mode {
  Normal,
  Strip
}

export interface Settings {
  mode?: Mode;
  maxWidth?: number;
  zoomLevel?: number;
}

export const settings: Settings = {
  mode: 0,
  maxWidth: 1366,
  zoomLevel: 1.0
};

export const saveSettings = () => {
  localStorage.setItem("settings", JSON.stringify(settings));
};

export const getSettings = () => {
  const localSettings = localStorage.getItem("settings");
  if (localSettings) {
    const obj = JSON.parse(localSettings);
    Object.assign(settings, obj);
  } else saveSettings();
};

export const setSetting = (key: string, value: any) => {
  if (!(key in settings) || (key === "zoomLevel" && (value < 0.1 || value > 5.0))) {
    return;
  }

  settings[key] = value;
  saveSettings();
};
