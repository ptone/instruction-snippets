<script lang="ts">
	import { onMount } from 'svelte';
	import { collection, getDocs, type DocumentData } from 'firebase/firestore';
	import { db } from '$lib/firebase';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';

	let snippets: DocumentData[] = [];

	onMount(async () => {
		const querySnapshot = await getDocs(collection(db, 'snippets'));
		snippets = querySnapshot.docs.map((doc) => doc.data());
	});
</script>

<div class="container mx-auto p-4">
	<h1 class="mb-4 text-2xl font-bold">Instruction Snippets</h1>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
		{#each snippets as snippet}
			<Card.Root>
				<Card.Header>
					<Card.Title>Snippet</Card.Title>
				</Card.Header>
				<Card.Content>
					<p>{snippet.text}</p>
				</Card.Content>
				<Card.Footer>
					<div class="mt-2 flex flex-wrap">
						{#each snippet.labels as label}
							<Badge class="mr-2 mb-2">{label}</Badge>
						{/each}
					</div>
				</Card.Footer>
			</Card.Root>
		{/each}
	</div>
</div>
