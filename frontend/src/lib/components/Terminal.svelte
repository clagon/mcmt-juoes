<script lang="ts">
	import { logs, serverStatus } from '$lib/store';
	import { onMount } from 'svelte';

	let logContainer: HTMLElement | undefined = $state(undefined);
	let command = $state('');

	// History state
	let commandHistory: string[] = $state([]);
	let historyIndex: number = $state(-1);
	let currentDraft: string = $state('');

	onMount(async () => {
		try {
			const res = await fetch('http://localhost:8080/api/server/logs');
			const initialLogs = await res.json();
			if (initialLogs && initialLogs.length > 0) {
				logs.set(initialLogs);
			}
		} catch (e) {
			console.error('Failed to fetch initial logs', e);
		}
	});

	$effect(() => {
		if ($logs.length && logContainer) {
			logContainer.scrollTop = logContainer.scrollHeight;
		}
	});

	function handleKeydown(e: KeyboardEvent) {
		const target = e.target as HTMLInputElement;
		if ($serverStatus !== 'Running') return;

		if (e.key === 'ArrowUp') {
			if (target.selectionStart === 0 && target.selectionEnd === 0) {
				e.preventDefault();
				if (commandHistory.length > 0) {
					if (historyIndex === -1) {
						currentDraft = command;
						historyIndex = commandHistory.length - 1;
					} else if (historyIndex > 0) {
						historyIndex--;
					}
					command = commandHistory[historyIndex];
				}
			} else {
				// Prevent default jump and manually move cursor to start
				e.preventDefault();
				target.setSelectionRange(0, 0);
			}
		} else if (e.key === 'ArrowDown') {
			const end = command.length;
			if (target.selectionStart === end && target.selectionEnd === end) {
				e.preventDefault();
				if (historyIndex !== -1) {
					if (historyIndex < commandHistory.length - 1) {
						historyIndex++;
						command = commandHistory[historyIndex];
					} else {
						historyIndex = -1;
						command = currentDraft;
					}
				}
			} else {
				// Prevent default jump and manually move cursor to end
				e.preventDefault();
				target.setSelectionRange(end, end);
			}
		}
	}

	async function sendCommand(e: Event) {
		e.preventDefault();
		const cmd = command.trim();
		if (!cmd || $serverStatus !== 'Running') return;

		try {
			await fetch('http://localhost:8080/api/server/command', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ command: cmd })
			});

			if (commandHistory.length === 0 || commandHistory[commandHistory.length - 1] !== cmd) {
				commandHistory.push(cmd);
			}

			historyIndex = -1;
			currentDraft = '';
			command = '';
		} catch (err) {
			console.error('Failed to send command', err);
		}
	}
</script>

<div
	class="col-span-1 md:col-span-2 flex flex-col bg-black/80 backdrop-blur-sm border-2 border-vapor-cyan shadow-[0_0_15px_var(--color-vapor-cyan)] rounded-lg overflow-hidden h-[600px] z-10 relative"
>
	<!-- CRT screen overlay -->
	<div
		class="absolute inset-0 pointer-events-none z-20"
		style="background: linear-gradient(rgba(18, 16, 16, 0) 50%, rgba(0, 0, 0, 0.25) 50%), linear-gradient(90deg, rgba(255, 0, 0, 0.06), rgba(0, 255, 0, 0.02), rgba(0, 0, 255, 0.06)); background-size: 100% 2px, 3px 100%;"
	></div>

	<!-- Terminal Header -->
	<div
		class="bg-vapor-bg/80 border-b-2 border-vapor-cyan px-4 py-2 flex items-center justify-between z-10"
	>
		<span
			class="text-vapor-cyan font-mono text-sm uppercase tracking-widest"
			style="text-shadow: 0 0 5px var(--color-vapor-cyan);"
		>
			>_ SERVER_TERMINAL.EXE
		</span>
		<div class="flex gap-2">
			<div class="w-3 h-3 rounded-full bg-red-500"></div>
			<div class="w-3 h-3 rounded-full bg-yellow-500"></div>
			<div class="w-3 h-3 rounded-full bg-green-500"></div>
		</div>
	</div>

	<!-- Log Output -->
	<div
		bind:this={logContainer}
		class="flex-grow p-4 font-mono text-sm overflow-y-auto bg-black text-green-400 z-10"
	>
		{#each $logs as log, i (i)}
			<div
				class="whitespace-pre-wrap break-words border-b border-green-900/30 pb-1 mb-1 font-mono text-xs"
			>
				{log}
			</div>
		{/each}
	</div>

	<!-- Command Input -->
	<form
		onsubmit={sendCommand}
		class="border-t-2 border-vapor-cyan flex items-center bg-gray-900 px-4 py-3 z-10"
	>
		<span
			class="text-vapor-pink mr-2 font-bold"
			style="text-shadow: 0 0 5px var(--color-vapor-pink);">></span
		>
		<input
			type="text"
			bind:value={command}
			onkeydown={handleKeydown}
			placeholder={$serverStatus === 'Running' ? 'Enter server command...' : 'Server is offline...'}
			disabled={$serverStatus !== 'Running'}
			class="flex-grow bg-transparent text-vapor-cyan placeholder-gray-600 outline-none font-mono disabled:opacity-50 disabled:cursor-not-allowed uppercase"
		/>
		<button
			type="submit"
			disabled={$serverStatus !== 'Running' || !command.trim()}
			class="ml-2 px-6 py-2 bg-vapor-purple text-black font-black uppercase tracking-wider text-sm hover:bg-vapor-pink hover:shadow-[0_0_15px_var(--color-vapor-pink)] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
		>
			EXEC
		</button>
	</form>
</div>

<style>
	/* Custom scrollbar for terminal */
	::-webkit-scrollbar {
		width: 10px;
	}
	::-webkit-scrollbar-track {
		background: #000;
		border-left: 1px solid var(--color-vapor-cyan);
	}
	::-webkit-scrollbar-thumb {
		background: var(--color-vapor-purple);
	}
</style>
