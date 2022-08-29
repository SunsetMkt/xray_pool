import { onMounted, reactive } from 'vue';
import SettingsApi from '@/apis/SettingsApi';
import type { ApiResponseSettings } from '@/interfaces/settings';

export interface SetupState {
  model: ApiResponseSettings | null;
  currentStep: number;
}

export const setupState = reactive<SetupState>({
  model: null,
  currentStep: 1,
});

export const useSetup = () => {
  const getDefaultSettings = async () => {
    const [res]: ApiResponseSettings[] = await SettingsApi.getDefaultSettings();
    setupState.model = res;
    return res;
  };

  onMounted(getDefaultSettings);
};

export const finishSetup = async () => {
  const [, err] = await SettingsApi.update(setupState.model);
  if (err !== null) {
    window.$message.error(err.message);
  }
  window.$message.success('初始化完成');
};
