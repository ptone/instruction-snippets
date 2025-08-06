<script lang="ts">
	import { onMount } from 'svelte';
	import { collection, getDocs, doc, updateDoc, increment, type DocumentData } from 'firebase/firestore';
	import { db } from '$lib/firebase';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { ThumbsUp, ThumbsDown, Copy } from 'lucide-svelte';

	let snippets: DocumentData[] = [];
	let searchTerm = '';
	let copiedState: { [key: string]: boolean } = {};

	onMount(async () => {
		const querySnapshot = await getDocs(collection(db, 'snippets'));
		snippets = querySnapshot.docs.map((doc) => {
			return { id: doc.id, ...doc.data() };
		});
	});

	$: filteredSnippets = snippets.filter((snippet) => {
		const lowerCaseSearch = searchTerm.toLowerCase();
		const inText = snippet.text.toLowerCase().includes(lowerCaseSearch);
		const inLabels = snippet.labels.some((label: string) =>
			label.toLowerCase().includes(lowerCaseSearch)
		);
		return inText || inLabels;
	});

	async function rateSnippet(id: string, rating: 'thumbs_up' | 'thumbs_down') {
		const snippetRef = doc(db, 'snippets', id);
		await updateDoc(snippetRef, {
			[rating]: increment(1)
		});
	}

	function copyToClipboard(id: string, text: string) {
		navigator.clipboard.writeText(text);
		copiedState[id] = true;
		setTimeout(() => {
			copiedState[id] = false;
		}, 2000);
	}
</script>

<div class="container mx-auto p-4">
	<h1 class="mb-4 text-2xl font-bold">Instruction Snippets</h1>

	<div class="mb-4">
		<Input placeholder="Search snippets..." bind:value={searchTerm} />
	</div>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
		{#each filteredSnippets as snippet}
			<Card.Root>
				<Card.Header>
					<Card.Title>Snippet</Card.Title>
				</Card.Header>
				<Card.Content>
					<p>{snippet.text}</p>
				</Card.Content>
				<Card.Footer class="flex justify-between">
					<div class="mt-2 flex flex-wrap">
						{#each snippet.labels as label}
							<Badge class="mr-2 mb-2">{label}</Badge>
						{/each}
					</div>
					<div class="flex items-center space-x-2">
						<Button variant="ghost" size="icon" onclick={() => rateSnippet(snippet.id, 'thumbs_up')}>
							<ThumbsUp class="h-4 w-4" />
						</Button>
						<span>{snippet.thumbs_up}</span>
						<Button variant="ghost" size="icon" onclick={() => rateSnippet(snippet.id, 'thumbs_down')}>
							<ThumbsDown class="h-4 w-4" />
						</Button>
						<span>{snippet.thumbs_down}</span>
						<div class="relative">
							<Button variant="ghost" size="icon" onclick={() => copyToClipboard(snippet.id, snippet.text)}>
								<Copy class="h-4 w-4" />
							</Button>
							{#if copiedState[snippet.id]}
								<span class="absolute -top-8 left-1/2 -translate-x-1/2 rounded-md bg-gray-900 px-2 py-1 text-xs text-white">Copied!</span>
							{/if}
						</div>
					</div>
				</Card.Footer>
			</Card.Root>
		{/each}
	</div>
</div>
