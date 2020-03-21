import App from './App.svelte';

const app = new App({
	target: document.body,
	props: {
		apiBaseUrl: "http://localhost:8080/api/v1"
	}
});

export default app;
