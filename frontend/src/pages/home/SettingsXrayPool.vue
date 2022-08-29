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
  >
    <n-form-item>
      <n-radio-group v-model:value="settingsState.mode" name="radiogroup">
        <n-space>
          <n-radio value="normal"> 简易模式 </n-radio>
          <n-radio value="pro"> 专业模式 </n-radio>
        </n-space>
      </n-radio-group>
    </n-form-item>

    <n-form-item label="目标网站" path="test_url">
      <n-input v-model:value="settingsState.model.test_url" clearable />
    </n-form-item>

    <n-form-item label="启动 Xray 的实例数量" path="test_url">
      <div class="w-full">
        <n-input-number v-model:value="settingsState.model.xray_instance_count" clearable />
        <div class="text-gray-500">PS：数量决定了同时开启节点的数量</div>
      </div>
    </n-form-item>

    <n-form-item v-if="isNormalMode" label="本机性能" path="test_url_thread">
      <n-radio-group v-model:value="settingsState.model.test_url_thread" name="radiogroup">
        <n-space>
          <n-radio :value="3"> 弱鸡 </n-radio>
          <n-radio :value="10"> 一般 </n-radio>
          <n-radio :value="20"> 很猛 </n-radio>
        </n-space>
      </n-radio-group>
    </n-form-item>

    <n-form-item v-if="isProMode" label="负载均衡策略" path="test_url"> </n-form-item>

    <n-form-item v-if="isProMode" label="Xray启动起始端口" path="xray_port_range">
      <n-input-number class="w-full" v-model:value="settingsState.model.xray_port_range" />
    </n-form-item>

    <n-form-item v-if="isProMode" label="Xray 是否开启 HTTP 端口" path="xray_open_socks_and_http">
      <n-switch v-model:value="settingsState.model.xray_open_socks_and_http" />
    </n-form-item>

    <n-form-item v-if="isProMode" label="单个节点 的测试超时时间（秒）" path="one_node_test_time_out">
      <n-input-number class="w-full" v-model:value="settingsState.model.one_node_test_time_out">
        <!--        <template #suffix> 秒 </template>-->
      </n-input-number>
    </n-form-item>

    <n-form-item v-if="isProMode" label="批量节点测试总超时时间（秒）" path="batch_node_test_max_time_out">
      <n-input-number class="w-full" v-model:value="settingsState.model.batch_node_test_max_time_out">
        <!--        <template #suffix> 秒 </template>-->
      </n-input-number>
    </n-form-item>

    <n-form-item v-if="isProMode" label="测速目标网站时，使用的并发线程数" path="test_url_thread">
      <n-input-number class="w-full" v-model:value="settingsState.model.test_url_thread" />
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import { settingsState, useSettings, formRules, isProMode, isNormalMode } from '@/composables/use-settings';

useSettings();
</script>
