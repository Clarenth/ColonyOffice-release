/* Libraries */
// React
import React from 'react';
import ReactDOM from 'react-dom/client';

// Mantine
//import { MantineProvider } from '@mantine/core'; 

// Context
import { AuthProvider } from './context/AuthContext';

// Components
import App from './App';
import { QueryProvider } from './lib/react-query/QueryProvider';
import { BrowserRouter } from 'react-router-dom';

// Styles


ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <QueryProvider>
        <AuthProvider>
          <App />
        </AuthProvider>
      </QueryProvider>
    </BrowserRouter>
  </React.StrictMode>
);