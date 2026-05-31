import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { RestfulProvider } from 'restful-react'
import App from './App'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RestfulProvider base="/api/v1">
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </RestfulProvider>
  </React.StrictMode>
)
