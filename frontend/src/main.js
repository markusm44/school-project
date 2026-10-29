import './app.css'
import App from './App.svelte'

export const BackendHost = "http://localhost:8080"
export const WSBackendHost = "ws://localhost:8080"

const app = new App({
  target: document.getElementById('app'),
})

export default app
