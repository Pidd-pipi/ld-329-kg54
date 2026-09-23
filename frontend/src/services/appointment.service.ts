import type {
  Appointment,
  AppointmentCreatePayload,
  AppointmentRevisePayload,
} from '../types/domain';
import { AppException } from '../errors/AppException';
import { APPOINTMENT_MESSAGES } from '../constants/appointment.constants';

const API_BASE = '/api';

async function parseResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let message: string = APPOINTMENT_MESSAGES.actionFailed;
    try {
      const body = (await response.json()) as { message?: string };
      if (body.message) {
        message = body.message;
      }
    } catch {
      // 响应不是 JSON 时使用默认提示
    }
    throw new AppException(message);
  }
  return response.json() as Promise<T>;
}

function send<T>(url: string, method: string, body?: unknown): Promise<T> {
  return fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  }).then((response) => parseResponse<T>(response));
}

export function fetchAppointments(): Promise<Appointment[]> {
  return send<Appointment[]>(`${API_BASE}/appointments`, 'GET');
}

export function createAppointment(payload: AppointmentCreatePayload): Promise<Appointment> {
  return send<Appointment>(`${API_BASE}/appointments`, 'POST', payload);
}

export function confirmAppointment(id: number, actor: string): Promise<Appointment> {
  return send<Appointment>(`${API_BASE}/appointments/${id}/confirm`, 'POST', { actor });
}

export function rejectAppointment(id: number, actor: string): Promise<Appointment> {
  return send<Appointment>(`${API_BASE}/appointments/${id}/reject`, 'POST', { actor });
}

export function reviseAppointment(id: number, payload: AppointmentRevisePayload): Promise<Appointment> {
  return send<Appointment>(`${API_BASE}/appointments/${id}`, 'PUT', payload);
}
