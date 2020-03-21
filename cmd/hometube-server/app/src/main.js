import App from './App.svelte';

const app = new App({
	target: document.body,
	props: {
		apiBaseUrl: "http://192.168.0.23:8080/api/v1"
	}
});

export default app;
