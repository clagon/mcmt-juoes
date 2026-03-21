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
	class="flex flex-col bg-[#110e19]/90 border-[3px] border-[#00f0ff] rounded-[12px] overflow-hidden h-[500px] z-10 relative shadow-[0_0_15px_rgba(0,240,255,0.4)]"
>
	<!-- CRT screen overlay -->
	<div
		class="absolute inset-0 pointer-events-none z-20 mix-blend-screen opacity-10"
		style="background: linear-gradient(rgba(18, 16, 16, 0) 50%, rgba(0, 0, 0, 0.25) 50%), linear-gradient(90deg, rgba(255, 0, 0, 0.06), rgba(0, 255, 0, 0.02), rgba(0, 0, 255, 0.06)); background-size: 100% 2px, 3px 100%;"
	></div>

	<!-- Terminal Header -->
	<div
		class="bg-[#1a1130] border-b-[3px] border-[#00f0ff] px-4 py-2 flex items-center justify-between z-10"
	>
		<span
			class="text-[#00f0ff] font-mono text-sm uppercase tracking-widest font-bold"
		>
			SERVER_TERMINAL.EXE
		</span>
		<div class="flex gap-2">
			<div class="w-3 h-3 rounded-full bg-[#3aff5c] shadow-[0_0_5px_#3aff5c]"></div>
			<div class="w-3 h-3 rounded-full bg-[#ffcf54] shadow-[0_0_5px_#ffcf54]"></div>
			<div class="w-3 h-3 rounded-full bg-[#ff0055] shadow-[0_0_5px_#ff0055]"></div>
		</div>
	</div>

	<!-- Log Output -->
	<div
		bind:this={logContainer}
		class="flex-grow p-4 font-mono text-sm overflow-y-auto bg-[#110e19] text-[#71d283] z-10 custom-scrollbar"
	>
		{#each $logs as log, i (i)}
			<div
				class="whitespace-pre-wrap break-words border-[#1d2720] pb-[2px] mb-[2px] font-mono text-[13px] leading-tight"
			>
				{log}
			</div>
		{/each}
        <div class="h-4 w-2 bg-[#71d283] animate-pulse mt-1 inline-block"></div>
	</div>

	<!-- Command Input -->
	<form
		onsubmit={sendCommand}
		class="border-t-[3px] border-[#00f0ff] flex items-center bg-[#1a1130] p-3 z-10"
	>
		<div class="flex-grow border-[2px] border-[#42395d] bg-transparent rounded px-3 py-2 flex items-center group focus-within:border-[#00f0ff] transition-colors">
            <span
                class="text-[#00f0ff] mr-3 font-bold font-mono"
                >></span>
            <input
                type="text"
                bind:value={command}
                onkeydown={handleKeydown}
                placeholder={$serverStatus === 'Running' ? 'ENTER SERVER COMMAND...' : 'SERVER IS OFFLINE...'}
                disabled={$serverStatus !== 'Running'}
                class="w-full bg-transparent text-[#00f0ff] placeholder-[#42395d] outline-none font-mono disabled:opacity-50 disabled:cursor-not-allowed uppercase text-sm"
            />
        </div>
		<button
			type="submit"
			disabled={$serverStatus !== 'Running' || !command.trim()}
			class="ml-3 px-6 py-2 bg-transparent border-[2px] border-[#00f0ff] text-[#00f0ff] rounded font-black uppercase tracking-wider text-sm hover:bg-[#00f0ff] hover:text-[#110e19] hover:shadow-[0_0_15px_#00f0ff] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
		>
			EXEC
		</button>
	</form>
</div>

<style>
	/* Custom scrollbar for terminal */
	.custom-scrollbar::-webkit-scrollbar {
		width: 12px;
	}
	.custom-scrollbar::-webkit-scrollbar-track {
		background: #110e19;
		border-left: 2px solid #28203c;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: #00f0ff;
        border-radius: 4px;
        border: 2px solid #110e19;
	}
</style>
