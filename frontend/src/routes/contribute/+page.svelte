<script lang="ts">
	import { authUser } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import * as ToggleGroup from '$lib/components/ui/toggle-group';

	let source_type = 'url'; // or 'file'
	let url = '';
	let content = '';
	let key = '';
	let limit = 0;
	let fileInput: HTMLInputElement;

	onMount(() => {
		if (!$authUser) {
			goto('/');
		}
	});

	const handleFileChange = () => {
		const file = fileInput.files?.[0];
		if (file) {
			const reader = new FileReader();
			reader.onload = (e) => {
				content = e.target?.result as string;
			};
			reader.readAsText(file);
		}
	};

	const handleSubmit = async () => {
		const body = {
			key: key || (source_type === 'url' ? url : fileInput.files?.[0].name),
			limit: limit || undefined,
			url: source_type === 'url' ? url : undefined,
			content: source_type === 'file' ? content : undefined
		};

		await fetch('/api/process', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
		});

		// maybe show a notification
		goto('/');
	};
</script>

<div class="mx-auto max-w-lg">
	<h2 class="mb-4 text-2xl font-bold">Add a new Contribution</h2>

	<form on:submit|preventDefault={handleSubmit}>
		<div class="mb-4">
			<label for="key">Source Key (optional)</label>
			<Input id="key" bind:value={key} placeholder="e.g., my-awesome-source" />
			<p class="mt-1 text-sm text-muted-foreground">
				A unique key for this source. If not provided, the URL or filename will be used.
			</p>
		</div>

		<div class="mb-4">
			<label for="source-type">Source Type</label>
			<ToggleGroup.Root id="source-type" type="single" bind:value={source_type} class="mt-2">
				<ToggleGroup.Item value="url">URL</ToggleGroup.Item>
				<ToggleGroup.Item value="file">File</ToggleGroup.Item>
			</ToggleGroup.Root>
		</div>

		{#if source_type === 'url'}
			<div class="mb-4">
				<label for="url-input">URL</label>
				<Input id="url-input" type="url" bind:value={url} required placeholder="https://..." />
			</div>
		{:else}
			<div class="mb-4">
				<label for="file-input">File</label>
				<input
					id="file-input"
					type="file"
					bind:this={fileInput}
					on:change={handleFileChange}
					required
				/>
			</div>
		{/if}

		<div class="mb-4">
			<label for="limit">Snippet Limit (optional)</label>
			<Input id="limit" type="number" bind:value={limit} placeholder="0" />
			<p class="mt-1 text-sm text-muted-foreground">
				Limit the number of snippets generated. 0 means no limit.
			</p>
		</div>

		<Button type="submit">Submit</Button>
	</form>
</div>
