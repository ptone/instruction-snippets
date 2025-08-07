<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { badgeVariants, type Variant } from '.';
	import { cn } from '$lib/utils';

	type $$Props = HTMLAttributes<HTMLDivElement> & {
		variant?: Variant;
	};

	let className: $$Props['class'] = undefined;
	let variant: $$Props['variant'] = 'default';
	export { className as class, variant };
</script>

<div class={cn(badgeVariants({ variant }), className)}>
	<slot />
</div>
