<template>
  <header class="app-header">
    <div>
      <p class="eyebrow">Skill Swap Campus</p>
      <h1>{{ APP_NAME }}</h1>
      <p>{{ APP_TAGLINE }}</p>
    </div>
    <div class="app-header__side">
      <PerspectiveSwitch
        :users="users"
        :model-value="user"
        @update:model-value="emit('update:user', $event)"
      />
      <el-badge :value="unread" class="message-badge">
        <el-button type="primary">站内消息</el-button>
      </el-badge>
    </div>
  </header>
</template>

<script setup lang="ts">
import { APP_NAME, APP_TAGLINE } from '../constants/app.constants';
import PerspectiveSwitch from '../features/appointment/PerspectiveSwitch.vue';

defineProps<{ unread: number; users: string[]; user: string }>();
const emit = defineEmits<{ 'update:user': [value: string] }>();
</script>
