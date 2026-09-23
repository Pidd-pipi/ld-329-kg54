<template>
  <div class="viewer-switcher">
    <span class="viewer-switcher__label">当前视角</span>
    <el-radio-group :model-value="modelValue" size="small" @change="onChange">
      <el-radio-button v-for="user in users" :key="user" :value="user">{{ user }}</el-radio-button>
    </el-radio-group>
  </div>
</template>

<script setup lang="ts">
import { VIEWER_USERS } from '../constants/appointment.constants';

defineProps<{ modelValue: string }>();
const emit = defineEmits<{ (e: 'update:modelValue', value: string): void }>();

const users = VIEWER_USERS;

function onChange(value: string | number | boolean): void {
  emit('update:modelValue', String(value));
}
</script>
