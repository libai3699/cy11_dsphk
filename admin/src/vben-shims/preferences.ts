// @vben/preferences 替代实现

export interface PreferencesConfig {
  namespace?: string;
  overrides?: any;
  extension?: any;
}

export async function initPreferences(config: PreferencesConfig) {
  console.log('Preferences initialized:', config.namespace);
  return Promise.resolve();
}

export function defineOverridesPreferences(overrides: any) {
  return overrides;
}

export function definePreferencesExtension(extension: any) {
  return extension;
}
