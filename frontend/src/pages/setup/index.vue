<template>
  <n-form
    ref="formRef"
    :model="model"
    :rules="rules"
    label-placement="left"
    label-width="auto"
    require-mark-placement="right-hanging"
    size="medium"
    :style="{
      maxWidth: '640px',
    }"
  >
    <n-form-item label="Xray启动起始端口" path="inputValue">
      <n-input-number v-model:value="model.inputValue" clearable />
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Settings } from '@/types/http-api';

const model = ref<Settings | null>(null);
</script>
