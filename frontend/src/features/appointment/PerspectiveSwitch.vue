<template>
  <div class="perspective-switch">
    <span class="perspective-switch__label">当前视角</span>
    <el-radio-group
      :model-value="modelValue"
      @update:model-value="emit('update:modelValue', String($event))"
    >
      <el-radio-button v-for="user in users" :key="user" :value="user">{{ user }}</el-radio-button>
    </el-radio-group>
  </div>
</template>

<script setup lang="ts">
defineProps<{ users: string[]; modelValue: string }>();
const emit = defineEmits<{ 'update:modelValue': [value: string] }>();
</script>
