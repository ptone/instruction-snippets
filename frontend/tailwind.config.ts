import type { Config } from 'tailwindcss';

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
	plugins: []
} satisfies Config;
