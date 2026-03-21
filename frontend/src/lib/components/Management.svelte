<script lang="ts">
	import { serverStatus } from '$lib/store';
	import { onMount } from 'svelte';

	let activeTab: 'players' | 'known' | 'settings' = $state('players');

	// Types based on backend structs
	type Player = {
		uuid: string;
		name: string;
	};

	type OpsEntry = {
		uuid: string;
		name: string;
		level: number;
		bypassesPlayerLimit: boolean;
	};

	type BannedPlayer = {
		uuid: string;
		name: string;
		created: string;
		source: string;
		expires: string;
		reason: string;
	};

	// State
	let onlinePlayers: string[] = $state([]);
	let whitelist: Player[] = $state([]);
	let ops: OpsEntry[] = $state([]);
	let bannedPlayers: BannedPlayer[] = $state([]);

	// Settings state
	let javaXms: string = $state('');
	let javaXmx: string = $state('');
	let javaArgs: string = $state('');
	let properties: string = $state('');
	let loadingSettings: boolean = $state(false);

	let inputPlayer: string = $state('');

	// Derived state (using functions in Runes to keep it simple, or derived if complex)
	let allPlayers = $derived.by(() => {
		const names = new Set<string>();
		whitelist.forEach((p) => names.add(p.name));
		ops.forEach((p) => names.add(p.name));
		bannedPlayers.forEach((p) => names.add(p.name));
		return Array.from(names).sort();
	});

	async function fetchData() {
		try {
			const [onlineRes, whitelistRes, opsRes, bannedRes] = await Promise.all([
				fetch('http://localhost:8080/api/server/players'),
				fetch('http://localhost:8080/api/server/whitelist'),
				fetch('http://localhost:8080/api/server/ops'),
				fetch('http://localhost:8080/api/server/banned-players')
			]);

			onlinePlayers = (await onlineRes.json()) || [];
			whitelist = (await whitelistRes.json()) || [];
			ops = (await opsRes.json()) || [];
			bannedPlayers = (await bannedRes.json()) || [];
		} catch (e) {
			console.error('Failed to fetch player data', e);
		}
	}

	async function fetchSettings() {
		try {
			const res = await fetch('http://localhost:8080/api/server/settings');
			const data = await res.json();
			javaXms = data.javaXms || '';
			javaXmx = data.javaXmx || '';
			javaArgs = data.javaArgs || '';
			properties = data.properties || '';
		} catch (e) {
			console.error('Failed to fetch settings', e);
		}
	}

	onMount(() => {
		fetchData();
		fetchSettings();
		// Poll for online players if running
		const interval = setInterval(() => {
			if ($serverStatus === 'Running') {
				fetchData();
			}
		}, 5000);
		return () => clearInterval(interval);
	});

	// Helpers
	function isOp(name: string) {
		return ops.some((p) => p.name.toLowerCase() === name.toLowerCase());
	}
	function isWhitelisted(name: string) {
		return whitelist.some((p) => p.name.toLowerCase() === name.toLowerCase());
	}
	function isBanned(name: string) {
		return bannedPlayers.some((p) => p.name.toLowerCase() === name.toLowerCase());
	}

	async function execute(cmd: string, target: string) {
		if ($serverStatus !== 'Running') return;
		try {
			await fetch('http://localhost:8080/api/server/command', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ command: `${cmd} ${target}` })
			});
			// Optimistically refresh after a short delay
			setTimeout(fetchData, 1000);
		} catch (e) {
			console.error(`Failed to execute ${cmd} on ${target}`, e);
		}
	}

	async function saveSettings() {
		loadingSettings = true;
		try {
			await fetch('http://localhost:8080/api/server/settings', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ javaXms, javaXmx, javaArgs })
			});
		} catch (e) {
			console.error('Failed to save settings', e);
		} finally {
			loadingSettings = false;
		}
	}

	function handleAddAction(action: 'whitelist' | 'op' | 'ban') {
		const target = inputPlayer.trim();
		if (!target) return;
		if (action === 'whitelist') execute('whitelist add', target);
		if (action === 'op') execute('op', target);
		if (action === 'ban') execute('ban', target);
		inputPlayer = '';
	}
</script>

<div
	class="flex flex-col bg-[#110e19]/90 border-[3px] border-[#ff0055] rounded-[12px] overflow-hidden h-[500px] z-10 shadow-[0_0_15px_rgba(255,0,85,0.4)]"
