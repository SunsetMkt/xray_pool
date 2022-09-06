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
    <n-form-item label="Mux多路复用" path="main_proxy_settings.Mux">
      <n-switch v-model:value="settingsState.model.main_proxy_settings.Mux" @change="updateSettings" />
    </n-form-item>

    <n-form-item label="开启UDP" path="main_proxy_settings.RelayUDP">
      <n-switch v-model:value="settingsState.model.main_proxy_settings.RelayUDP" @change="updateSettings" />
    </n-form-item>

    <n-form-item label="开启流量探测" path="main_proxy_settings.Sniffing">
      <n-switch v-model:value="settingsState.model.main_proxy_settings.Sniffing" @change="updateSettings" />
    </n-form-item>

    <n-form-item label="国外的DNS" path="main_proxy_settings.DNSForeign">
      <n-input v-model:value="settingsState.model.main_proxy_settings.DNSForeign" clearable @blur="updateSettings" />
    </n-form-item>

    <n-form-item label="国内的DNS" path="main_proxy_settings.DNSDomestic">
      <n-input v-model:value="settingsState.model.main_proxy_settings.DNSDomestic" clearable @blur="updateSettings" />
    </n-form-item>

    <n-form-item label="国内的DNS备用" path="main_proxy_settings.DNSDomesticBackup">
      <n-input
        v-model:value="settingsState.model.main_proxy_settings.DNSDomesticBackup"
        clearable
        @blur="updateSettings"
      />
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import { settingsState, formRules, updateSettings } from '@/composables/use-settings';
</script>
