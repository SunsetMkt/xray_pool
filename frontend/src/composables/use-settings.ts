import { computed, onMounted, reactive } from 'vue';
import type { ApiResponseSettings } from '@/interfaces/settings';
import SettingsApi from '@/apis/SettingsApi';

export type SettingsMode = 'normal' | 'pro';

export interface SettingsState {
  settings: ApiResponseSettings | null;
  model: ApiResponseSettings | null;
  mode: SettingsMode;
}

export const settingsState = reactive<SettingsState>({
  settings: null,
  model: null,
  mode: 'normal',
});

export const isNormalMode = computed(() => settingsState.mode === 'normal');
export const isProMode = computed(() => settingsState.mode === 'pro');

export const getSettings = async () => {
  const [res] = await SettingsApi.get();
  if (res === null) return;
  settingsState.settings = res;
  settingsState.model = res;
};

export const updateSettings = async () => {
  const [, err] = await SettingsApi.update(settingsState.model);
  if (err !== null) {
    window.$message.error(err.message);
  }
  await getSettings();
  window.$message.success('更新成功');
};

export const formRules = {};

export const useSettings = () => {
  onMounted(getSettings);
};
