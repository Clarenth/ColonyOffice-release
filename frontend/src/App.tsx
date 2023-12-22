// Libraries
import React from 'react'
import { Route, Routes } from 'react-router-dom'

// Components
import Login from "./Auth/forms/Login";
import Signup from "./Auth/forms/Signup";
import { Toaster } from './components/ui/toaster'

// Layouts
import AuthLayout from "./Auth/AuthLayout";
import RootLayout from "./_root/RootLayout";

// Pages
import { Account, AI, Chat, Docs, DocsCreate, DocsDetails, DocsEdit, Home, Search } from "./_root/pages";

// Routes
//import routes from './routes/routes'

// styles
import './global.css'

const App = () => {
  return (
    <main className='flex h-screen'>
      <Routes>
        <Route>
          {/* public routes */}
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />
          </Route>
          
          {/* private routes */}
          <Route element={<RootLayout />}>
            <Route index element={<Home />} />
            <Route path='/account/:id_code' element={<Account />} />
            <Route path='/account-update/:document_id' element={<Account />} />
            <Route path='/docs' element={<Docs />} />
            <Route path='/docs/:document_id' element={<DocsDetails />} />
            <Route path='/create-docs' element={<DocsCreate />} />
            <Route path='/edit-docs/:document_id' element={<DocsEdit />} />
            <Route path='/search' element={<Search />} />
            <Route path='/chat' element={<Chat />} />
            <Route path='/ai' element={<AI />} />
          </Route>
        </Route>
      </Routes>
      {/*<RouterProvider router={ routes } />*/}
      <Toaster />
    </main>
  )
}

export default App