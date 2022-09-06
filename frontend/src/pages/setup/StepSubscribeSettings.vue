<template>
  <n-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-placement="top"
    label-width="auto"
    require-mark-placement="right-hanging"
    size="medium"
    inline
    :style="{
      maxWidth: '640px',
    }"
  >
    <n-form-item label="订阅源名称" path="name">
      <n-input v-model:value="form.name" clearable placeholder="订阅源名称" />
    </n-form-item>

    <n-form-item label="订阅源URL" path="url">
      <n-input v-model:value="form.url" clearable placeholder="订阅源URL" />
    </n-form-item>

    <n-form-item>
      <n-button type="primary" @click="addSubscribe">添加</n-button>
    </n-form-item>
  </n-form>

  <div class="font-bold">订阅源列表</div>
  <div v-if="state.subscribeList?.length === 0" class="text-gray-500">当前没有订阅源，请通过上方表单添加</div>
  <n-list v-else class="border-1 mt-2" hoverable clickable>
    <n-list-item v-for="(item, i) in state.subscribeList" :key="item.name">
      <div class="flex row">
        <div class="flex-1">
          <div>{{ item.name }}</div>
          <div class="text-gray-500">{{ item.url }}</div>
        </div>

        <n-popconfirm @positive-click="removeSubscribe(i)">
          <template #trigger>
            <n-button size="tiny" type="error">删除</n-button>
          </template>
          确定删除该订阅源？
        </n-popconfirm>
        <div></div>
      </div>
    </n-list-item>
  </n-list>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import type { FormInst } from 'naive-ui';
import SubscribeApi from '@/apis/SubscribeApi';
import { getSubscribeList, subscribeState as state } from '@/composables/use-subscribe';

const formRef = ref<FormInst | null>(null);
const form = reactive({
  name: '',
  url: '',
});

const rules = {
  name: {
    required: true,
    message: '请输入名称',
    trigger: 'blur',
  },

  url: {
    required: true,
    message: '请输入URL',
    trigger: 'blur',
  },
};

const addSubscribe = async () => {
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      const [, err] = await SubscribeApi.add(form);
      if (err !== null) {
        window.$message.success('添加成功');
      }
      form.name = '';
      form.url = '';
      getSubscribeList();
    }
  });
};

const removeSubscribe = async (id: number) => {
  const [, err] = await SubscribeApi.remove(id + 1);
  if (err === null) {
    window.$message.success('删除成功');
  }
  getSubscribeList();
};

onMounted(() => {
  getSubscribeList();
});
</script>
