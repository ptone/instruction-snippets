import type { Config } from 'tailwindcss';
import typography from '@tailwindcss/typography';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				primary: {
					DEFAULT: '#000000',
					foreground: '#ffffff'
				}
			}
		}
	},
	plugins: [typography]
} satisfies Config;
