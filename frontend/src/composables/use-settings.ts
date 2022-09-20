import { computed, onMounted, reactive, watch } from 'vue';
import type { FormInst, FormItemRule, FormRules } from 'naive-ui';
import type { ApiResponseSettings } from '@/apis/SettingsApi';
import SettingsApi from '@/apis/SettingsApi';

export type SettingsMode = 'normal' | 'pro' | 'gfw';

export interface SettingsState {
  settings: ApiResponseSettings | null;
  model: ApiResponseSettings | null;
  mode: SettingsMode;
}

export const settingsState = reactive<SettingsState>({
  settings: null,
  model: null,
  mode: (localStorage.getItem('settingMode') as SettingsMode) || 'normal',
});

export const isNormalMode = computed(() => settingsState.mode === 'normal');
export const isProMode = computed(() => settingsState.mode === 'pro');
export const isGfwMode = computed(() => settingsState.mode === 'gfw');

export const getSettings = async () => {
  const [res] = await SettingsApi.get();
  if (res === null) return;
  settingsState.settings = res;
  settingsState.model = res;
};

export const updateSettings = async (form: FormInst | null = null) => {
  if (settingsState.model === null) return;
  if (form !== null) {
    const validate = new Promise((resolve) => {
      form.validate((errors) => {
        resolve(errors);
      });
    });
    const errors = await validate;
    if (errors !== undefined) return;
  }
  const [, err] = await SettingsApi.update(settingsState.model);
  if (err !== null) {
    window.$message.error(err.message);
  }
  await getSettings();
  // window.$message.success('更新成功');
};

watch(
  () => settingsState.mode,
  (val) => {
    localStorage.setItem('settingMode', settingsState.mode);
    if (settingsState.model === null) {
      return;
    }
    if (val === 'gfw') {
      settingsState.model.glider_strategy = 'ha';
      settingsState.model.health_check_url = settingsState.model.health_check_url || 'https://google.com';
      settingsState.model.manual_lb_port = settingsState.model.manual_lb_port || 10808;
      updateSettings();
    }
  }
);

export const formRules = computed((): FormRules => {
  if (isGfwMode.value) {
    return {
      health_check_url: [
        {
          // required: true,
          validator: (rule: FormItemRule, value: string) => {
            if (value === '') return new Error('请输入健康检查地址');
            return true;
          },
          trigger: 'blur',
        },
      ],
      manual_lb_port: [
        {
          validator: (rule: FormItemRule, value: number) => {
            if (!value) return new Error('请输入手动负载均衡端口');
            return true;
          },
          trigger: 'blur',
        },
      ],
    };
  }
  return {};
});

export const useSettings = () => {
  onMounted(getSettings);
};
