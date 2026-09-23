import type { AppointmentStatus } from '../types/domain';

export const APPOINTMENT_STATUS_LABELS: Record<AppointmentStatus, string> = {
  pending: '等待对方确认',
  confirmed: '双方已确认',
  rejected: '已撤下',
};

export const APPOINTMENT_STATUS_TAG_TYPES: Record<AppointmentStatus, 'warning' | 'success' | 'info'> = {
  pending: 'warning',
  confirmed: 'success',
  rejected: 'info',
};

export const APPOINTMENT_MESSAGES = {
  created: '预约已发起，等待对方确认',
  confirmed: '已确认，预约生效',
  rejected: '已拒绝，预约撤下，可重新发起',
  revised: '预约已修改，等待对方重新确认',
  actionFailed: '预约操作失败，请稍后重试',
} as const;

export const APPOINTMENT_REJECT_CONFIRM = '拒绝后预约将撤下，双方可重新发起，确定拒绝吗？';

export const APPOINTMENT_PLACE_SUGGESTIONS = ['线上会议室', '图书馆研讨间', '东校区湖边', '中心校区咖啡厅'] as const;

export const APPOINTMENT_FORM_LIMITS = {
  placeMax: 60,
  agendaMax: 200,
} as const;
