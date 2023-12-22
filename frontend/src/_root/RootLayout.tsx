// Libraries
import React from 'react'
import { Outlet } from 'react-router-dom'

// Components
import Header from '@/components/shared/Header'
import SidebarLeft from '@/components/shared/SidebarLeft'
import Footer from '@/components/shared/Footer'

const RootLayout = () => {
  return (
    <div className='w-full md:flex'>
      <Header />
      <SidebarLeft />
      <section className='flex flex-1 h-full'>
        <Outlet />
      </section>
      <Footer />
    </div>
  )
}

export default RootLayout