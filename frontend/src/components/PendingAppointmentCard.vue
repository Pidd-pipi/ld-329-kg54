<template>
  <FeatureCard :title="appointment.pair" :description="appointment.agenda">
    <template #tag>
      <el-tag type="warning">{{ statusLabel }} · v{{ appointment.version }}</el-tag>
    </template>
    <p class="muted">🕑 {{ appointment.time }} · 📍 {{ appointment.place }}</p>
    <p v-if="showMatch" class="muted">
      {{ match?.offerSkill }} ↔ {{ match?.wantedSkill }}
    </p>
    <small v-if="isInitiator">
      你发起的预约，对方确认前可修改；修改后需要对方重新确认。
    </small>
    <small v-else>{{ appointment.initiator }} 发起，等待你的回应。</small>
    <AppointmentActions
      :appointment="appointment"
      :current-user="currentUser"
      @confirm="emit('confirm', appointment)"
      @reject="emit('reject', appointment)"
      @revise="emit('revise', appointment)"
    />
  </FeatureCard>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Appointment, Match } from '../types/domain';
import FeatureCard from './FeatureCard.vue';
import AppointmentActions from './AppointmentActions.vue';
import { APPOINTMENT_STATUS_LABEL } from '../constants/appointment.constants';

const props = defineProps<{
  appointment: Appointment;
  currentUser: string;
  match?: Match;
  showMatch?: boolean;
}>();

const emit = defineEmits<{
  (e: 'confirm', appointment: Appointment): void;
  (e: 'reject', appointment: Appointment): void;
  (e: 'revise', appointment: Appointment): void;
}>();

const statusLabel = computed(
  () => APPOINTMENT_STATUS_LABEL[props.appointment.status],
);
const isInitiator = computed(() => props.appointment.initiator === props.currentUser);
</script>
