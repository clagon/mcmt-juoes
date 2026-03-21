<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { connectWebSocket } from '$lib/store';

	onMount(() => {
		connectWebSocket();
	});

	let { children } = $props();
</script>

<div
	class="min-h-screen bg-vapor-bg text-vapor-cyan font-sans relative overflow-hidden flex flex-col items-center py-6 px-4"
>
	<!-- Background grid -->
	<div
		class="absolute inset-0 z-[-1] opacity-20 pointer-events-none"
		style="background-image: linear-gradient(var(--color-vapor-pink) 1px, transparent 1px), linear-gradient(90deg, var(--color-vapor-pink) 1px, transparent 1px); background-size: 50px 50px; transform: perspective(500px) rotateX(60deg) translateY(100px) translateZ(-200px);"
	></div>

	<!-- Top Bar Container (Header + Controls will fit within layout flow or inside page)
         In the current architecture, +layout provides the frame.
         We'll move the header content entirely into +page.svelte to make the top bar
         a single cohesive unit with the Controls.
         Wait, let's keep the global status here but styled differently, or better yet,
         move all header UI into +page.svelte so Controls and Header are siblings in a flex container.
         For now, let's render children without a large global header,
         and let +page.svelte handle the top bar. -->

	<main class="w-full max-w-6xl flex-grow z-10 flex flex-col gap-6">
		{@render children()}
	</main>
</div>
