<template>
  <template v-if="appointments.length">
    <PendingAppointmentCard
      v-for="appointment in appointments"
      :key="appointment.id"
      :appointment="appointment"
      :current-user="currentUser"
      :match="matchById(appointment.matchId)"
      show-match
      @confirm="emit('confirm', $event)"
      @reject="emit('reject', $event)"
      @revise="emit('revise', $event)"
    />
  </template>
  <el-empty v-else description="没有待回应的预约" :image-size="72" />
</template>

<script setup lang="ts">
import type { Appointment, Match } from '../types/domain';
import PendingAppointmentCard from './PendingAppointmentCard.vue';

const props = defineProps<{
  appointments: Appointment[];
  matches: Match[];
  currentUser: string;
}>();

const emit = defineEmits<{
  (e: 'confirm', appointment: Appointment): void;
  (e: 'reject', appointment: Appointment): void;
  (e: 'revise', appointment: Appointment): void;
}>();

function matchById(id: number): Match | undefined {
  return props.matches.find((match) => match.id === id);
}
</script>
