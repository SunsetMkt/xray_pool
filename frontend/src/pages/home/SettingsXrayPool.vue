<template>
  <n-form
    v-if="settingsState.model"
    ref="formRef"
    :model="settingsState.model"
    :rules="formRules"
    label-placement="left"
    label-width="auto"
    require-mark-placement="right-hanging"
    size="small"
    :key="settingsState.mode"
    :disabled="!isStopped"
  >
    <n-form-item>
      <n-radio-group v-model:value="settingsState.mode" name="radiogroup">
        <n-space>
          <n-radio value="normal"> 简易模式 </n-radio>
          <n-radio value="gfw"> 科学上网 </n-radio>
          <n-radio value="pro"> 专业模式 </n-radio>
        </n-space>
      </n-radio-group>
    </n-form-item>

    <n-form-item label="目标网站" path="test_url">
      <n-input v-model:value="settingsState.model.test_url" @blur="handleUpdateSettings" />
    </n-form-item>

    <n-form-item v-if="isProMode || isGfwMode" label="健康检测网站" path="health_check_url">
      <n-input v-model:value="settingsState.model.health_check_url" @blur="handleUpdateSettings" />
    </n-form-item>

    <n-form-item v-if="isProMode || isGfwMode" label="健康检测间隔（秒）" path="health_check_interval">
      <n-input-number
        class="w-full"
        v-model:value="settingsState.model.health_check_interval"
        @blur="handleUpdateSettings"
      />
    </n-form-item>

    <n-form-item v-if="isNormalMode || isGfwMode" label="本机性能" path="test_url_thread">
      <n-radio-group
        v-model:value="settingsState.model.test_url_thread"
        name="radiogroup"
        @change="handleUpdateSettings"
      >
        <n-space>
          <n-radio :value="3"> 弱鸡（3线程） </n-radio>
          <n-radio :value="10"> 一般（10线程） </n-radio>
          <n-radio :value="20"> 很猛（20线程） </n-radio>
        </n-space>
      </n-radio-group>
    </n-form-item>

    <n-form-item label="负载均衡策略" path="glider_strategy">
      <n-select
        v-model:value="settingsState.model.glider_strategy"
        :options="gliderStrategyOptions"
        @update:value="handleUpdateSettings"
      ></n-select>
    </n-form-item>

    <n-form-item v-if="isProMode || isGfwMode" label="自定义负载均衡端口（0为动态端口）" path="manual_lb_port">
      <n-input-number class="w-full" v-model:value="settingsState.model.manual_lb_port" @blur="handleUpdateSettings" />
    </n-form-item>

    <n-form-item v-if="isProMode" label="Xray启动起始端口" path="xray_port_range">
      <n-input-number class="w-full" v-model:value="settingsState.model.xray_port_range" @blur="handleUpdateSettings" />
    </n-form-item>

    <n-form-item v-if="isProMode" label="Xray 是否开启 HTTP 端口" path="xray_open_socks_and_http">
      <n-switch
        v-model:value="settingsState.model.xray_open_socks_and_http"
        @change="handleUpdateSettings"
        :disabled="settingsState.model.test_url_hard_way"
      />
    </n-form-item>

    <n-form-item v-if="isProMode" label="单个节点 的测试超时时间（秒）" path="one_node_test_time_out">
      <n-input-number
        class="w-full"
        v-model:value="settingsState.model.one_node_test_time_out"
        @blur="handleUpdateSettings"
      >
        <!--        <template #suffix> 秒 </template>-->
      </n-input-number>
    </n-form-item>

    <n-form-item v-if="isProMode" label="批量节点测试总超时时间（秒）" path="batch_node_test_max_time_out">
      <n-input-number
        class="w-full"
        v-model:value="settingsState.model.batch_node_test_max_time_out"
        @blur="handleUpdateSettings"
      >
        <!--        <template #suffix> 秒 </template>-->
      </n-input-number>
    </n-form-item>

    <n-form-item v-if="isProMode" path="test_url_thread">
      <template #label>
        <div>
          <span>测速目标网站时，使用的并发线程数</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><help-circle /></n-icon>
            </template>
            <span>请根据实际带宽设置并发线程，否则可能导致测速不准</span>
          </n-tooltip>
        </div>
      </template>
      <div class="w-full">
        <n-input-number
          class="w-full"
          v-model:value="settingsState.model.test_url_thread"
          @blur="handleUpdateSettings"
        />
      </div>
    </n-form-item>

    <n-form-item v-if="isProMode" path="test_url_hard_way">
      <template #label>
        <div>
          <span>启动浏览器进行测速</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><help-circle /></n-icon>
            </template>
            <span>应对一些网站的反爬虫策略</span>
          </n-tooltip>
        </div>
      </template>
      <n-switch v-model:value="settingsState.model.test_url_hard_way" @change="handleTestUrlHardWayChange" />
    </n-form-item>

    <n-form-item v-if="isProMode" path="test_url_hard_way_load_picture">
      <template #label>
        <div>
          <span>浏览器测速是否加载图片</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><help-circle /></n-icon>
            </template>
            <span>有些爬虫任务无需图片，这样效率更高，流量更低</span>
          </n-tooltip>
        </div>
      </template>
      <n-switch v-model:value="settingsState.model.test_url_hard_way_load_picture" @change="handleUpdateSettings" />
    </n-form-item>

    <n-form-item v-if="isProMode" path="test_url_failed_words">
      <template #label>
        <div>
          <span>测速成功的关键字</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><help-circle /></n-icon>
            </template>
            <span>如果测速目标网站返回结果中包含这些关键字，则认为测速成功，不区分大小写</span>
          </n-tooltip>
        </div>
      </template>

      <n-switch v-model:value="settingsState.model.test_url_succeed_words_enable" @change="handleUpdateSettings" />
      <n-select
        :disabled="!settingsState.model.test_url_succeed_words_enable"
        class="w-full ml-2"
        v-model:value="settingsState.model.test_url_succeed_words"
        multiple
        tag
        filterable
        :options="[]"
        placeholder="输入后按回车新增一条关键字"
        @update:value="handleUpdateSettings"
      />
    </n-form-item>

    <n-form-item v-if="isProMode" path="test_url_failed_words">
      <template #label>
        <div>
          <span>测速失败的关键字</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><help-circle /></n-icon>
            </template>
            <span>如果测速目标网站返回结果中包含这些关键字，则认为测速失败，不区分大小写</span>
          </n-tooltip>
        </div>
      </template>

      <n-switch v-model:value="settingsState.model.test_url_failed_words_enable" @change="handleUpdateSettings" />
      <n-select
        class="w-full ml-2"
        v-model:value="settingsState.model.test_url_failed_words_enable"
        :disabled="!settingsState.model.test_url_failed_words_enable"
        multiple
        tag
        filterable
        :options="[]"
        placeholder="输入后按回车新增一条关键字"
        @update:value="handleUpdateSettings"
      />
    </n-form-item>

    <n-form-item v-if="isProMode" label="测速失败的关键字（正则）" path="test_url_failed_regex">
      <div class="w-full">
        <n-input
          class="w-full"
          v-model:value="settingsState.model.test_url_failed_regex"
          @blur="handleUpdateSettings"
        />
      </div>
    </n-form-item>

    <n-form-item v-if="isProMode" path="test_url_status_code">
      <template #label>
        <div>
          <span>期望返回的HTTP状态码</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><help-circle /></n-icon>
            </template>
            <span>若返回的状态码不一致，则认为测速失败</span>
          </n-tooltip>
        </div>
      </template>
      <n-input-number
        class="w-full"
        v-model:value="settingsState.model.test_url_status_code"
        @blur="handleUpdateSettings"
      />
    </n-form-item>

    <div v-if="isProMode">
      <n-divider class="!my-3" />
      <div class="font-bold">代理设置</div>
      <div class="text-gray-400">
        此处代理设置的作用：软件添加“订阅源”的时候，可能这时候就需要代理才能够访问这些链接，以及在“浏览器”加载插件的时候，需要代理才能够取下载
        Adblock 插件。
      </div>

      <n-form-item v-if="isProMode" label="是否启用代理" path="proxy_info_settings.enable">
        <n-switch v-model:value="settingsState.model.proxy_info_settings.enable" @change="handleUpdateSettings" />
      </n-form-item>

      <n-form-item label="代理地址" path="test_url">
        <n-input
          v-model:value="settingsState.model.proxy_info_settings.http_url"
          :disabled="!settingsState.model.proxy_info_settings.enable"
          @blur="handleUpdateSettings"
        />
      </n-form-item>

      <n-form-item label="代理端口" path="health_check_url">
        <n-input-number
          class="w-full"
          v-model:value="settingsState.model.proxy_info_settings.http_port"
          :disabled="!settingsState.model.proxy_info_settings.enable"
          @blur="handleUpdateSettings"
        />
      </n-form-item>
    </div>
  </n-form>
</template>

<script setup lang="ts">
import { HelpCircle } from '@vicons/ionicons5';
import { ref } from 'vue';
import { NForm } from 'naive-ui';
import type { FormInst } from 'naive-ui';
import {
  settingsState,
  formRules,
  isProMode,
  isNormalMode,
  updateSettings,
  isGfwMode,
} from '@/composables/use-settings';
import { isStopped } from '@/composables/use-proxy-pool';

const formRef = ref<FormInst | null>(null);

const gliderStrategyOptions = [
  { label: 'rr(round robin)', value: 'rr' },
  { label: 'ha(high availability)', value: 'ha' },
  { label: 'lha(latency based high availability)', value: 'lha' },
  { label: 'dh(destination hashing)', value: 'dh' },
];

const handleUpdateSettings = () => {
  updateSettings(formRef.value);
};

const handleTestUrlHardWayChange = () => {
  handleUpdateSettings();
  if (settingsState.model === null) return;
  if (settingsState.model.test_url_hard_way === true) {
    settingsState.model.xray_open_socks_and_http = true;
  }
};
</script>
