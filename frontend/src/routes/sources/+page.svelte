<script lang="ts">
	import { onMount } from 'svelte';
	import { authUser } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { db } from '$lib/firebase';
	import { collection, query, where, getDocs, orderBy } from 'firebase/firestore';
	import type { DocumentData, QuerySnapshot } from 'firebase/firestore';

	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let sources: DocumentData[] = [];
	let loading = true;
	let showMySources = true;

	onMount(async () => {
		if (!$authUser) {
			goto('/');
		} else {
			await fetchSources();
		}
	});

	async function fetchSources() {
		loading = true;
		let q = query(collection(db, 'sources'), orderBy('last_refreshed', 'desc'));

		if (showMySources && $authUser) {
			q = query(q, where('submitterId', '==', $authUser.uid));
		}

		const querySnapshot: QuerySnapshot<DocumentData> = await getDocs(q);
		sources = querySnapshot.docs.map((doc) => ({ id: doc.id, ...doc.data() }));
		loading = false;
	}

	function handleToggle() {
		showMySources = !showMySources;
		fetchSources();
	}
</script>

<div class="mx-auto max-w-4xl">
	<div class="mb-4 flex items-center justify-between">
		<h2 class="text-2xl font-bold">Sources</h2>
		<div
			class="flex items-center space-x-2"
			on:click={handleToggle}
			on:keydown={(e) => e.key === 'Enter' && handleToggle()}
			role="button"
			tabindex="0"
		>
			<Button variant={showMySources ? 'secondary' : 'outline'}>
				{showMySources ? 'Show all submissions' : 'Show my submissions only'}
			</Button>
		</div>
	</div>

	<div class="mb-4">
		<a href="/contribute">
			<Button>Add Contribution</Button>
		</a>
	</div>

	{#if loading}
		<p>Loading...</p>
	{:else if sources.length === 0}
		<p>No sources found.</p>
	{:else}
		<div class="grid gap-4">
			{#each sources as source}
				<Card.Root>
					<Card.Header>
						<Card.Title class="flex items-center justify-between">
							<span>{source.key}</span>
							<Badge>{source.status}</Badge>
						</Card.Title>
						<Card.Description>
							Submitted by: {source.submitterEmail}
						</Card.Description>
					</Card.Header>
					<Card.Content>
						{#if source.type === 'url'}
							<a
								href={source.url}
								target="_blank"
								rel="noopener noreferrer"
								class="text-blue-500 hover:underline"
							>
								{source.url}
							</a>
						{:else}
							<p>File upload</p>
						{/if}
					</Card.Content>
					<Card.Footer>
						<p class="text-sm text-muted-foreground">
							Last refreshed: {new Date(source.last_refreshed.seconds * 1000).toLocaleString()}
						</p>
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	{/if}
</div>
