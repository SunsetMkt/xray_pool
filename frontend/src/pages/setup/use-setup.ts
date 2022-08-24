import { onMounted, reactive } from 'vue';
import type { SubscribeModel } from '@/interfaces/subscribe';
import SettingsApi from '@/apis/SettingsApi';
import type { ApiResponseSettings } from '@/interfaces/settings';
import SubscribeApi from '@/apis/SubscribeApi';

export interface SetupState {
  model: ApiResponseSettings | null;
  currentStep: number;
  subscribeList: SubscribeModel[] | null;
}

export const setupState = reactive<SetupState>({
  model: null,
  currentStep: 1,
  subscribeList: [],
});

export const getSubscribeList = async () => {
  const [res] = await SubscribeApi.getList();
  setupState.subscribeList = res;
};

export const useSetup = () => {
  const getDefaultSettings = async () => {
    const [res]: ApiResponseSettings[] = await SettingsApi.get();
    setupState.model = res;
    return res;
  };

  onMounted(getDefaultSettings);
};

export const finishSetup = () => {
  window.$message.success('初始化完成');
};
