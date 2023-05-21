import { ThemeProvider } from '@material-tailwind/react'
import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { MaterialTailwindControllerProvider } from "./context/index"
import App from './App'
import './index.css'
import { ToastContainer } from 'react-toastify';

import 'react-toastify/dist/ReactToastify.css';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
    <BrowserRouter>
    <ThemeProvider>
      <MaterialTailwindControllerProvider>
        <App />
        <ToastContainer />
      </MaterialTailwindControllerProvider>
    </ThemeProvider>
  </BrowserRouter>
)