>
	<!-- Tabs -->
	<div class="flex border-b-[3px] border-[#ff0055] bg-[#1a1130] p-2 gap-2">
		<button
			onclick={() => (activeTab = 'players')}
			class="flex-1 py-2 font-bold text-sm tracking-widest uppercase rounded border-[2px] transition-all
			{activeTab === 'players'
				? 'bg-[#ff0055] text-white border-[#ff0055] shadow-[0_0_10px_rgba(255,0,85,0.6)]'
				: 'bg-transparent text-[#ff0055] border-[#ff0055] hover:bg-[#ff0055]/20'}"
		>
			PLAYERS
		</button>
		<button
			onclick={() => (activeTab = 'known')}
			class="flex-1 py-2 font-bold text-sm tracking-widest uppercase rounded border-[2px] transition-all
			{activeTab === 'known'
				? 'bg-[#00f0ff] text-[#110e19] border-[#00f0ff] shadow-[0_0_10px_rgba(0,240,255,0.6)]'
				: 'bg-transparent text-[#00f0ff] border-[#00f0ff] hover:bg-[#00f0ff]/20'}"
		>
			KNOWN
		</button>
		<button
			onclick={() => (activeTab = 'settings')}
			class="flex-1 py-2 font-bold text-sm tracking-widest uppercase rounded border-[2px] transition-all
			{activeTab === 'settings'
				? 'bg-[#a355ff] text-white border-[#a355ff] shadow-[0_0_10px_rgba(163,85,255,0.6)]'
				: 'bg-transparent text-[#a355ff] border-[#a355ff] hover:bg-[#a355ff]/20'}"
		>
			SETTINGS
		</button>
	</div>

	<div class="flex-grow overflow-y-auto p-4 custom-scrollbar text-sm bg-[#110e19]">
		{#if activeTab === 'players'}
			<div class="space-y-4">
                <div class="border-[2px] border-[#42395d] rounded bg-[#1a1130] p-3">
                    <div class="text-[#7f74a8] text-xs font-bold mb-2 uppercase tracking-wider">TARGET PLAYER: USERNAME...</div>
                </div>

				<div
					class="text-[#ff0055] font-bold border-b-[2px] border-[#2a1b3d] pb-2 flex justify-between items-center bg-[#1a1130] px-3 py-2 rounded-t border-[2px] border-b-0 border-[#3a448a]"
				>
					<span class="tracking-widest">PLAYER</span>
					<span class="text-white text-xs tracking-widest">
						ACTIONS
					</span>
				</div>

                <div class="border-[2px] border-[#3a448a] border-t-0 rounded-b p-3 bg-[#110e19]">
				{#if onlinePlayers.length === 0}
					<div class="text-[#42395d] italic text-center py-4">No players currently online.</div>
				{:else}
					<ul class="space-y-3">
						{#each onlinePlayers as player (player)}
							<li
								class="flex justify-between items-center border-[2px] border-[#42395d] p-3 rounded bg-[#1a1130]"
							>
								<div class="flex flex-col">
									<span class="text-white font-bold flex items-center gap-3">
										{player}
									</span>
									<div class="flex items-center gap-1 mt-1">
                                        <span class="text-[#ff0055] text-[10px] tracking-wider">ONLINE: </span>
                                        <div class="w-2 h-2 rounded-full bg-[#ff0055] shadow-[0_0_4px_#ff0055]"></div>
									</div>
								</div>

								<div class="flex gap-2">
                                    <button
                                        disabled={$serverStatus !== 'Running'}
                                        onclick={() => execute('whitelist add', player)}
                                        class="border-[2px] border-[#3aff5c] text-[#3aff5c] rounded px-3 py-1 text-xs font-bold tracking-wider uppercase hover:bg-[#3aff5c] hover:text-[#110e19] disabled:opacity-50"
                                        >WHITELIST</button
                                    >
									{#if isOp(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('deop', player)}
											class="border-[2px] border-[#00f0ff] text-[#00f0ff] rounded px-3 py-1 text-xs font-bold tracking-wider uppercase hover:bg-[#00f0ff] hover:text-[#110e19] disabled:opacity-50"
											>DE-OP</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('op', player)}
											class="border-[2px] border-[#00f0ff] text-[#00f0ff] rounded px-3 py-1 text-xs font-bold tracking-wider uppercase hover:bg-[#00f0ff] hover:text-[#110e19] disabled:opacity-50"
											>OP</button
										>
									{/if}

									<button
										disabled={$serverStatus !== 'Running'}
										onclick={() => execute('kick', player)}
										class="border-[2px] border-[#ff0055] text-[#ff0055] rounded px-3 py-1 text-xs font-bold tracking-wider uppercase hover:bg-[#ff0055] hover:text-white disabled:opacity-50"
										>KICK</button
									>
                                    <button
										disabled={$serverStatus !== 'Running'}
										onclick={() => execute('ban', player)}
										class="border-[2px] border-[#ff0055] text-[#ff0055] rounded px-3 py-1 text-xs font-bold tracking-wider uppercase hover:bg-[#ff0055] hover:text-white disabled:opacity-50"
										>BAN</button
									>
								</div>
							</li>
						{/each}
					</ul>
				{/if}
                </div>
			</div>
		{/if}

		{#if activeTab === 'known'}
			<div class="space-y-6">
				<!-- Add Player Action Box -->
				<div class="border-[2px] border-[#00f0ff] p-4 rounded bg-[#1a1130] shadow-[0_0_10px_rgba(0,240,255,0.2)]">
					<div class="text-[#00f0ff] mb-2 font-bold tracking-widest text-xs">TARGET PLAYER:</div>
					<form
						onsubmit={(e) => {
							e.preventDefault();
							handleAddAction('whitelist');
						}}
                        class="mb-3"
					>
						<input
							type="text"
							bind:value={inputPlayer}
							placeholder="USERNAME..."
							class="w-full bg-[#110e19] border-[2px] border-[#42395d] rounded p-2 text-[#00f0ff] outline-none uppercase placeholder-[#42395d] focus:border-[#00f0ff] transition-colors"
						/>
					</form>
					<div class="flex gap-2 justify-between">
						<button
							disabled={$serverStatus !== 'Running'}
							onclick={() => handleAddAction('whitelist')}
							class="flex-1 border-[2px] border-[#3aff5c] text-[#3aff5c] rounded text-xs font-bold hover:bg-[#3aff5c] hover:text-[#110e19] py-2 disabled:opacity-50"
							>WHITELIST</button
						>
						<button
							disabled={$serverStatus !== 'Running'}
							onclick={() => handleAddAction('op')}
							class="flex-1 border-[2px] border-[#00f0ff] text-[#00f0ff] rounded text-xs font-bold hover:bg-[#00f0ff] hover:text-[#110e19] py-2 disabled:opacity-50"
							>OP</button
						>
						<button
							disabled={$serverStatus !== 'Running'}
							onclick={() => handleAddAction('ban')}
							class="flex-1 border-[2px] border-[#ff0055] text-[#ff0055] rounded text-xs font-bold hover:bg-[#ff0055] hover:text-white py-2 disabled:opacity-50"
							>BAN</button
						>
					</div>
				</div>

				<div class="text-[#00f0ff] font-bold border-b-[2px] border-[#00f0ff] pb-1 tracking-widest drop-shadow-[0_0_5px_rgba(0,240,255,0.5)]">
					KNOWN PLAYERS
				</div>

				{#if allPlayers.length === 0}
					<div class="text-[#42395d] italic text-center py-4">No data found in JSON files.</div>
				{:else}
					<ul class="space-y-3">
						{#each allPlayers as player (player)}
							<li class="flex flex-col border-[2px] border-[#42395d] p-3 rounded bg-[#1a1130]">
								<div class="flex justify-between items-center mb-3 border-b-[2px] border-[#2a1b3d] pb-2">
									<span class="text-white font-bold tracking-wider">
										{player}
									</span>
									<div class="flex gap-1 text-[10px] font-bold">
										{#if isOp(player)}
											<span
												class="bg-[#00f0ff]/20 text-[#00f0ff] px-2 py-1 rounded border-[1px] border-[#00f0ff]/50"
												>OP</span
											>
										{/if}
										{#if isWhitelisted(player)}
											<span class="bg-[#3aff5c]/20 text-[#3aff5c] px-2 py-1 rounded border-[1px] border-[#3aff5c]/50"
												>WL</span
											>
										{/if}
										{#if isBanned(player)}
											<span class="bg-[#ff0055]/20 text-[#ff0055] px-2 py-1 rounded border-[1px] border-[#ff0055]/50"
												>BAN</span
											>
										{/if}
									</div>
								</div>

								<div class="grid grid-cols-2 gap-2 text-xs font-bold">
									{#if isWhitelisted(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('whitelist remove', player)}
											class="border-[2px] border-[#3aff5c] text-[#3aff5c] rounded hover:bg-[#3aff5c] hover:text-[#110e19] py-1 disabled:opacity-50"
											>UN-WHITELIST</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('whitelist add', player)}
											class="border-[2px] border-[#3aff5c] text-[#3aff5c] rounded hover:bg-[#3aff5c] hover:text-[#110e19] py-1 disabled:opacity-50"
											>WHITELIST</button
										>
									{/if}

									{#if isOp(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('deop', player)}
											class="border-[2px] border-[#00f0ff] text-[#00f0ff] rounded hover:bg-[#00f0ff] hover:text-[#110e19] py-1 disabled:opacity-50"
											>DE-OP</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('op', player)}
											class="border-[2px] border-[#00f0ff] text-[#00f0ff] rounded hover:bg-[#00f0ff] hover:text-[#110e19] py-1 disabled:opacity-50"
											>OP</button
										>
									{/if}

									{#if isBanned(player)}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('pardon', player)}
											class="border-[2px] border-[#ff0055] text-[#ff0055] rounded hover:bg-[#ff0055] hover:text-white py-1 disabled:opacity-50 col-span-2 mt-1"
											>UNBAN</button
										>
									{:else}
										<button
											disabled={$serverStatus !== 'Running'}
											onclick={() => execute('ban', player)}
											class="border-[2px] border-[#ff0055] text-[#ff0055] rounded hover:bg-[#ff0055] hover:text-white py-1 disabled:opacity-50 col-span-2 mt-1"
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
				<div class="text-[#a355ff] font-bold border-b-[2px] border-[#a355ff] pb-1 tracking-widest drop-shadow-[0_0_5px_rgba(163,85,255,0.5)]">JVM ARGUMENTS</div>
				<form
					onsubmit={(e) => {
						e.preventDefault();
						saveSettings();
					}}
					class="space-y-4 bg-[#1a1130] p-4 rounded border-[2px] border-[#42395d]"
				>
					<div>
						<label for="initial-memory" class="block text-[#7f74a8] font-bold mb-1 text-xs uppercase"
							>Initial Memory (-Xms)</label
						>
						<input
							id="initial-memory"
							type="text"
							bind:value={javaXms}
							placeholder="E.G. 2G"
							class="w-full bg-[#110e19] border-[2px] border-[#42395d] rounded p-2 text-[#a355ff] outline-none uppercase placeholder-[#42395d] focus:border-[#a355ff] transition-colors"
						/>
					</div>
					<div>
						<label for="max-memory" class="block text-[#7f74a8] font-bold mb-1 text-xs uppercase"
							>Max Memory (-Xmx)</label
						>
						<input
							id="max-memory"
							type="text"
							bind:value={javaXmx}
							placeholder="E.G. 4G"
							class="w-full bg-[#110e19] border-[2px] border-[#42395d] rounded p-2 text-[#a355ff] outline-none uppercase placeholder-[#42395d] focus:border-[#a355ff] transition-colors"
						/>
					</div>
					<div>
						<label for="additional-args" class="block text-[#7f74a8] font-bold mb-1 text-xs uppercase"
							>Additional Arguments</label
						>
						<input
							id="additional-args"
							type="text"
							bind:value={javaArgs}
							placeholder="E.G. -XX:+UseG1GC"
							class="w-full bg-[#110e19] border-[2px] border-[#42395d] rounded p-2 text-[#a355ff] outline-none uppercase placeholder-[#42395d] focus:border-[#a355ff] transition-colors"
						/>
					</div>
					<button
						type="submit"
						disabled={loadingSettings}
						class="w-full border-[2px] border-[#a355ff] text-[#a355ff] hover:bg-[#a355ff] hover:text-white font-bold py-2 rounded transition-all disabled:opacity-50 uppercase tracking-widest mt-2"
					>
						{loadingSettings ? 'SAVING...' : 'SAVE SETTINGS'}
					</button>
					<div class="text-xs text-[#ffcf54] mt-2 italic text-center tracking-wider font-bold">
						* REQUIRES SERVER RESTART *
					</div>
				</form>

				<div class="text-[#a355ff] font-bold border-b-[2px] border-[#a355ff] pb-1 mt-6 tracking-widest drop-shadow-[0_0_5px_rgba(163,85,255,0.5)]">
					SERVER.PROPERTIES
				</div>
				<textarea
					readonly
					class="w-full h-48 bg-[#110e19] border-[2px] border-[#42395d] rounded p-3 text-[#a355ff] text-[10px] outline-none custom-scrollbar whitespace-pre font-mono"
					value={properties}
				></textarea>
				<div class="text-[10px] text-[#7f74a8] text-center font-bold uppercase tracking-widest">
					(READ-ONLY. EDIT FILE DIRECTLY ON HOST.)
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.custom-scrollbar::-webkit-scrollbar {
		width: 12px;
	}
	.custom-scrollbar::-webkit-scrollbar-track {
		background: #110e19;
        border-left: 2px solid #28203c;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: #ff0055;
        border-radius: 4px;
        border: 2px solid #110e19;
	}
    /* Dynamic tab scrollbar based on active tab would require inline styles or svelte classes, keeping it pink for now as default/fallback or using activeTab color in a style tag via variable is complex in standard svelte without CSS variables, so sticking to pink for the scrollbar thumb in management panel is a safe bet, or making it a neutral vaporwave color. Let's make it the dark purple to match the border */
    .custom-scrollbar::-webkit-scrollbar-thumb {
		background: #42395d;
	}
</style>
