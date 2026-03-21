<script lang="ts">
	import { serverStatus } from '$lib/store';
	import { onMount, onDestroy } from 'svelte';

	type Player = { uuid: string; name: string };

	let whitelist: Player[] = $state([]);
	let ops: Player[] = $state([]);
	let banned: Player[] = $state([]);
	let properties: string = $state('');
	let onlinePlayers: string[] = $state([]);

	let activeTab: 'players' | 'known' | 'settings' = $state('players');

	// Refresh timer
	let interval: ReturnType<typeof setInterval>;

	async function fetchData() {
		try {
			const [wl, op, ban, props, online] = await Promise.all([
				fetch('http://localhost:8080/api/server/whitelist').then((r) => r.json()),
				fetch('http://localhost:8080/api/server/ops').then((r) => r.json()),
				fetch('http://localhost:8080/api/server/banned-players').then((r) => r.json()),
				fetch('http://localhost:8080/api/server/properties').then((r) => r.json()),
				fetch('http://localhost:8080/api/server/online').then((r) => r.json())
			]);
			whitelist = wl || [];
			ops = op || [];
			banned = ban || [];
			properties = props.content || '';
			onlinePlayers = online || [];
		} catch (e) {
			console.error('Failed to fetch server data', e);
		}
	}

	onMount(() => {
		fetchData();
		interval = setInterval(fetchData, 5000); // Poll every 5 seconds
	});

	onDestroy(() => {
		clearInterval(interval);
	});

	async function sendCommand(cmd: string) {
		if ($serverStatus !== 'Running') return;
		try {
			await fetch('http://localhost:8080/api/server/command', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ command: cmd })
			});
			// Optimistically wait and fetch data
			setTimeout(fetchData, 1000);
		} catch (e) {
			console.error('Failed to execute management command', e);
		}
	}

	function execute(action: string, player: string) {
		sendCommand(`${action} ${player}`);
	}

	// Simple helpers to check status
	const isOp = (name: string) => ops.some((p) => p.name === name);
	const isWhitelisted = (name: string) => whitelist.some((p) => p.name === name);
	const isBanned = (name: string) => banned.some((p) => p.name === name);
	const isOnline = (name: string) => onlinePlayers.includes(name);

	let inputPlayer = $state('');
	function handleAddAction(action: 'whitelist' | 'op' | 'ban') {
		if (!inputPlayer) return;
		if (action === 'whitelist') execute('whitelist add', inputPlayer);
		if (action === 'op') execute('op', inputPlayer);
		if (action === 'ban') execute('ban', inputPlayer);
		inputPlayer = '';
	}

	// Settings
	let javaXms = $state('');
	let javaXmx = $state('');
	let javaArgs = $state('');
	let loadingSettings = $state(false);

	async function loadSettings() {
		try {
			const res = await fetch('http://localhost:8080/api/settings');
			const data = await res.json();
			javaXms = data.java_xms;
			javaXmx = data.java_xmx;
			javaArgs = data.java_args;
		} catch (e) {
			console.error('Failed to load settings', e);
		}
	}

	async function saveSettings() {
		loadingSettings = true;
		try {
			await fetch('http://localhost:8080/api/settings', {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ java_xms: javaXms, java_xmx: javaXmx, java_args: javaArgs })
			});
		} catch (e) {
			console.error('Failed to save settings', e);
		} finally {
			loadingSettings = false;
		}
	}

	// Load settings on mount as well
	onMount(() => {
		loadSettings();
	});

	// Derived list of all known players across files
	let allPlayers: string[] = $derived.by(() => {
		const names = new Set([
			...whitelist.map((p) => p.name),
			...ops.map((p) => p.name),
			...banned.map((p) => p.name)
		]);
		return Array.from(names);
	});
</script>

<div
	class="col-span-1 border-2 border-vapor-pink rounded-lg bg-black/80 backdrop-blur-sm shadow-[0_0_15px_var(--color-vapor-pink)] flex flex-col h-[600px] z-10 overflow-hidden font-mono"
