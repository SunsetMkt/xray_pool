<template>
  <div>
    <n-form
      v-if="settingsState.model"
      ref="formRef"
      :model="settingsState.model"
      label-placement="left"
      label-width="auto"
      require-mark-placement="right-hanging"
      size="small"
      :key="settingsState.mode"
    >
      <n-form-item label="Mux多路复用" path="main_proxy_settings.RoutingStrategy">
        <n-select
          v-model:value="settingsState.model.main_proxy_settings.RoutingStrategy"
          :options="routingStrategyOptions"
          @update:value="updateSettings"
        />
      </n-form-item>

      <n-form-item label="直连局域网和大陆" path="main_proxy_settings.BypassLANAndMainLand">
        <n-switch
          v-model:value="settingsState.model.main_proxy_settings.BypassLANAndMainLand"
          @update:value="updateSettings"
        />
      </n-form-item>
    </n-form>

    <div class="flex row justify-between items-center">
      <div>路由规则：</div>
      <div>
        <n-radio-group v-model:value="routingType" size="small">
          <n-radio-button value="Proxy" label="代理" />
          <n-radio-button value="Direct" label="直连" />
          <n-radio-button value="Block" label="阻止" />
        </n-radio-group>
      </div>
    </div>
    <n-list v-if="rules" class="rule-list border-1 p-2 mt-1" hoverable clickable :show-divider="false">
      <n-list-item v-for="(item, i) in currentRuleList" :key="item">
        <div class="flex row">
          <div class="flex-1">
            <div>{{ item }}</div>
          </div>

          <n-popconfirm @positive-click="removeRule(i)">
            <template #trigger>
              <n-button size="tiny" type="error">删除</n-button>
            </template>
            确定删除该规则？
          </n-popconfirm>
        </div>
      </n-list-item>
    </n-list>

    <div class="flex row mt-1 items-center">
      <div class="flex-1">
        <n-input v-model:value="newRule" placeholder="输入新规则" />
      </div>
      <div class="ml-1">
        <n-button type="primary" size="small" @click="addRule">新增</n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { settingsState, updateSettings } from '@/composables/use-settings';
import RoutingApi from '@/apis/RoutingApi';
import type { RoutingResponseModel, RoutingType } from '@/interfaces/Routing';

const routingStrategyOptions = [
  { label: 'AsIs', value: 'AsIs' },
  { label: 'IPIfNonMatch', value: 'IPIfNonMatch' },
  { label: 'IPOnDemand', value: 'IPOnDemand' },
];

const rules = ref<RoutingResponseModel | null>(null);
const routingType = ref<RoutingType>('Proxy');
const newRule = ref('');

const currentRuleList = computed(() => {
  if (rules.value === null) {
    return [];
  }

  switch (routingType.value) {
    case 'Proxy':
      return rules.value.proxy_list.rules;
    case 'Direct':
      return rules.value.direct_list.rules;
    case 'Block':
      return rules.value.block_list.rules;
    default:
      return [];
  }
});

const getRoutingRules = async () => {
  const [res, err] = await RoutingApi.getList();
  if (err) {
    window.$message.error(err.message);
    return;
  }
  rules.value = res;
};

const removeRule = async (index: number) => {
  const [, err] = await RoutingApi.remove({
    routing_type: routingType.value,
    index_list: [index + 1],
  });
  if (err) {
    window.$message.error(err.message);
    return;
  }

  getRoutingRules();
};

const addRule = async () => {
  if (!newRule.value) {
    window.$message.error('请输入规则');
    return;
  }
  const [, err] = await RoutingApi.add({
    routing_type: routingType.value,
    rules: [newRule.value],
  });
  if (err) {
    window.$message.error(err.message);
    return;
  }

  newRule.value = '';

  getRoutingRules();
};

onMounted(getRoutingRules);
</script>

<style>
.rule-list {
  height: 200px;
  overflow-x: auto;
}
</style>
