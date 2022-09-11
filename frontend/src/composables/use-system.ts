import { computed, reactive } from 'vue';
import type { ApiResponseSystemStatus } from '@/apis/CommonApi';
import CommonApi from '@/apis/CommonApi';
import useInterval from '@/composables/use-interval';

export interface SystemState {
  status: ApiResponseSystemStatus | null;
}

export const systemState = reactive<SystemState>({
  status: null,
});

export const isSetup = computed(() => systemState.status?.glider_downloaded && systemState.status?.xray_downloaded);

export const getSystemStatus = async () => {
  const [res, err] = await CommonApi.getSystemStatus();
  if (err !== null) {
    window.$message.error(err.message);
    return;
  }
  systemState.status = res;
};

export const useSystem = () => {
  getSystemStatus();
  const { stopInterval } = useInterval(getSystemStatus, 2000);
  return {
    stopInterval,
  };
};
