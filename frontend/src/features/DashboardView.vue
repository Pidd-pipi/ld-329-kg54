<template>
  <main class="page-shell" v-loading="loading">
    <AppHeader
      :unread="overview?.metrics.unread ?? 0"
      :users="perspectiveUsers"
      :user="store.currentUser"
      @update:user="store.currentUser = $event"
    />

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

      <div class="panel">
        <h2>智能匹配</h2>
        <FeatureCard v-for="match in overview.matches" :key="match.id" :title="`${match.provider} × ${match.learner}`" :description="match.recommendation">
          <template #tag><el-tag type="warning">{{ match.score }}%</el-tag></template>
          <p class="muted">{{ match.offerSkill }} ↔ {{ match.wantedSkill }}</p>
          <div class="tag-row">
            <el-tag v-for="slot in match.commonSlots" :key="slot">{{ slot }}</el-tag>
          </div>
          <div class="match-card__footer">
            <el-button size="small" type="primary" :disabled="!canInitiate(match)" @click="openCreateDialog(match)">
              发起预约
            </el-button>
            <small v-if="initiateHint(match)" class="muted">{{ initiateHint(match) }}</small>
          </div>
        </FeatureCard>
      </div>

      <AppointmentPanel
        :current-user="store.currentUser"
        @confirm="handleConfirm"
        @reject="handleReject"
        @revise="openReviseDialog"
      />

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
      v-model:visible="dialog.visible"
      :match="dialog.match"
      :appointment="dialog.appointment"
      :submitting="dialog.submitting"
      @submit="submitDialog"
    />
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import AppHeader from '../components/AppHeader.vue';
import FeatureCard from '../components/FeatureCard.vue';
import MetricCard from '../components/MetricCard.vue';
import RadarChart from '../components/RadarChart.vue';
import AppointmentFormDialog from './appointment/AppointmentFormDialog.vue';
import AppointmentPanel from './appointment/AppointmentPanel.vue';
import { fetchOverview } from '../services/storage.service';
import { useAppointmentStore } from '../stores/appointment.store';
import { APPOINTMENT_MESSAGES, APPOINTMENT_REJECT_CONFIRM } from '../constants/appointment.constants';
import type { Appointment, AppointmentFormPayload, Match, Overview } from '../types/domain';

const overview = ref<Overview | null>(null);
const loading = ref(true);
const error = ref('');
const store = useAppointmentStore();

const dialog = reactive<{
  visible: boolean;
  match: Match | null;
  appointment: Appointment | null;
  submitting: boolean;
}>({ visible: false, match: null, appointment: null, submitting: false });

const perspectiveUsers = computed(() => {
  if (!overview.value) {
    return [];
  }
  const names = new Set<string>([overview.value.profile.name]);
  for (const match of overview.value.matches) {
    names.add(match.provider);
    names.add(match.learner);
  }
  return [...names];
});

function isMatchParticipant(match: Match) {
  return store.currentUser === match.provider || store.currentUser === match.learner;
}

function canInitiate(match: Match) {
  return isMatchParticipant(match) && !store.activeMatchIds.has(match.id);
}

function initiateHint(match: Match) {
  if (!isMatchParticipant(match)) {
    return '仅匹配双方可发起';
  }
  if (store.activeMatchIds.has(match.id)) {
    return '已有进行中的预约';
  }
  return '';
}

function openCreateDialog(match: Match) {
  dialog.match = match;
  dialog.appointment = null;
  dialog.visible = true;
}

function openReviseDialog(appointment: Appointment) {
  dialog.match = overview.value?.matches.find((item) => item.id === appointment.matchId) ?? null;
  dialog.appointment = appointment;
  dialog.visible = true;
}

async function submitDialog(form: AppointmentFormPayload) {
  dialog.submitting = true;
  try {
    if (dialog.appointment) {
      await store.revise(dialog.appointment.id, form);
      ElMessage.success(APPOINTMENT_MESSAGES.revised);
    } else if (dialog.match) {
      await store.create({ ...form, matchId: dialog.match.id, actor: store.currentUser });
      ElMessage.success(APPOINTMENT_MESSAGES.created);
    }
    dialog.visible = false;
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : APPOINTMENT_MESSAGES.actionFailed);
  } finally {
    dialog.submitting = false;
  }
}

async function handleConfirm(appointment: Appointment) {
  try {
    await store.confirm(appointment.id);
    ElMessage.success(APPOINTMENT_MESSAGES.confirmed);
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : APPOINTMENT_MESSAGES.actionFailed);
  }
}

async function handleReject(appointment: Appointment) {
  try {
    await ElMessageBox.confirm(APPOINTMENT_REJECT_CONFIRM, '拒绝预约', { type: 'warning' });
  } catch {
    return;
  }
  try {
    await store.reject(appointment.id);
    ElMessage.success(APPOINTMENT_MESSAGES.rejected);
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : APPOINTMENT_MESSAGES.actionFailed);
  }
}

onMounted(async () => {
  try {
    overview.value = await fetchOverview();
    if (!store.currentUser) {
      store.currentUser = overview.value.profile.name;
    }
    await store.load();
  } catch (err) {
    error.value = err instanceof Error ? err.message : '加载失败';
  } finally {
    loading.value = false;
  }
});
</script>
