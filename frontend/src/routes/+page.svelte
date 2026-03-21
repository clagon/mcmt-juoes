<script lang="ts">
	import Controls from '$lib/components/Controls.svelte';
	import Terminal from '$lib/components/Terminal.svelte';
	import Management from '$lib/components/Management.svelte';
</script>

<div class="col-span-1 md:col-span-3">
    <Controls />
</div>
<div class="col-span-1 md:col-span-2">
    <Terminal />
</div>
<div class="col-span-1 md:col-span-1">
    <Management />
</div>
