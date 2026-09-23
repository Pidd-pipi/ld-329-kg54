import type { AppointmentStatus } from '../types/appointment';

// 可切换的双方视角（演示账号）
export const VIEWER_USERS = ['林澈', '孟野', '周芮', '许安'] as const;

export const DEFAULT_VIEWER = '林澈';

export const APPOINTMENT_STATUS_PENDING: AppointmentStatus = 'pending';
export const APPOINTMENT_STATUS_CONFIRMED: AppointmentStatus = 'confirmed';

export const APPOINTMENT_STATUS_LABEL: Record<AppointmentStatus, string> = {
  pending: '等待对方确认',
  confirmed: '双方已确认',
};

// 地点快捷选项（线上/线下）
export const PLACE_SUGGESTIONS = ['线上会议室', '东校区湖边', '中心校区图书馆', '西校区琴房'] as const;

export const APPOINTMENT_MESSAGES = {
  initiated: '预约已发起，等待对方回应',
  revised: '协商内容已更新，需要对方重新确认',
  confirmed: '预约已确认生效',
  rejected: '对方拒绝了预约，已撤下，可重新发起',
  loadFailed: '无法加载校园技能交换数据',
  actionFailed: '操作失败，请稍后重试',
} as const;
