import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { RestfulProvider } from 'restful-react'
import App from './App'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    {/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
    {/* @ts-ignore restful-react + React 18 类型兼容 */}
    <RestfulProvider base="">
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </RestfulProvider>
  </React.StrictMode>
)
