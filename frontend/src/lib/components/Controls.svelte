<script lang="ts">
	import { serverStatus } from '$lib/store';
	import { writable } from 'svelte/store';

	const loading = writable(false);

	async function callApi(endpoint: string) {
		loading.set(true);
		try {
			await fetch(`http://localhost:8080/api/server/${endpoint}`, {
				method: 'POST'
			});
		} catch (e) {
			console.error(`Failed to ${endpoint} server`, e);
		} finally {
			loading.set(false);
		}
	}

	function start() {
		callApi('start');
	}

	function stop() {
		callApi('stop');
	}

	function restart() {
		// Restart implementation (stop then start when stopped is implemented via WS)
		callApi('stop').then(() => {
			setTimeout(() => {
				callApi('start');
			}, 3000); // basic wait, ideally wait for 'Stopped' status
		});
	}
</script>

<div
	class="col-span-1 md:col-span-3 flex justify-center gap-6 p-6 bg-black/40 border-2 border-vapor-purple rounded-lg backdrop-blur-sm shadow-[0_0_20px_var(--color-vapor-purple)] z-10"
>
	<button
		onclick={start}
		disabled={$serverStatus !== 'Stopped' || $loading}
		class="px-8 py-3 bg-transparent border-2 border-vapor-cyan text-vapor-cyan font-black tracking-[0.2em] uppercase transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-vapor-cyan hover:text-black hover:shadow-[0_0_20px_var(--color-vapor-cyan)] hover:scale-105"
	>
		START
	</button>
	<button
		onclick={stop}
		disabled={$serverStatus !== 'Running' || $loading}
		class="px-8 py-3 bg-transparent border-2 border-vapor-pink text-vapor-pink font-black tracking-[0.2em] uppercase transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-vapor-pink hover:text-black hover:shadow-[0_0_20px_var(--color-vapor-pink)] hover:scale-105"
	>
		STOP
	</button>
	<button
		onclick={restart}
		disabled={$serverStatus === 'Stopped' || $loading}
		class="px-8 py-3 bg-transparent border-2 border-vapor-purple text-vapor-purple font-black tracking-[0.2em] uppercase transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-vapor-purple hover:text-black hover:shadow-[0_0_20px_var(--color-vapor-purple)] hover:scale-105"
	>
		RESTART
	</button>
</div>
