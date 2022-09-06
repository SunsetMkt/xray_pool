<template>
  <n-form
    v-if="state.model"
    ref="formRef"
    :model="state.model"
    :rules="rules"
    label-placement="left"
    label-width="auto"
    require-mark-placement="right-hanging"
    size="medium"
    :style="{
      maxWidth: '640px',
    }"
  >
    <n-form-item label="Xray启动起始端口" path="xray_port_range">
      <n-input-number v-model:value="state.model.xray_port_range" />
    </n-form-item>

    <n-form-item label="Xray 是否开启 HTTP 端口" path="xray_open_socks_and_http">
      <n-switch v-model:value="state.model.xray_open_socks_and_http" />
    </n-form-item>

    <n-form-item label="单个节点 的测试超时时间（秒）" path="one_node_test_time_out">
      <n-input-number v-model:value="state.model.one_node_test_time_out" clearable>
        <!--        <template #suffix> 秒 </template>-->
      </n-input-number>
    </n-form-item>

    <n-form-item label="批量节点测试总超时时间（秒）" path="batch_node_test_max_time_out">
      <n-input-number v-model:value="state.model.batch_node_test_max_time_out" clearable>
        <!--        <template #suffix> 秒 </template>-->
      </n-input-number>
    </n-form-item>

    <n-form-item label="测速目标网站时，使用的并发线程数" path="test_url_thread">
      <n-input-number v-model:value="state.model.test_url_thread" clearable />
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import { setupState as state } from '@/pages/setup/use-setup';

const rules = {};
</script>
