<script lang="ts">
	import { authUser } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import { Label } from '$lib/components/ui/label';

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

<div class="max-w-lg mx-auto">
	<h2 class="text-2xl font-bold mb-4">Add a new Contribution</h2>

	<form on:submit|preventDefault={handleSubmit}>
		<div class="mb-4">
			<Label for="key">Source Key (optional)</Label>
			<Input id="key" bind:value={key} placeholder="e.g., my-awesome-source" />
			<p class="text-sm text-muted-foreground mt-1">
				A unique key for this source. If not provided, the URL or filename will be used.
			</p>
		</div>

		<div class="mb-4">
			<Label>Source Type</Label>
			<RadioGroup.Root bind:value={source_type} class="mt-2">
				<div class="flex items-center space-x-2">
					<RadioGroup.Item value="url" id="url" />
					<Label for="url">URL</Label>
				</div>
				<div class="flex items-center space-x-2">
					<RadioGroup.Item value="file" id="file" />
					<Label for="file">File</Label>
				</div>
			</RadioGroup.Root>
		</div>

		{#if source_type === 'url'}
			<div class="mb-4">
				<Label for="url-input">URL</Label>
				<Input id="url-input" type="url" bind:value={url} required placeholder="https://..." />
			</div>
		{:else}
			<div class="mb-4">
				<Label for="file-input">File</Label>
				<Input id="file-input" type="file" bind:this={fileInput} on:change={handleFileChange} required />
			</div>
		{/if}

		<div class="mb-4">
			<Label for="limit">Snippet Limit (optional)</Label>
			<Input id="limit" type="number" bind:value={limit} placeholder="0" />
			<p class="text-sm text-muted-foreground mt-1">
				Limit the number of snippets generated. 0 means no limit.
			</p>
		</div>

		<Button type="submit">Submit</Button>
	</form>
</div>