>
	<!-- Header / Tabs -->
	<div class="flex border-b-2 border-vapor-pink text-xs sm:text-sm">
		<button
			onclick={() => (activeTab = 'players')}
			class="flex-1 py-3 font-bold text-center tracking-widest transition-colors {activeTab ===
			'players'
				? 'bg-vapor-pink text-black'
				: 'text-vapor-pink hover:bg-vapor-pink/20'}"
		>
			PLAYERS
		</button>
		<button
			onclick={() => (activeTab = 'known')}
			class="flex-1 py-3 font-bold text-center tracking-widest border-l-2 border-vapor-pink transition-colors {activeTab ===
			'known'
				? 'bg-vapor-pink text-black'
				: 'text-vapor-pink hover:bg-vapor-pink/20'}"
		>
			KNOWN
		</button>
		<button
			onclick={() => (activeTab = 'settings')}
			class="flex-1 py-3 font-bold text-center tracking-widest border-l-2 border-vapor-pink transition-colors {activeTab ===
			'settings'
				? 'bg-vapor-pink text-black'
				: 'text-vapor-pink hover:bg-vapor-pink/20'}"
		>
			SETTINGS
		</button>
	</div>

	<div class="flex-grow overflow-y-auto p-4 custom-scrollbar text-sm">
		{#if activeTab === 'players'}
			<div class="space-y-6">
				<div
					class="text-vapor-pink glow-pink font-bold border-b border-vapor-pink pb-1 flex justify-between items-center"
				>
					<span>ONLINE PLAYERS</span>
					<span class="text-green-400 text-xs">
						{$serverStatus === 'Running' ? `ONLINE: ${onlinePlayers.length}` : 'SERVER OFFLINE'}
					</span>
				</div>

				{#if onlinePlayers.length === 0}
					<div class="text-gray-500 italic text-center py-4">No players currently online.</div>
				{:else}
					<ul class="space-y-3">
						{#each onlinePlayers as player (player)}
							<li
								class="flex flex-col border border-green-500/50 p-2 rounded bg-green-900/10 shadow-[0_0_5px_rgba(34,197,94,0.2)]"
							>
								<div class="flex justify-between items-center mb-2">
									<span class="text-vapor-cyan font-bold flex items-center gap-3">
										<img
											src={`https://mc-heads.net/avatar/${player}/32`}
											alt={player}
											class="w-8 h-8 rounded pixelated border border-gray-600"
										/>
										{player}
									</span>
									<div class="flex gap-1 text-[10px]">
										{#if isOp(player)}
											<span
												class="bg-yellow-400/20 text-yellow-400 px-1 border border-yellow-400/50"
												>OP</span
											>
										{/if}
										{#if isWhitelisted(player)}
											<span class="bg-green-400/20 text-green-400 px-1 border border-green-400/50"
												>WL</span
											>
										{/if}
									</div>
								</div>

								<div class="grid grid-cols-2 gap-1 text-[10px]">
									{#if isOp(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('deop', player)}
											class="border border-yellow-400 text-yellow-400 hover:bg-yellow-400 hover:text-black py-1 disabled:opacity-50"
											>DE-OP</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('op', player)}
											class="border border-yellow-400 text-yellow-400 hover:bg-yellow-400 hover:text-black py-1 disabled:opacity-50"
											>OP</button
										>
									{/if}

									<button
										disabled={$serverStatus !== 'Running'}
										onclick={() => execute('kick', player)}
										class="border border-orange-400 text-orange-400 hover:bg-orange-400 hover:text-black py-1 disabled:opacity-50"
										>KICK</button
									>
								</div>
							</li>
						{/each}
					</ul>
				{/if}
			</div>
		{/if}

		{#if activeTab === 'known'}
			<div class="space-y-6">
				<!-- Add Player Action Box -->
				<div class="border border-vapor-purple p-3 rounded bg-black/50">
					<div class="text-vapor-cyan mb-2">TARGET PLAYER:</div>
					<form
						onsubmit={(e) => {
							e.preventDefault();
							handleAddAction('whitelist');
						}}
					>
						<input
							type="text"
							bind:value={inputPlayer}
							placeholder="Username... (Press Enter to Whitelist)"
							class="w-full bg-gray-900 border border-vapor-cyan p-2 text-vapor-cyan outline-none mb-3 uppercase placeholder-gray-600"
						/>
					</form>
					<div class="grid grid-cols-3 gap-2">
						<button
							disabled={$serverStatus !== 'Running'}
							onclick={() => handleAddAction('whitelist')}
							class="border border-vapor-cyan text-vapor-cyan hover:bg-vapor-cyan hover:text-black py-1 disabled:opacity-50"
							>WHITELIST</button
						>
						<button
							disabled={$serverStatus !== 'Running'}
							onclick={() => handleAddAction('op')}
							class="border border-yellow-400 text-yellow-400 hover:bg-yellow-400 hover:text-black py-1 disabled:opacity-50"
							>OP</button
						>
						<button
							disabled={$serverStatus !== 'Running'}
							onclick={() => handleAddAction('ban')}
							class="border border-red-500 text-red-500 hover:bg-red-500 hover:text-black py-1 disabled:opacity-50"
							>BAN</button
						>
					</div>
				</div>

				<div class="text-vapor-pink glow-pink font-bold border-b border-vapor-pink pb-1">
					KNOWN PLAYERS
				</div>

				{#if allPlayers.length === 0}
					<div class="text-gray-500 italic text-center py-4">No data found in JSON files.</div>
				{:else}
					<ul class="space-y-3">
						{#each allPlayers as player (player)}
							<li class="flex flex-col border border-gray-800 p-2 rounded bg-gray-900/50">
								<div class="flex justify-between items-center mb-2">
									<span class="text-vapor-cyan font-bold flex items-center gap-3">
										<img
											src={`https://mc-heads.net/avatar/${player}/32`}
											alt={player}
											class="w-8 h-8 rounded pixelated border border-gray-600"
										/>
										{player}
									</span>
									<div class="flex gap-1 text-[10px]">
										{#if isOp(player)}
											<span
												class="bg-yellow-400/20 text-yellow-400 px-1 border border-yellow-400/50"
												>OP</span
											>
										{/if}
										{#if isWhitelisted(player)}
											<span class="bg-green-400/20 text-green-400 px-1 border border-green-400/50"
												>WL</span
											>
										{/if}
										{#if isBanned(player)}
											<span class="bg-red-500/20 text-red-500 px-1 border border-red-500/50"
												>BAN</span
											>
										{/if}
									</div>
								</div>

								<div class="grid grid-cols-2 gap-1 text-[10px]">
									{#if isWhitelisted(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('whitelist remove', player)}
											class="border border-vapor-cyan text-vapor-cyan hover:bg-vapor-cyan hover:text-black py-1 disabled:opacity-50"
											>UN-WHITELIST</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('whitelist add', player)}
											class="border border-vapor-cyan text-vapor-cyan hover:bg-vapor-cyan hover:text-black py-1 disabled:opacity-50"
											>WHITELIST</button
										>
									{/if}

									{#if isOp(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('deop', player)}
											class="border border-yellow-400 text-yellow-400 hover:bg-yellow-400 hover:text-black py-1 disabled:opacity-50"
											>DE-OP</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('op', player)}
											class="border border-yellow-400 text-yellow-400 hover:bg-yellow-400 hover:text-black py-1 disabled:opacity-50"
											>OP</button
										>
									{/if}

									{#if isBanned(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('pardon', player)}
											class="border border-red-500 text-red-500 hover:bg-red-500 hover:text-black py-1 disabled:opacity-50 col-span-2"
											>UNBAN</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('ban', player)}
											class="border border-red-500 text-red-500 hover:bg-red-500 hover:text-black py-1 disabled:opacity-50 col-span-2"
											>BAN</button
										>
									{/if}
								</div>
							</li>
						{/each}
					</ul>
				{/if}
			</div>
		{/if}

		{#if activeTab === 'settings'}
			<div class="space-y-6">
				<div class="text-vapor-cyan font-bold border-b border-vapor-cyan pb-1">JVM ARGUMENTS</div>
				<form
					onsubmit={(e) => {
						e.preventDefault();
						saveSettings();
					}}
					class="space-y-4"
				>
					<div>
						<label for="initial-memory" class="block text-gray-400 mb-1 text-xs"
							>Initial Memory (-Xms)</label
						>
						<input
							id="initial-memory"
							type="text"
							bind:value={javaXms}
							placeholder="e.g. 2G"
							class="w-full bg-gray-900 border border-vapor-cyan p-2 text-vapor-cyan outline-none"
						/>
					</div>
					<div>
						<label for="max-memory" class="block text-gray-400 mb-1 text-xs"
							>Max Memory (-Xmx)</label
						>
						<input
							id="max-memory"
							type="text"
							bind:value={javaXmx}
							placeholder="e.g. 4G"
							class="w-full bg-gray-900 border border-vapor-cyan p-2 text-vapor-cyan outline-none"
						/>
					</div>
					<div>
						<label for="additional-args" class="block text-gray-400 mb-1 text-xs"
							>Additional Arguments</label
						>
						<input
							id="additional-args"
							type="text"
							bind:value={javaArgs}
							placeholder="e.g. -XX:+UseG1GC"
							class="w-full bg-gray-900 border border-vapor-cyan p-2 text-vapor-cyan outline-none"
						/>
					</div>
					<button
						type="submit"
						disabled={loadingSettings}
						class="w-full bg-vapor-purple hover:bg-vapor-pink text-black font-bold py-2 mt-2 transition-colors disabled:opacity-50"
					>
						{loadingSettings ? 'SAVING...' : 'SAVE SETTINGS'}
					</button>
					<div class="text-xs text-red-400 mt-1 italic text-center">
						* Requires server restart *
					</div>
				</form>

				<div class="text-vapor-cyan font-bold border-b border-vapor-cyan pb-1 mt-6">
					SERVER.PROPERTIES
				</div>
				<textarea
					readonly
					class="w-full h-48 bg-gray-900 border border-gray-700 p-2 text-green-400 text-xs outline-none custom-scrollbar whitespace-pre font-mono"
					value={properties}
				></textarea>
				<div class="text-xs text-gray-500 text-center">
					(Read-only. Edit file directly on host.)
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.custom-scrollbar::-webkit-scrollbar {
		width: 8px;
	}
	.custom-scrollbar::-webkit-scrollbar-track {
		background: #000;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: var(--color-vapor-pink);
	}
</style>
