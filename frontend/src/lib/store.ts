import { writable } from 'svelte/store';

export type ServerStatus = 'Stopped' | 'Starting' | 'Running';

export const serverStatus = writable<ServerStatus>('Stopped');
export const logs = writable<string[]>([]);

let ws: WebSocket | null = null;

export function connectWebSocket() {
	if (ws) return;

	ws = new WebSocket('ws://localhost:8080/ws');

	ws.onmessage = (event) => {
		try {
			const msg = JSON.parse(event.data);
			if (msg.type === 'status') {
				serverStatus.set(msg.data);
			} else if (msg.type === 'log') {
				logs.update((current) => [...current, msg.data]);
			}
		} catch (e) {
			console.error('Failed to parse websocket message', e);
		}
	};

	ws.onclose = () => {
		console.log('WebSocket disconnected, attempting to reconnect...');
		ws = null;
		setTimeout(connectWebSocket, 3000);
	};

	ws.onerror = (err) => {
		console.error('WebSocket error:', err);
		ws?.close();
	};
}
