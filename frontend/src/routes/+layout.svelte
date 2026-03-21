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
	class="min-h-screen bg-[#0d0814] text-[#00f0ff] font-sans relative overflow-hidden flex flex-col items-center py-8 px-4"
>
	<!-- Sun -->
	<div class="absolute top-10 left-20 z-[-2] pointer-events-none opacity-60">
		<div class="w-48 h-48 rounded-full bg-gradient-to-b from-[#ffcf54] to-[#ff0055] blur-[2px]"></div>
	</div>
    <!-- Triangle -->
	<div class="absolute top-10 right-32 z-[-2] pointer-events-none opacity-40">
		<div class="w-0 h-0 border-l-[60px] border-r-[60px] border-t-[100px] border-l-transparent border-r-transparent border-t-[#a355ff] blur-[4px]"></div>
	</div>

	<!-- Background grid -->
	<div
		class="absolute bottom-0 left-0 right-0 h-[60%] z-[-1] opacity-30 pointer-events-none"
		style="background-image: linear-gradient(#ff0055 2px, transparent 2px), linear-gradient(90deg, #ff0055 2px, transparent 2px); background-size: 50px 50px; transform: perspective(500px) rotateX(60deg) translateY(50px) translateZ(0px); transform-origin: top center;"
	></div>
    <!-- Horizon line -->
    <div class="absolute top-[40%] left-0 right-0 h-[2px] bg-[#ff0055] z-[-1] shadow-[0_0_15px_#ff0055]"></div>

	<main class="w-full max-w-7xl flex-grow gap-6 grid grid-cols-1 md:grid-cols-3 z-10 mt-2">
		{@render children()}
	</main>
</div>
