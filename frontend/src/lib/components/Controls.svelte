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
	class="col-span-1 md:col-span-3 flex justify-between items-center px-10 py-6 bg-[#1a1130]/80 border-[3px] border-[#3a448a] rounded-[16px] backdrop-blur-md shadow-[0_0_15px_rgba(88,101,242,0.4)] z-10 mb-2"
>
	<div class="flex flex-col">
		<div class="flex items-center gap-4 text-[#ffffff]">
			<h1 class="text-4xl font-bold tracking-widest drop-shadow-[0_0_8px_rgba(255,255,255,0.8)]">SERVER STATUS:</h1>
			<span
				class="text-4xl font-black tracking-widest uppercase border-[3px] px-3 py-1 drop-shadow-[0_0_8px_currentColor]
				{$serverStatus === 'Running'
					? 'text-[#00f0ff] border-[#00f0ff]/40 shadow-[0_0_15px_rgba(0,240,255,0.5)]'
					: $serverStatus === 'Starting'
						? 'text-[#ffcf54] border-[#ffcf54]/40 shadow-[0_0_15px_rgba(255,207,84,0.5)]'
						: 'text-[#ff0055] border-[#ff0055]/40 shadow-[0_0_15px_rgba(255,0,85,0.5)]'}"
			>
				{$serverStatus === 'Running' ? 'ONLINE' : $serverStatus}
			</span>
		</div>
		<div class="mt-3 text-[#00f0ff] font-bold text-sm tracking-wider uppercase drop-shadow-[0_0_4px_rgba(0,240,255,0.6)]">
			CLUSTER: US-EAST-01 // IP: 192.168.1.104:25565
		</div>
	</div>

	<div class="flex gap-4">
		<button
			onclick={start}
			disabled={$serverStatus !== 'Stopped' || $loading}
			class="flex items-center gap-2 px-8 py-3 bg-[#00f0ff] text-[#1a1130] font-black tracking-widest uppercase transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white hover:shadow-[0_0_20px_#00f0ff] hover:scale-105"
		>
			<span class="text-xs">▶</span> START
		</button>
		<button
			onclick={stop}
			disabled={$serverStatus !== 'Running' || $loading}
			class="flex items-center gap-2 px-8 py-3 bg-[#ff0055] text-white font-black tracking-widest uppercase transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white hover:text-[#ff0055] hover:shadow-[0_0_20px_#ff0055] hover:scale-105"
		>
			<span class="text-xs">■</span> STOP
		</button>
		<button
			onclick={restart}
			disabled={$serverStatus === 'Stopped' || $loading}
			class="flex items-center gap-2 px-8 py-3 bg-[#a355ff] text-white font-black tracking-widest uppercase transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white hover:text-[#a355ff] hover:shadow-[0_0_20px_#a355ff] hover:scale-105"
		>
			<span class="text-lg leading-none">↻</span> RESTART
		</button>
	</div>
</div>
