<script lang="ts">
	import Controls from '$lib/components/Controls.svelte';
	import Terminal from '$lib/components/Terminal.svelte';
	import Management from '$lib/components/Management.svelte';
</script>

<div class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
	<Controls />
	<Terminal />
	<Management />
</div>
