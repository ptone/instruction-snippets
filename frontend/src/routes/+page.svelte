<script lang="ts">
	import { onMount, tick } from 'svelte';
	import {
		collection,
		getDocs,
		doc,
		updateDoc,
		increment,
		type DocumentData
	} from 'firebase/firestore';
	import { db } from '$lib/firebase';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import { ThumbsUp, ThumbsDown, Copy, X } from 'lucide-svelte';

	let snippets: DocumentData[] = [];
	let searchTerm = '';
	let sortBy = 'newest';
	let activeFilters = new Set<string>();
	let copiedState: { [key: string]: boolean } = {};

	onMount(async () => {
		const querySnapshot = await getDocs(collection(db, 'snippets'));
		snippets = querySnapshot.docs.map((doc) => {
			const data = doc.data();
			return {
				id: doc.id,
				...data,
				createdAt: data.createdAt.toDate() // Convert Firestore Timestamp to JS Date
			};
		});
	});

	$: sortedAndFilteredSnippets = (() => {
		let result = snippets.filter((snippet) => {
			const lowerCaseSearch = searchTerm.toLowerCase();
			const inText = snippet.content.toLowerCase().includes(lowerCaseSearch);
			const inLabels = snippet.labels.some((label: string) =>
				label.toLowerCase().includes(lowerCaseSearch)
			);

			const hasActiveLabels = activeFilters.size === 0 ||
				[...activeFilters].every(filter => snippet.labels.includes(filter));

			return (inText || inLabels) && hasActiveLabels;
		});

		result.sort((a, b) => {
			switch (sortBy) {
				case 'newest':
					return b.createdAt - a.createdAt;
				case 'oldest':
					return a.createdAt - b.createdAt;
				case 'most_liked':
					return (b.thumbs_up - b.thumbs_down) - (a.thumbs_up - a.thumbs_down);
				case 'least_liked':
					return (a.thumbs_up - a.thumbs_down) - (b.thumbs_up - b.thumbs_down);
				default:
					return 0;
			}
		});

		return result;
	})();

	async function rateSnippet(id: string, rating: 'thumbs_up' | 'thumbs_down') {
		const snippetRef = doc(db, 'snippets', id);
		await updateDoc(snippetRef, {
			[rating]: increment(1)
		});
		// Optimistically update the UI
		const snippet = snippets.find(s => s.id === id);
		if(snippet) {
			snippet[rating]++;
			snippets = [...snippets];
		}
	}

	function copyToClipboard(id: string, text: string) {
		navigator.clipboard.writeText(text);
		copiedState[id] = true;
		setTimeout(() => {
			copiedState[id] = false;
			// Use tick to ensure the UI updates before we potentially re-render
			tick();
		}, 2000);
	}

	function toggleFilter(label: string) {
		if (activeFilters.has(label)) {
			activeFilters.delete(label);
		} else {
			activeFilters.add(label);
		}
		activeFilters = new Set(activeFilters);
	}
</script>

<div class="container mx-auto p-4">
	<div class="flex justify-between items-center mb-4">
		<h1 class="text-2xl font-bold">Instruction Snippets</h1>
	</div>

	<div class="flex space-x-4 mb-4">
		<div class="flex-grow">
			<Input placeholder="Search snippets..." bind:value={searchTerm} />
		</div>
		<Select.Root bind:value={sortBy}>
			<Select.Trigger class="w-[180px]">
				<Select.Value placeholder="Sort by" />
			</Select.Trigger>
			<Select.Content>
				<Select.Item value="newest">Newest</Select.Item>
				<Select.Item value="oldest">Oldest</Select.Item>
				<Select.Item value="most_liked">Most Liked</Select.Item>
				<Select.Item value="least_liked">Least Liked</Select.Item>
			</Select.Content>
		</Select.Root>
	</div>

	<div class="mb-4 flex flex-wrap gap-2">
		{#each [...activeFilters] as filter}
			<Badge class="cursor-pointer bg-blue-500 text-white hover:bg-blue-600"
				   on:click={() => toggleFilter(filter)}>
				{filter} <X class="ml-2 h-4 w-4" />
			</Badge>
		{/each}
	</div>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
		{#each sortedAndFilteredSnippets as snippet}
			<Card.Root>
				<Card.Header>
					<Card.Title>Snippet</Card.Title>
				</Card.Header>
				<Card.Content>
					<p class="whitespace-pre-wrap">{snippet.content}</p>
				</Card.Content>
				<Card.Footer class="flex-col items-start">
					<div class="mt-2 flex flex-wrap">
						{#each snippet.labels as label}
							<Badge
								class="mr-2 mb-2 cursor-pointer"
								variant={activeFilters.has(label) ? 'secondary' : 'outline'}
								on:click={() => toggleFilter(label)}>{label}</Badge
							>
						{/each}
					</div>
					<div class="flex items-center space-x-2 self-end">
						<Button
							variant="ghost"
							size="icon"
							on:click={() => rateSnippet(snippet.id, 'thumbs_up')}
						>
							<ThumbsUp class="h-4 w-4" />
						</Button>
						<span>{snippet.thumbs_up}</span>
						<Button
							variant="ghost"
							size="icon"
							on:click={() => rateSnippet(snippet.id, 'thumbs_down')}
						>
							<ThumbsDown class="h-4 w-4" />
						</Button>
						<span>{snippet.thumbs_down}</span>
						<div class="relative">
							<Button
								variant="ghost"
								size="icon"
								on:click={() => copyToClipboard(snippet.id, snippet.content)}
							>
								<Copy class="h-4 w-4" />
							</Button>
							{#if copiedState[snippet.id]}
								<span
									class="absolute -top-8 left-1/2 -translate-x-1/2 rounded-md bg-gray-900 px-2 py-1 text-xs text-white"
									>Copied!</span
								>
							{/if}
						</div>
					</div>
				</Card.Footer>
			</Card.Root>
		{/each}
	</div>
</div>
