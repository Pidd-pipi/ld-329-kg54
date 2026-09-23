<template>
  <FeatureCard :title="`${match.provider} × ${match.learner}`" :description="match.recommendation">
    <template #tag><el-tag type="warning">{{ match.score }}%</el-tag></template>
    <p class="muted">{{ match.offerSkill }} ↔ {{ match.wantedSkill }}</p>
    <div class="tag-row">
      <el-tag v-for="slot in match.commonSlots" :key="slot">{{ slot }}</el-tag>
    </div>

    <PendingAppointmentCard
      v-if="appointment && appointment.status === 'pending'"
      :appointment="appointment"
      :current-user="currentUser"
      @confirm="emit('confirm', $event)"
      @reject="emit('reject', $event)"
      @revise="emit('revise', $event)"
    />

    <template v-else-if="appointment && appointment.status === 'confirmed'">
      <el-tag type="success">已确认：{{ appointment.time }} · {{ appointment.place }}</el-tag>
      <p class="muted">{{ appointment.agenda }}</p>
    </template>

    <template v-else-if="involved">      <el-button size="small" type="primary" @click="emit('initiate', match)">发起预约</el-button>
      <small class="hint">选择共同可约时间，填写地点与议程</small>
    </template>

    <el-tag v-else type="info" effect="plain">切换到 {{ match.provider }} 或 {{ match.learner }} 视角可预约</el-tag>
  </FeatureCard>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Appointment, Match } from '../types/domain';
import FeatureCard from './FeatureCard.vue';
import PendingAppointmentCard from './PendingAppointmentCard.vue';

const props = defineProps<{
  match: Match;
  appointment?: Appointment;
  currentUser: string;
}>();

const emit = defineEmits<{
  (e: 'initiate', match: Match): void;
  (e: 'confirm', appointment: Appointment): void;
  (e: 'reject', appointment: Appointment): void;
  (e: 'revise', appointment: Appointment): void;
}>();

const involved = computed(
  () => props.match.provider === props.currentUser || props.match.learner === props.currentUser,
);
</script>

<style scoped>
.hint { margin-left: 8px; color: #667085; }
</style>
