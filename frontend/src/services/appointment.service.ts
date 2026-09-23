import type { Appointment } from '../types/domain';
import type { AppointmentDraft } from '../types/appointment';

const API_BASE = '/api';
const APPOINTMENTS_ENDPOINT = `${API_BASE}/appointments`;

async function parseError(response: Response, fallback: string): Promise<string> {
  try {
    const body = (await response.json()) as { message?: string };
    return body.message || fallback;
  } catch {
    return fallback;
  }
}

export async function createAppointment(draft: AppointmentDraft): Promise<Appointment> {
  const response = await fetch(APPOINTMENTS_ENDPOINT, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(draft),
  });
  if (!response.ok) {
    throw new Error(await parseError(response, '发起预约失败'));
  }
  return response.json() as Promise<Appointment>;
}

export async function reviseAppointment(
  id: number,
  draft: AppointmentDraft,
): Promise<Appointment> {
  const response = await fetch(`${APPOINTMENTS_ENDPOINT}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(draft),
  });
  if (!response.ok) {
    throw new Error(await parseError(response, '修改预约失败'));
  }
  return response.json() as Promise<Appointment>;
}

export async function confirmAppointment(id: number, actor: string): Promise<Appointment> {
  const response = await fetch(`${APPOINTMENTS_ENDPOINT}/${id}/confirm`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ actor }),
  });
  if (!response.ok) {
    throw new Error(await parseError(response, '确认预约失败'));
  }
  return response.json() as Promise<Appointment>;
}

export async function rejectAppointment(id: number, actor: string): Promise<void> {
  const response = await fetch(`${APPOINTMENTS_ENDPOINT}/${id}/reject`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ actor }),
  });
  if (!response.ok) {
    throw new Error(await parseError(response, '拒绝预约失败'));
  }
}
