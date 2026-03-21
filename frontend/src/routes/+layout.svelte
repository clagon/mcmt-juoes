<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { connectWebSocket, serverStatus } from '$lib/store';

	onMount(() => {
		connectWebSocket();
	});

	let { children } = $props();
</script>

<div
	class="min-h-screen bg-vapor-bg text-vapor-cyan font-sans relative overflow-hidden flex flex-col items-center py-10 px-4"
>
	<!-- Background grid -->
	<div
		class="absolute inset-0 z-[-1] opacity-20 pointer-events-none"
		style="background-image: linear-gradient(var(--color-vapor-pink) 1px, transparent 1px), linear-gradient(90deg, var(--color-vapor-pink) 1px, transparent 1px); background-size: 50px 50px; transform: perspective(500px) rotateX(60deg) translateY(100px) translateZ(-200px);"
	></div>

	<header class="mb-10 text-center relative z-10 w-full">
		<h1
			class="text-4xl md:text-6xl font-black text-vapor-pink glow-pink tracking-widest uppercase italic border-b-4 border-vapor-purple pb-2 inline-block"
		>
			Server Manager
		</h1>
		<div class="flex items-center justify-center mt-4 gap-4">
			<span class="text-vapor-cyan glow-cyan tracking-widest text-lg font-mono">STATUS:</span>
			<span
				class="px-4 py-1 border-2 font-bold tracking-widest uppercase
				{$serverStatus === 'Running'
					? 'text-green-400 border-green-400'
					: $serverStatus === 'Starting'
						? 'text-yellow-400 border-yellow-400'
						: 'text-red-400 border-red-400'}"
			>
				{$serverStatus}
			</span>
		</div>
	</header>

	<main class="w-full max-w-6xl flex-grow gap-6 grid grid-cols-1 md:grid-cols-3 z-10">
		{@render children()}
	</main>
</div>
