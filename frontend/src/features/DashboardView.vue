<template>
  <main class="page-shell" v-loading="loading">
    <AppHeader :unread="overview?.metrics.unread ?? 0">
      <template #viewer>
        <ViewerSwitcher
          :model-value="store.currentUser"
          @update:model-value="store.switchViewer"
        />
      </template>
    </AppHeader>

    <section v-if="overview" class="metrics-grid">
      <MetricCard label="已发布技能" :value="overview.metrics.skills" />
      <MetricCard label="活跃需求" :value="overview.metrics.needs" />
      <MetricCard label="智能匹配" :value="overview.metrics.matches" />
      <MetricCard label="评价记录" :value="overview.metrics.reviews" />
    </section>

    <el-alert v-if="error" :title="error" type="error" show-icon />

    <section v-if="overview" class="workspace-grid">
      <div class="panel">
        <h2>技能发布</h2>
        <FeatureCard v-for="skill in overview.skills" :key="skill.id" :title="skill.title" :description="skill.description">
          <template #tag><el-tag>{{ skill.category }} {{ skill.level }}%</el-tag></template>
          <div class="tag-row">
            <el-tag v-for="slot in skill.timeSlots" :key="slot" effect="plain">{{ slot }}</el-tag>
            <el-tag v-for="reward in skill.rewards" :key="reward" type="success" effect="plain">{{ reward }}</el-tag>
          </div>
          <small>{{ skill.owner }} · {{ skill.campus }} · {{ skill.portfolio }}</small>
        </FeatureCard>
      </div>

      <div class="panel">
        <h2>需求浏览</h2>
        <el-table :data="overview.needs" size="small">
          <el-table-column prop="title" label="需求" min-width="170" />
          <el-table-column prop="category" label="类别" width="82" />
          <el-table-column prop="campus" label="校区" width="96" />
          <el-table-column prop="responses" label="响应" width="72" sortable />
        </el-table>
      </div>

      <div class="panel negotiation-panel">
        <h2>智能匹配 · 预约协商</h2>
        <MatchNegotiationCard
          v-for="match in overview.matches"
          :key="match.id"
          :match="match"
          :appointment="store.appointmentByMatch.get(match.id)"
          :current-user="store.currentUser"
          @initiate="openCreate"
          @confirm="onConfirm"
          @reject="onReject"
          @revise="openRevise"
        />
      </div>

      <div class="panel">
        <h2>待回应预约</h2>
        <PendingAppointmentList
          :appointments="store.pendingAppointments"
          :matches="overview.matches"
          :current-user="store.currentUser"
          @confirm="onConfirm"
          @reject="onReject"
          @revise="openRevise"
        />
      </div>

      <div class="panel">
        <h2>预约确认</h2>
        <ConfirmedAppointmentList :appointments="store.confirmedAppointments" />
      </div>

      <div class="panel profile-panel">
        <div>
          <h2>个人主页与技能墙</h2>
          <h3>{{ overview.profile.name }}</h3>
          <p>{{ overview.profile.major }} · {{ overview.profile.creditLevel }}</p>
          <el-progress :percentage="overview.profile.creditScore" />
          <ul>
            <li v-for="item in overview.profile.history" :key="item">{{ item }}</li>
          </ul>
        </div>
        <RadarChart :radar="overview.profile.radar" />
      </div>

      <div class="panel">
        <h2>评价信用</h2>
        <FeatureCard v-for="review in overview.reviews" :key="review.id" :title="`${review.from} → ${review.to}`" :description="review.content">
          <template #tag><el-rate :model-value="review.rating" disabled size="small" /></template>
        </FeatureCard>
      </div>

      <div class="panel">
        <h2>消息通知</h2>
        <FeatureCard v-for="conversation in overview.messages" :key="conversation.id" :title="conversation.withUser" :description="conversation.messages.join(' / ')">
          <template #tag><el-badge :value="conversation.unread" /></template>
        </FeatureCard>
      </div>
    </section>

    <AppointmentFormDialog
      v-model="dialogVisible"
      :match="dialogMatchValue"
      :actor="store.currentUser"
      :existing="editingAppointment"
      :loading="store.busy"
      @submit="onSubmitDraft"
    />
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import AppHeader from '../components/AppHeader.vue';
import FeatureCard from '../components/FeatureCard.vue';
import MetricCard from '../components/MetricCard.vue';
import RadarChart from '../components/RadarChart.vue';
import ViewerSwitcher from '../components/ViewerSwitcher.vue';
import MatchNegotiationCard from '../components/MatchNegotiationCard.vue';
import PendingAppointmentList from '../components/PendingAppointmentList.vue';
import ConfirmedAppointmentList from '../components/ConfirmedAppointmentList.vue';
import AppointmentFormDialog from '../components/AppointmentFormDialog.vue';
import { fetchOverview } from '../services/storage.service';
import { useNegotiationStore } from '../stores/negotiation.store';
import { APPOINTMENT_MESSAGES } from '../constants/appointment.constants';
import type { Appointment, Match, Overview } from '../types/domain';
import type { AppointmentDraft } from '../types/appointment';

const overview = ref<Overview | null>(null);
const loading = ref(true);
const error = ref('');

const store = useNegotiationStore();

const dialogVisible = ref(false);
const dialogMatch = ref<Match | null>(null);
const editingAppointment = ref<Appointment | null>(null);

const dialogMatchValue = computed<Match>(() => dialogMatch.value ?? ({} as Match));

onMounted(async () => {
  try {
    overview.value = await fetchOverview();
    store.hydrate(overview.value.matches, overview.value.appointments);
  } catch (err) {
    error.value = err instanceof Error ? err.message : APPOINTMENT_MESSAGES.loadFailed;
  } finally {
    loading.value = false;
  }
});

function openCreate(match: Match): void {
  editingAppointment.value = null;
  dialogMatch.value = match;
  dialogVisible.value = true;
}

function openRevise(appointment: Appointment): void {
  const match = overview.value?.matches.find((item) => item.id === appointment.matchId);
  if (!match) return;
  editingAppointment.value = appointment;
  dialogMatch.value = match;
  dialogVisible.value = true;
}

async function onSubmitDraft(draft: AppointmentDraft): Promise<void> {
  try {
    if (editingAppointment.value) {
      await store.revise(editingAppointment.value.id, draft);
      ElMessage.success(APPOINTMENT_MESSAGES.revised);
    } else {
      await store.initiate(draft);
      ElMessage.success(APPOINTMENT_MESSAGES.initiated);
    }
    dialogVisible.value = false;
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : APPOINTMENT_MESSAGES.actionFailed);
  }
}

async function onConfirm(appointment: Appointment): Promise<void> {
  try {
    await store.confirm(appointment.id);
    ElMessage.success(APPOINTMENT_MESSAGES.confirmed);
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : APPOINTMENT_MESSAGES.actionFailed);
  }
}

async function onReject(appointment: Appointment): Promise<void> {
  try {
    await store.reject(appointment.id);
    ElMessage.info(APPOINTMENT_MESSAGES.rejected);
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : APPOINTMENT_MESSAGES.actionFailed);
  }
}
</script>
