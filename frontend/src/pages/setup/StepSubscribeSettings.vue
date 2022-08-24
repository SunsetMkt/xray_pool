<template>
  <n-form
    v-if="state.model"
    ref="formRef"
    :model="state.model"
    :rules="rules"
    label-placement="top"
    label-width="auto"
    require-mark-placement="right-hanging"
    size="medium"
    :style="{
      maxWidth: '640px',
    }"
  >
    <n-form-item label="添加订阅源" path="app_start_port">
      <n-input v-model:value="state.model.app_start_port" clearable />
    </n-form-item>
  </n-form>

  <div>订阅源列表</div>
  <n-list class="border-1" hoverable clickable>
    <n-list-item>
      <div></div>
    </n-list-item>
  </n-list>
</template>

<script setup lang="ts">
import { setupState as state } from '@/pages/setup/use-setup';
</script>
